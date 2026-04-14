package pg

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/adamchenEpm/ym3-go/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgUtil struct {
	db *sql.DB
}

var (
	instance *PgUtil
	once     sync.Once
)

// GetInstance 获取 PostgreSQL 单例（自动初始化连接池）
func GetInstance() *PgUtil {
	once.Do(func() {
		instance = &PgUtil{}
		if err := instance.init(); err != nil {
			panic(fmt.Errorf("初始化 PostgreSQL 失败: %w", err))
		}
	})
	return instance
}

// init 初始化数据库连接池
func (p *PgUtil) init() error {
	cfg := config.Get() // 从你的配置中读取 PG 配置
	// 连接字符串格式：postgres://user:pass@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Pg.User, cfg.Pg.Password, cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.Name)
	var err error
	p.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	// 连接池配置
	p.db.SetMaxOpenConns(50)
	p.db.SetMaxIdleConns(10)
	p.db.SetConnMaxLifetime(30 * time.Minute)
	if err = p.db.Ping(); err != nil {
		return err
	}
	return nil
}

// GetDB 返回原始 *sql.DB 对象
func (p *PgUtil) GetDB() *sql.DB {
	return p.db
}

// Close 关闭数据库连接
func (p *PgUtil) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// ================== 基础查询方法 ==================

// QueryRow 查询单行，扫描到 dest 指针列表
func (p *PgUtil) QueryRow(query string, dest ...interface{}) error {
	row := p.db.QueryRow(query)
	return row.Scan(dest...)
}

// QueryRows 查询多行，逐行调用回调函数
func (p *PgUtil) QueryRows(query string, callback func(*sql.Rows) error) error {
	rows, err := p.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := callback(rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

// QueryToStructs 将查询结果映射到结构体切片（使用反射）
// 结构体字段 tag 为 `pg:"column_name"`
func (p *PgUtil) QueryToStructs(query string, dest interface{}, args ...interface{}) error {
	rows, err := p.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 获取列名
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	// 反射 dest（必须是切片指针）
	slicePtr := reflect.ValueOf(dest)
	if slicePtr.Kind() != reflect.Ptr || slicePtr.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest 必须是切片的指针")
	}
	sliceVal := slicePtr.Elem()
	elemType := sliceVal.Type().Elem()

	// 构建列名到字段索引的映射
	colToFieldIdx := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("pg")
		if tag == "" {
			tag = field.Name
		}
		colToFieldIdx[tag] = i
	}

	// 遍历结果集
	for rows.Next() {
		elemPtr := reflect.New(elemType)
		elem := elemPtr.Elem()
		// 准备扫描用的 interface{} 切片
		scanArgs := make([]interface{}, len(cols))
		for i, col := range cols {
			if idx, ok := colToFieldIdx[col]; ok {
				scanArgs[i] = elem.Field(idx).Addr().Interface()
			} else {
				var tmp interface{}
				scanArgs[i] = &tmp
			}
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, elem))
	}
	return rows.Err()
}

// ================== 执行命令方法 ==================

// Exec 执行 INSERT/UPDATE/DELETE，返回 sql.Result
func (p *PgUtil) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(query, args...)
}

// Insert 插入并返回自增ID（PostgreSQL 使用 RETURNING id）
func (p *PgUtil) Insert(query string, args ...interface{}) (int64, error) {
	var id int64
	err := p.db.QueryRow(query+" RETURNING id", args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Update 执行更新，返回受影响行数
func (p *PgUtil) Update(query string, args ...interface{}) (int64, error) {
	result, err := p.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ================== 事务支持 ==================

// Transaction 执行事务，回调函数中执行数据库操作
func (p *PgUtil) Transaction(fn func(tx *sql.Tx) error) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// ================== 批量插入辅助 ==================

// BatchInsert 批量插入数据（使用 PostgreSQL 的 COPY 或拼接多行 VALUES）
// 这里使用多行 VALUES 方式，注意参数占位符为 $1,$2...
func (p *PgUtil) BatchInsert(table string, columns []string, rows [][]interface{}) error {
	if len(rows) == 0 {
		return nil
	}
	// 构建 VALUES 占位符：($1,$2),($3,$4)...
	placeholders := make([]string, 0, len(rows))
	paramIndex := 1
	for _, row := range rows {
		ph := make([]string, len(row))
		for i := range row {
			ph[i] = fmt.Sprintf("$%d", paramIndex)
			paramIndex++
		}
		placeholders = append(placeholders, "("+strings.Join(ph, ",")+")")
	}
	// 扁平化参数列表
	flatArgs := make([]interface{}, 0, len(rows)*len(columns))
	for _, row := range rows {
		flatArgs = append(flatArgs, row...)
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		table, strings.Join(columns, ","), strings.Join(placeholders, ","))
	_, err := p.db.Exec(query, flatArgs...)
	return err
}

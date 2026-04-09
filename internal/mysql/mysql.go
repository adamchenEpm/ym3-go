package mysql

import (
	"database/sql"
	"fmt"
	"github.com/adamchenEpm/ym3-go/internal/config"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLUtil MySQL 操作工具类
type MySQLUtil struct {
	db *sql.DB
}

var (
	instance *MySQLUtil
	once     sync.Once
)

// GetInstance 获取 MySQL 工具单例（自动初始化连接池）
func GetInstance() *MySQLUtil {
	once.Do(func() {
		instance = &MySQLUtil{}
		if err := instance.init(); err != nil {
			panic(fmt.Errorf("初始化 MySQL 失败: %w", err))
		}
	})
	return instance
}

// init 内部初始化数据库连接
func (m *MySQLUtil) init() error {
	// 从全局配置中读取 MySQL 配置（可根据实际调整）
	cfg := config.Get()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	var err error
	m.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 设置连接池参数
	m.db.SetMaxOpenConns(50)                  // 最大打开连接数
	m.db.SetMaxIdleConns(10)                  // 最大空闲连接数
	m.db.SetConnMaxLifetime(30 * time.Minute) // 连接最大生存时间
	if err = m.db.Ping(); err != nil {
		return err
	}
	return nil
}

// GetDB 返回原始 *sql.DB 对象（供高级操作）
func (m *MySQLUtil) GetDB() *sql.DB {
	return m.db
}

// Close 关闭数据库连接（优雅退出时调用）
func (m *MySQLUtil) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// ================== 基础查询方法 ==================

// QueryRow 查询单行，自动扫描到 dest
// dest 为指针列表，例如：&id, &name
func (m *MySQLUtil) QueryRow(query string, dest ...interface{}) error {
	row := m.db.QueryRow(query)
	return row.Scan(dest...)
}

// QueryRows 查询多行，逐行调用回调函数
// callback 接收 *sql.Rows，需在内部 Scan 并处理
func (m *MySQLUtil) QueryRows(query string, callback func(*sql.Rows) error) error {
	rows, err := m.db.Query(query)
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
// 结构体字段 tag 为 `db:"column_name"`
func (m *MySQLUtil) QueryToStructs(query string, dest interface{}, args ...interface{}) error {
	rows, err := m.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 获取列信息
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	// 反射处理 dest（必须是切片指针）
	slicePtr := reflect.ValueOf(dest)
	if slicePtr.Kind() != reflect.Ptr || slicePtr.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest 必须是切片的指针")
	}
	sliceVal := slicePtr.Elem()
	elemType := sliceVal.Type().Elem()

	// 构建列到字段索引的映射
	colToFieldIdx := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			tag = field.Name
		}
		colToFieldIdx[tag] = i
	}

	// 准备扫描用的 interface{} 切片
	for rows.Next() {
		elemPtr := reflect.New(elemType)
		elem := elemPtr.Elem()
		// 为每列准备地址
		colValues := make([]interface{}, len(cols))
		for i, col := range cols {
			if fieldIdx, ok := colToFieldIdx[col]; ok {
				colValues[i] = elem.Field(fieldIdx).Addr().Interface()
			} else {
				var tmp interface{}
				colValues[i] = &tmp
			}
		}
		if err := rows.Scan(colValues...); err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, elem))
	}
	return rows.Err()
}

// ================== 执行命令方法 ==================

// Exec 执行 INSERT/UPDATE/DELETE，返回 sql.Result
func (m *MySQLUtil) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.db.Exec(query, args...)
}

// Insert 插入并返回自增ID
func (m *MySQLUtil) Insert(query string, args ...interface{}) (int64, error) {
	result, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 执行更新，返回受影响行数
func (m *MySQLUtil) Update(query string, args ...interface{}) (int64, error) {
	result, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ================== 事务支持 ==================

// Transaction 执行事务，回调函数中执行数据库操作，自动提交/回滚
func (m *MySQLUtil) Transaction(fn func(tx *sql.Tx) error) error {
	tx, err := m.db.Begin()
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

// BatchInsert 批量插入数据（适用于同一表的多行插入）
// table: 表名
// columns: 列名切片，如 []string{"name", "age"}
// rows: 每行数据的切片，每行是与 columns 对应的 []interface{}
func (m *MySQLUtil) BatchInsert(table string, columns []string, rows [][]interface{}) error {
	if len(rows) == 0 {
		return nil
	}
	// 构建占位符 (?, ?...)
	placeholders := make([]byte, 0, len(columns)*2)
	for i := 0; i < len(columns); i++ {
		if i > 0 {
			placeholders = append(placeholders, ',')
		}
		placeholders = append(placeholders, '?')
	}
	phStr := string(placeholders)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, strings.Join(columns, ","), phStr)
	// 逐行执行（也可拼接多行，但注意参数个数）
	for _, row := range rows {
		if _, err := m.db.Exec(query, row...); err != nil {
			return err
		}
	}
	return nil
}

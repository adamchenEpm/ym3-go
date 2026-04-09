package main_test

import (
	"github.com/adamchenEpm/ym3-go/internal/config"
	"github.com/adamchenEpm/ym3-go/internal/mysql"
	"testing"
	"time"
)

/*
 * 测试 config.NewConfig
 */
func Test_config_NewConfig(t *testing.T) {

	cfg := config.Get()
	//t.Assert(cfg != nil)

	t.Logf("Config.name: %v,  code :%v", cfg.Name, cfg.Code)
}

/*
 * 测试 mysql.QueryToStructs
 */
func Test_mysql_QueryToStructs(t *testing.T) {
	db := mysql.GetInstance()
	defer db.Close()

	// 4. 查询并扫描到结构体
	type User struct {
		ID         int64      `db:"id"`
		Name       string     `db:"name"`
		UpdateTime *time.Time `db:"update_time"`
	}
	var users []User
	err := db.QueryToStructs("SELECT id, name,update_time FROM sys_user WHERE id = ?", &users, 138)
	if err != nil {
		t.Fatalf("QueryToStructs失败: %v", err)
	}
	if len(users) < 1 {
		t.Logf("查询结果是空的 ")
	} else {
		t.Logf("查询结果正确: %+v", users)
		if users[0].UpdateTime != nil {
			t.Logf("格式化 update_time: %v", users[0].UpdateTime.Format("2006-01-02 15:04:05"))
		}
	}

}

// TestIntegration_BatchInsert 测试批量插入
func TestIntegration_BatchInsert(t *testing.T) {
	db := mysql.GetInstance()
	defer db.Close()

	rows := [][]interface{}{
		{"John", 22},
		{"Jane", 27},
	}
	err := db.BatchInsert("users", []string{"name", "age"}, rows)
	if err != nil {
		t.Fatalf("BatchInsert失败: %v", err)
	}
}

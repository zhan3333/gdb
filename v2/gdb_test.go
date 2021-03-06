package gdb_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhan3333/gdb/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestMain(m *testing.M) {
	maxLiftTime := 30 * time.Second
	parseTime := true
	gdb.DefaultName = "default"
	gdb.ConnConfigs = map[string]gdb.MySQLConf{
		gdb.DefaultName: {
			Host:        "127.0.0.1",
			Port:        "3306",
			Username:    "root",
			Password:    "123456",
			Database:    "test",
			ParseTime:   &parseTime,
			MaxLiftTime: &maxLiftTime,
			GORMConfig: &gorm.Config{
				Logger: logger.New(
					log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
					logger.Config{
						SlowThreshold: time.Second, // 慢 SQL 阈值
						LogLevel:      logger.Info, // Log level
						Colorful:      true,        // 禁用彩色打印
					},
				),
			},
		},
	}
	m.Run()
}

func TestConn(t *testing.T) {
	var err error
	_, err = gdb.InitDef()
	assert.Nil(t, err)
	err = gdb.DB.SQLDB.Ping()
	assert.Nil(t, err)
	err = gdb.Def().SQLDB.Ping()
	assert.Nil(t, err)
}

func TestQuery(t *testing.T) {
	var err error
	_, err = gdb.InitDef()
	assert.Nil(t, err)
	assert.Nil(t, gdb.Def().SQLDB.Ping())
}

func TestMigrate(t *testing.T) {
	var err error
	err = gdb.Def().AutoMigrate(&User{})
	assert.Nil(t, err)
}

func TestGetTables(t *testing.T) {
	var err error
	err = gdb.Def().AutoMigrate(&User{})
	assert.Nil(t, err)
	var tables []string
	err = gdb.Def().Raw("show tables").Scan(&tables).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "users", tables[0])
}

func TestInsertMany(t *testing.T) {
	var err error
	err = gdb.Def().AutoMigrate(&User{})
	assert.Nil(t, err)
	var usersCount int64
	err = gdb.Def().Table("users").Count(&usersCount).Error
	assert.Nil(t, err)
	var users = []User{
		{
			Name:     "a",
			Email:    "a",
			Password: "a",
		}, {
			Name:     "a",
			Email:    "a",
			Password: "a",
		},
	}
	err = gdb.Def().CreateInBatches(&users, 1).Error
	assert.Nil(t, err)
	var newUsersCount int64
	err = gdb.Def().Table("users").Count(&newUsersCount).Error
	assert.Equal(t, newUsersCount, usersCount+2)
}

// 测试批量操作
func TestChunk(t *testing.T) {
	var results []User
	result := gdb.Def().FindInBatches(&results, 2, func(tx *gorm.DB, batch int) error {
		for _, r := range results {
			t.Logf("%+v \n", r)
		}
		return nil
	})
	assert.Nil(t, result.Error)
	t.Log(result.RowsAffected)
	t.Log(results)
}

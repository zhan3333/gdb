package gdb_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhan3333/gdb"
	"testing"
)

func TestMain(m *testing.M) {
	gdb.ConnConfigs = map[string]gdb.MysqlConf{
		"default": {
			Host:        "127.0.0.1",
			Port:        "3306",
			Username:    "root",
			Password:    "123456",
			Database:    "test",
			MaxLiftTime: 30,
			LogMode:     true,
		},
	}
	m.Run()
	gdb.Close()
}

func TestConn(t *testing.T) {
	var err error
	_, err = gdb.InitDef()
	assert.Nil(t, err)
	err = gdb.Def().DB().Ping()
	assert.Nil(t, err)
}

func TestQuery(t *testing.T) {
	var err error
	_, err = gdb.InitDef()
	assert.Nil(t, err)
	assert.Nil(t, gdb.Def().DB().Ping())
}

func TestGetTables(t *testing.T) {
	var err error
	tables := []string{}
	err = gdb.Def().Raw("show tables").Pluck("Tables_in_mysql", &tables).Error
	assert.Nil(t, err)
	t.Log(tables)
}

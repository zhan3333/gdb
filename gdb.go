package gdb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"time"
)

var (
	DefaultName = "default"
	ConnConfigs map[string]MysqlConf
)

type MysqlConf struct {
	Host        string
	Port        string
	Username    string
	Password    string
	Database    string
	Loc         string
	ParseTime   *string
	Timeout     *uint
	Charset     string
	MaxLiftTime time.Duration
	LogMode     bool
	Log         *log.Logger
}

func (c MysqlConf) String() string {
	charset := "utf8"
	parseTime := "True"
	loc := "Local"
	timeout := uint(15)
	if c.Charset != "" {
		charset = c.Charset
	}
	if c.ParseTime != nil {
		parseTime = *c.ParseTime
	}
	if c.Loc != "" {
		loc = c.Loc
	}
	if c.Timeout != nil {
		timeout = *c.Timeout
	}
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s&timeout=%ds",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		charset,
		parseTime,
		loc,
		timeout,
	)
}

var connections = map[string]*gorm.DB{}

func InitDef() (*gorm.DB, error) {
	db, err := InitConn(ConnConfigs["default"])
	if err != nil {
		return db, err
	}
	connections["default"] = db
	return db, nil
}

func InitConn(c MysqlConf) (*gorm.DB, error) {
	var conn *gorm.DB
	var err error
	conn, err = gorm.Open("mysql", c.String())
	if err != nil {
		return conn, err
	}
	conn.LogMode(c.LogMode)
	if c.Log != nil {
		conn.SetLogger(c.Log)
	} else {
		conn.SetLogger(gorm.Logger{log.New(os.Stdout, "\r\n", 0)})
	}
	conn.DB().SetConnMaxLifetime(c.MaxLiftTime)
	return conn, nil
}

func Close() {
	for k, conn := range connections {
		if err := conn.Close(); err != nil {
			log.Printf("Close mysql conn %s err: %+v", k, err)
		}
	}
}

// get default conn
func Def() *gorm.DB {
	return Conn(DefaultName)
}

// 获取指定的连接, 当创建新连接失败时, 将会返回默认连接
func Conn(name string) *gorm.DB {
	var err error
	if conn, ok := connections[name]; ok {
		return conn
	}
	if c, ok := ConnConfigs[name]; ok {
		connections[name], err = InitConn(c)
		if err != nil {
			panic(fmt.Sprintf("Connect mysql (%s: %s) failed: %+v", name, c.String(), err))
			return nil
		}
		return connections[name]
	}
	panic(fmt.Sprintf("Can't read mysql config: %s", name))
	return nil
}

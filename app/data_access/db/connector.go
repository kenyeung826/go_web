package db

import (
	"app/util"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

type ConnMap struct {
	mu sync.Mutex
	v  map[string]*sql.DB
}

var (
	cm     *ConnMap
	cmOnce sync.Once
)

func (c *ConnMap) newConnection(name string) *sql.DB {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := c.v[name]; ok {
		return conn
	} else {
		config := &mysql.Config{
			User:                 os.Getenv(fmt.Sprintf("DB_%s_USER", name)), // Username
			Passwd:               os.Getenv(fmt.Sprintf("DB_%s_PASS", name)), // Password (requires User)
			Net:                  os.Getenv(fmt.Sprintf("DB_%s_NET", name)),  // Network (e.g. "tcp", "tcp6", "unix". default: "tcp")
			Addr:                 os.Getenv(fmt.Sprintf("DB_%s_HOST", name)), // Address (default: "127.0.0.1:3306" for "tcp" and "/tmp/mysql.sock" for "unix")
			DBName:               os.Getenv(fmt.Sprintf("DB_%s_DBNAME", name)),
			AllowNativePasswords: true,
		}

		dConn, err := mysql.NewConnector(config)
		util.CheckError(err, nil)

		conn := sql.OpenDB(dConn)
		conn.SetConnMaxLifetime(time.Minute * 3)
		conn.SetMaxOpenConns(10)
		conn.SetMaxIdleConns(10)

		c.v[name] = conn
		return cm.v[name]
	}
}

func (c *ConnMap) Close(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if conn, ok := c.v[name]; ok {
		delete(c.v, name)
		defer func() {
			_ = conn.Close()
		}()

	}
}

func CloseAll() {
	if cm != nil {
		for name := range cm.v {
			cm.Close(name)
		}
	}
}

func GetConn(name string) *sql.DB {
	return cm.newConnection(name)
}

func init() {
	cm = &ConnMap{v: make(map[string]*sql.DB)}
}

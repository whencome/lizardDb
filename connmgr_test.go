package lizardDb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"sync"
)

var db *sql.DB
var dsn = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8"
var myConnMgr *MyConnManager

func init() {
	myConnMgr = NewMyConnManager()
	// 注册连接管理器
	RegisterConnectionManager(myConnMgr)
}

// 定义一个ConnectionManager
type MyConnManager struct {
	Connections 		map[string]*sql.DB
	mutex 				sync.RWMutex
}

func NewMyConnManager() *MyConnManager {
	return &MyConnManager{
		Connections : make(map[string]*sql.DB),
	}
}

// 获取读连接
func (mcm *MyConnManager) GetReadConn(dbName string) (*sql.DB, error) {
	cfgKey := fmt.Sprintf("%s_read", dbName)
	if conn, ok := mcm.Connections[cfgKey]; ok {
		return conn, nil
	}
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	mcm.mutex.Lock()
	mcm.Connections[cfgKey] = conn
	mcm.mutex.Unlock()
	return conn, nil
}

// 获取写连接
func (mcm *MyConnManager) GetWriteConn(dbName string) (*sql.DB, error) {
	cfgKey := fmt.Sprintf("%s_write", dbName)
	if conn, ok := mcm.Connections[cfgKey]; ok {
		return conn, nil
	}
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	mcm.mutex.Lock()
	mcm.Connections[cfgKey] = conn
	mcm.mutex.Unlock()
	return conn, nil
}
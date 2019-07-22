package lizardDb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"sync"
	"testing"
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

// 测试查询单个值
func TestSqlExecutor_QueryRawScalar(t *testing.T) {
	query := "select count(0) from book"
	executor := NewSqlExecutor(false)
	v, err := executor.QueryRawScalar("test", query)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("RESULT: ", v)
}

// 测试查询一条记录
func TestSqlExecutor_QueryRawRow(t *testing.T) {
	query := "select * from book where id = ?"
	executor := NewSqlExecutor(false)
	result, err := executor.QueryRawRow("test", query, 2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("RESULT: %v\n", result)
	for k, v := range result {
		fmt.Printf("%s : %s\n", k, v)
	}
}

func TestSqlExecutor_QueryRawRows(t *testing.T) {
	query := "select * from book"
	executor := NewSqlExecutor(false)
	result, err := executor.QueryRawRows("test", query)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("RESULT: %v\n", result)
	for _, book := range result {
		for k, v := range book {
			fmt.Printf("%s : %s\n", k, v)
		}
		fmt.Println("---------------")
	}
}


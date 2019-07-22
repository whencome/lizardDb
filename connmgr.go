package lizardDb

import (
	"database/sql"
)

// 连接管理对象
type ConnectionManager interface {
	// 获取读连接
	GetReadConn(dbName string) (*sql.DB, error)
	// 获取写连接
	GetWriteConn(dbName string) (*sql.DB, error)
}
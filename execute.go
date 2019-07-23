package lizardDb

import (
	"database/sql"
	"regexp"
	"strings"
	"errors"
	"github.com/whencome/lizardDb/types"
)

// 数据库执行对象
type SqlExecutor struct {
	UseReadConn   			bool		// 是否使用读连接
}

// 创建一个新的sql执行对象
// useReadConn ： 是否使用读连接，为false是使用写连接
func NewSqlExecutor(useReadConn bool) *SqlExecutor {
	return &SqlExecutor{
		UseReadConn : useReadConn,
	}
}

// 使用读库（默认）
func (se *SqlExecutor) Read() {
	se.UseReadConn = true
}

// 使用写库
func (se *SqlExecutor) Write() {
	se.UseReadConn = false
}

// connect to the database which the dbName mapped to
func (se *SqlExecutor) Connect(dbName string) (*sql.DB, error) {
	// check whether is a connection manager registered or not
	if DbManager == nil {
		return nil, errors.New("new connection manager registered")
	}
	// get a connection
	if se.UseReadConn {
		return DbManager.GetReadConn(dbName)
	}
	return DbManager.GetWriteConn(dbName)
}

// 查询单个值
func (se *SqlExecutor) QueryRawScalar(dbName string, query string, args ...interface{}) (*types.Value, error) {
	// 获取连接
	conn, err := se.Connect(dbName)
	if err != nil {
		return nil, err
	}
	// 查询
	row := conn.QueryRow(query, args...)
	var ret interface{}
	err = row.Scan(&ret)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return types.NewValue(ret), nil
}

// 查询单行数据
func (se *SqlExecutor) QueryRawRow(dbName string, query string, args ...interface{}) (map[string]*types.Value, error) {
	// 获取连接
	conn, err := se.Connect(dbName)
	if err != nil {
		return nil, err
	}
	// 只查询一条数据
	matched, err := regexp.Match(`\s+limit\s+(\d+\s*,\s*)?\d+$`, []byte(strings.ToLower(strings.TrimSpace(query))))
	if err != nil {
		return nil, err
	}
	if !matched {
		query += " LIMIT 1"
	}
	// 使用查询多条数据的方式查询
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	// 获取查询的字段列表
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	fieldsCount := len(fields)
	var result= make(map[string]*types.Value)
	for rows.Next() {
		var data= make([]interface{}, fieldsCount)
		for i := 0; i < fieldsCount; i++ {
			var tmp interface{}
			data[i] = &tmp
		}
		err = rows.Scan(data...)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		// 组装数据
		for idx, field := range fields {
			value := data[idx]
			result[field] = types.NewValue(*value.(*interface{}))
		}
		// 只读取一条
		break
	}
	return result, nil
}

// 查询多行数据
func (se *SqlExecutor) QueryRawRows(dbName string, query string, args ...interface{}) ([]map[string]*types.Value, error) {
	// 获取连接
	conn, err := se.Connect(dbName)
	if err != nil {
		return nil, err
	}
	// 使用查询多条数据的方式查询
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	// 获取查询的字段列表
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	fieldsCount := len(fields)
	var result = make([]map[string]*types.Value, 0)
	// 读取数据
	for rows.Next() {
		var data= make([]interface{}, fieldsCount)
		for i := 0; i < fieldsCount; i++ {
			var tmp interface{}
			data[i] = &tmp
		}
		err = rows.Scan(data...)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		// 组装数据
		var rowMapData = make(map[string]*types.Value)
		for idx, field := range fields  {
			value := data[idx]
			rowMapData[field] = types.NewValue(*value.(*interface{}))
		}
		result = append(result, rowMapData)
	}
	// 如果没有数据则设置为nil
	if len(result) == 0 {
		result = nil
	}
	return result, nil
}
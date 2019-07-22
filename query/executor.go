package query

import (
	"database/sql"
)

type Executor struct {
	Conn 			sql.DB	// 连接管理器
	Querier 		Querier 	// 查询对象
	UseReadConn   	bool		// 是否使用读连接
}

func NewExecutor(q Querier, read bool) *Executor {
	return &Executor{
		Querier : q,
		UseReadConn : read,
	}
}

func (e *Executor) QueryOne() (map[string]string, error) {
	// 获取查询
	query, err := e.Querier.GetQuery()
	if err != nil {
		return nil, err
	}
	// 执行
	stmt, err := e.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	// 开始查询
	queryParams := e.Querier.GetParams().Params()
	row := stmt.QueryRow(queryParams...)
	result := make(map[string]string, 0)
	err = row.Scan(result)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
/*
func (e *Executor) QueryAll() interface{} {}

func (e *Executor) QueryPage() interface{} {}

func (e *Executor) QueryScalar() interface{} {}

func (e *Executor) Execute() interface{} {}
*/

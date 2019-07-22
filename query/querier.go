package query

import (
	"fmt"
	"bytes"
	"strings"
	"errors"
)

// 查询对象
type Querier interface {
	GetDbName() string
	GetCommand() string
	GetQuery() (string, error)
	GetFields() []string
	GetParams() *QueryParams
}

// 原始查询对象
type RawQuerier struct {
	DbName 			string   		// 注册的数据库别名，用户获取连接管理器
	command 		string
	Sql 			string   		// 查询的SQL
	params 			*QueryParams 	// 查询参数列表
}

func NewRawQuerier(db, sql string, params ...interface{}) *RawQuerier {
	sql = strings.TrimSpace(sql)
	firstSpacePos := strings.Index(sql, " ")
	cmd := strings.ToUpper(sql[0:firstSpacePos])
	queryParams := NewQueryParams()
	queryParams.AddBatch(params)
	return &RawQuerier{
		DbName	: db,
		command : cmd,
		Sql 	: sql,
		params 	: queryParams,
	}

}

func (this *RawQuerier) GetDbName() string {
	return this.DbName
}

func (this *RawQuerier) GetCommand() string {
	return this.command
}

func (this *RawQuerier) GetQuery() (string, error) {
	return this.Sql, nil
}

func (this *RawQuerier) GetParams() *QueryParams {
	return this.params
}


// 查询对象
type SimpleQuerier struct {
	DbName 			string   		// 注册的数据库别名，用户获取连接管理器
	command   		string 			// 定义执行的命令
	tableName 		string	 		// 数据表名（真实表名）
	fields 			string   		// 查询的字段列表
	conditions      *Conditions     // 查询条件
	orderBy 	    string   		// 排序方式
	groupBy 		string   		// 分组方式
	having			*Conditions     // 分组过滤条件
	limit 			string 	 		// limit参数
	params 			*QueryParams   	// 查询参数列表
	joinTables 		*JoinTables   	// 连表信息
	updateData 		map[string]interface{}	// 更新数据，用于插入或者更新时
}

func NewSimpleQuerier(db string) *SimpleQuerier {
	return &SimpleQuerier{
		DbName : db,
		joinTables : NewJoinTables(),
		params : NewQueryParams(),
	}
}

func (this *SimpleQuerier) Select(fields string) *SimpleQuerier {
	this.command = SqlCmdSelect
	this.fields = fields
	return this
}

func (this *SimpleQuerier) Update(table string, updateData map[string]interface{}, updateCond *Conditions) {
	this.command = SqlCmdUpdate
	this.updateData = updateData
	this.conditions = updateCond
}

func (this *SimpleQuerier) Insert(table string, insertData map[string]interface{}) {
	this.command = SqlCmdInsert
	this.updateData = insertData
}

func (this *SimpleQuerier) Delete(table string, deleteCond *Conditions) {
	this.command = SqlCmdDelete
	this.conditions = deleteCond
}

func (this *SimpleQuerier) From(tableName string) *SimpleQuerier {
	this.tableName = tableName
	return this
}

func (this *SimpleQuerier) Where(cond *Conditions, args... interface{}) *SimpleQuerier {
	this.conditions = cond
	this.params.AddBatch(args)
	return this
}

func (this *SimpleQuerier) GroupBy(groupBy string) *SimpleQuerier {
	this.groupBy = groupBy
	return this
}

func (this *SimpleQuerier) Having(cond *Conditions, args... interface{}) *SimpleQuerier {
	this.having = cond
	this.params.AddBatch(args)
	return this
}

func (this *SimpleQuerier) OrderBy(orderBy string) *SimpleQuerier {
	this.orderBy = orderBy
	return this
}

func (this *SimpleQuerier) Limit(offset, size int) *SimpleQuerier {
	this.limit = fmt.Sprintf("%d, %d", offset, size)
	return this
}

// 连表查询
func (this *SimpleQuerier) Join(joinTable *JoinTable, args... interface{}) *SimpleQuerier {
	this.joinTables.Add(joinTable)
	this.params.AddBatch(args)
	return this
}

func (this *SimpleQuerier) buildCommonQuery() (string, error) {
	var buf bytes.Buffer
	// table
	buf.WriteString("FROM ")
	buf.WriteString(this.tableName)
	// join info
	joinStr, err := this.joinTables.Build()
	if err != nil {
		return "", err
	}
	buf.WriteString(joinStr)
	buf.WriteString(" ")
	// conditions
	if !this.conditions.IsEmpty() {
		condStr, err := this.conditions.Build()
		if err != nil {
			return "", err
		}
		buf.WriteString(" WHERE ")
		buf.WriteString(condStr)
		buf.WriteString(" ")
	}
	// group info
	if this.groupBy != "" {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(this.groupBy)
		buf.WriteString(" ")
		// filter info
		if !this.having.IsEmpty() {
			havingCond, err := this.having.Build()
			if err != nil {
				return "", err
			}
			buf.WriteString(" HAVING ")
			buf.WriteString(havingCond)
			buf.WriteString(" ")
		}
	}
	return buf.String(), nil
}

func (this *SimpleQuerier) buildQuery() (string, error) {
	var buf bytes.Buffer
	// command
	buf.WriteString("SELECT ")
	// fields
	if this.fields == "" {
		buf.WriteString("*")
	} else {
		buf.WriteString(this.fields)
	}
	buf.WriteString(" ")

	// common query
	commonQuery, err := this.buildCommonQuery()
	if err != nil {
		return "", err
	}
	buf.WriteString(commonQuery)

	// order info
	if this.orderBy != "" {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(this.orderBy)
		buf.WriteString(" ")
	}
	// limit info
	if this.limit != "" {
		buf.WriteString(" LIMIT ")
		buf.WriteString(this.limit)
	}
	// not finished yet
	return buf.String(), nil
}

func (this *SimpleQuerier) buildCountQuery() (string, error) {
	var buf bytes.Buffer
	// command
	buf.WriteString("SELECT COUNT(0) AS num ")
	// common query
	commonQuery, err := this.buildCommonQuery()
	if err != nil {
		return "", err
	}
	buf.WriteString(commonQuery)
	// not finished yet
	return buf.String(), nil
}

func (this *SimpleQuerier) buildInsert() (string, error) {
	var buf bytes.Buffer
	buf.WriteString("INSERT INTO ")
	buf.WriteString(this.tableName)
	fields := make([]string, 0)
	values := make([]string, 0)
	for k, v := range this.updateData {
		fields = append(fields, k)
		this.params.Add(v)
		values = append(values, "?")
	}
	buf.WriteString("(")
	buf.WriteString(strings.Join(fields, ","))
	buf.WriteString(") ")
	buf.WriteString(" VALUES(")
	buf.WriteString(strings.Join(values, ","))
	buf.WriteString(") ")
	return buf.String(), nil
}

func (this *SimpleQuerier) buildUpdate() (string, error) {
	var buf bytes.Buffer
	buf.WriteString("UPDATE ")
	buf.WriteString(this.tableName)
	buf.WriteString(" SET ")
	count := 0
	for field, value := range this.updateData {
		if count > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(field)
		buf.WriteString(" = ? ")
		this.params.Add(value)
	}
	upCond, err := this.conditions.Build()
	if err != nil {
		return "", nil
	}
	buf.WriteString(" WHERE ")
	buf.WriteString(upCond)
	return buf.String(), nil
}

func (this *SimpleQuerier) buildDelete() (string, error) {
	delCond, err := this.conditions.Build()
	if err != nil {
		return "", nil
	}
	var buf bytes.Buffer
	buf.WriteString("DELETE FROM ")
	buf.WriteString(this.tableName)
	buf.WriteString(" WHERE ")
	buf.WriteString(delCond)
	return buf.String(), nil
}



//////////// 接口定义方法 //////////////

func (this *SimpleQuerier) GetDbName() string {
	return this.DbName
}

func (this *SimpleQuerier) GetCommand() string {
	return this.command
}

func (this *SimpleQuerier) GetQuery() (string, error) {
	// select查询
	if this.command == SqlCmdSelect {
		return this.buildQuery()
	}
	// insert、update、delete
	switch this.command {
	case SqlCmdSelect:
		return this.buildQuery()
	case SqlCmdInsert:
		return this.buildInsert()
	case SqlCmdUpdate:
		return this.buildUpdate()
	case SqlCmdDelete:
		return this.buildDelete()
	default:
		return "", errors.New("sql command not supported")
	}
}

func (this *SimpleQuerier) GetParams() *QueryParams {
	return this.params
}
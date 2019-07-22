package query

import (
	"testing"
)

var dbName = "test_db"
var tableName = "test_table"

func TestSimpleQuerier_Select(t *testing.T) {
	// 构造查询条件
	queryConds := NewConditions()
	conds1 := NewConditions()
	conds2 := NewConditions()
	err := conds1.AndPlain("id > ?")
	if err != nil {
		t.Error(err)
	}
	err = conds1.And(NewPlainCondition("score > ?"))
	if err != nil {
		t.Error(err)
	}
	err = conds1.And(NewFieldCondition("uid", "in", []int{100001,100002,100003}))
	if err != nil {
		t.Error(err)
	}
	err = conds2.And(NewFieldCondition("class", "not in", []string{"classA","classB","classC"}))
	if err != nil {
		t.Error(err)
	}
	queryConds.And(conds1)
	queryConds.Or(conds2)

	// 构造分组条件
	havingConds := NewConditions()
	havingConds.And(NewFieldCondition("fieldC", "between", []int{70, 90}))

	// 构造查询对象
	querier := NewSimpleQuerier(dbName)
	querier.Select("fieldA,fieldB,fieldC,fieldD").From(tableName).Where(queryConds, 100, 60)
	querier.GroupBy("fieldA").Having(havingConds)
	querier.OrderBy("fieldD desc").Limit(0, 20)

	query, err := querier.GetQuery()
	if err != nil {
		t.Error(err)
	}
	t.Log(query)
}
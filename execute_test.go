package lizardDb

import (
	"fmt"
	"testing"
)

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


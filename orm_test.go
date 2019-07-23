package lizardDb

/**
 * 测试命令：go test -v defines.go conf.go lizard_db.go connmgr.go connmgr_test.go execute.go orm.go orm_manager.go orm_test.go
 */

import (
	"testing"
	"fmt"
)

type Book struct {
	DatabaseName 	DatabaseName 		`orm:"test"`
	TableName 		TableName 			`orm:"book"`
	Id 				int 				`orm:"id"`
	Name 			string 				`orm:"name"`
	Author 			string 				`orm:"author"`
	Price 			float64 			`orm:"price"`
}

func (b *Book) New() DbObject {
	return &Book{}
}

func (b *Book) GetDbName() string {
	return "test"
}

// 获取数据表名称
func (b *Book) GetTableName() string {
	return "book"
}

func TestORM_FetchOneOnRaw(t *testing.T) {
	var book *Book = &Book{}
	querySql := "select * from book where id = ?"
	om := NewObjectManager(book)
	data, err := om.Read().FetchOne(querySql, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("RESULT: %#v\n", data)
	fmt.Println("Book Name: ", data.(*Book).Name)
}

func TestORM_FetchAllOnRaw(t *testing.T) {
	var book *Book = &Book{}
	querySql := "select * from book"
	om := NewObjectManager(book)
	books, err := om.Read().FetchAll(querySql)
	if err != nil {
		t.Fatal(err)
	}
	if books == nil {
		fmt.Println("no book data")
	}
	fmt.Printf("RESULT: %#v\n", books)
	for _, data := range books {
		tmpBook := data.(*Book)
		fmt.Println("ID: ", tmpBook.Id)
		fmt.Println("Name: ", tmpBook.Name)
		fmt.Println("Author: ", tmpBook.Author)
		fmt.Println("Price: ", tmpBook.Price)
		fmt.Println("-------------------------------")
	}
}


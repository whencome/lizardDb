package lizardDb

import (
	"strings"
)

// 数据库
type DatabaseName string
// 表名
type TableName string

// 元数据
type MetaData struct {
	// 存储原始数据
	RawData					map[string]interface{}
	// 表字段和struct字段的映射关系
	FieldPropertyMappings	map[string]string
	PropertyFieldMappings	map[string]string
	PropertyOrmTags			map[string]*OrmTag
	// 数据库
	DatabaseName 			string
	// 数据表
	TableName 				string
}

// orm基类
type Object struct {
	// Meta 					*MetaData
}

func (o Object) GetDbName() string {
	return ""
}

// 获取数据表名称
func (o Object) GetTableName() string {
	return ""
}


type OrmTag struct {
	content 				string
	FieldName 				string	// 数据库字段名称
}

func NewOrmTag(tag string) *OrmTag{
	ot := &OrmTag{
		content : tag,
	}
	ot.Resolve()
	return ot
}

func (ot *OrmTag) Resolve() {
	if ot.content == "" {
		return
	}
	// 解析tag内容
	vals := strings.Split(ot.content, ",")
	// 设置字段名称
	ot.FieldName = vals[0]
}

func (ot *OrmTag) getFieldName() string {
	return ot.FieldName
}

func NewMetaData() *MetaData {
	return &MetaData{
		RawData 				: make(map[string]interface{}),
		FieldPropertyMappings 	: make(map[string]string),
		PropertyFieldMappings 	: make(map[string]string),
		PropertyOrmTags 		: make(map[string]*OrmTag),
		DatabaseName 			: "",
		TableName 				: "",
	}
}
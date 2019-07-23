package lizardDb

import (
	"github.com/whencome/lizardDb/query"
	"github.com/whencome/lizardDb/types"
	"reflect"
)

type ObjectManager struct {
	Meta 					*MetaData
	Obj 					DbObject
	UseReadConn   			bool		// 是否使用读连接
	Querier					*query.Querier	// 查询对象
}

// 创建一个新的ObjectManager
func NewObjectManager(o DbObject) *ObjectManager {
	om := &ObjectManager{
		UseReadConn : true,
		Meta 		: NewMetaData(),
		Obj 		: o,
	}
	om.Resolve(o)
	return om
}

// 使用读库（默认）
func (om *ObjectManager) Read() *ObjectManager {
	om.UseReadConn = true
	return om
}

// 使用写库
func (om *ObjectManager) Write() *ObjectManager {
	om.UseReadConn = false
	return om
}

// 解析对象
func (om *ObjectManager) Resolve(obj DbObject) {
	om.Meta.DatabaseName = obj.GetDbName()
	om.Meta.TableName = obj.GetTableName()
	rt := reflect.TypeOf(obj).Elem()
	fieldsNum := rt.NumField()
	for i := 0; i < fieldsNum; i++ {
		propertyName := rt.Field(i).Name
		tag := rt.Field(i).Tag.Get(TagName)
		ormTag := NewOrmTag(tag)
		// 获取数据库和表名
		if rt.Field(i).Type.Name() == "DatabaseName" {}
		if rt.Field(i).Type.Name() == "TableName" {}
		// 未设置字段名表示不在db字段映射之中，该属性为独立的字段
		if ormTag.FieldName == "" {
			continue
		}
		// 设置值
		om.Meta.PropertyFieldMappings[propertyName] = ormTag.FieldName
		om.Meta.FieldPropertyMappings[ormTag.FieldName] = propertyName
		om.Meta.PropertyOrmTags[propertyName] = ormTag
	}
}

// 获取字段对应的属性名称
func (om *ObjectManager) GetFieldPropertyName(field string) string {
	property, ok := om.Meta.FieldPropertyMappings[field]
	if !ok {
		return ""
	}
	return property
}

// 获取数据库名称
func (om *ObjectManager) GetDbName() string {
	return om.Meta.DatabaseName
}

// 获取数据库名称
func (om *ObjectManager) GetTableName() string {
	return om.Meta.TableName
}

// 根据数据表字段设置属性的值
// 如果不存在则忽略
func (om *ObjectManager) Bind(obj interface{}, row map[string]*types.Value) {
	if row == nil {
		return
	}
	rv := reflect.ValueOf(obj).Elem()
	for field, val := range row {
		propertyName := om.GetFieldPropertyName(field)
		if propertyName == "" {
			continue
		}
		prop := rv.FieldByName(propertyName)
		if !prop.IsValid() {
			continue
		}
		switch prop.Type().Kind() {
		case reflect.String:
			prop.SetString(val.String())
		case reflect.Int: fallthrough
		case reflect.Int8: fallthrough
		case reflect.Int16: fallthrough
		case reflect.Int32: fallthrough
		case reflect.Int64:
			prop.SetInt(val.Int())
		case reflect.Uint: fallthrough
		case reflect.Uint8: fallthrough
		case reflect.Uint16: fallthrough
		case reflect.Uint32: fallthrough
		case reflect.Uint64:
			prop.SetUint(uint64(val.Uint()))
		case reflect.Bool:
			prop.SetBool(val.Bool())
		case reflect.Float32:
		case reflect.Float64:
			prop.SetFloat(val.Float())
		}
	}
}

func (om *ObjectManager) BindBatch(result []interface{}, rows []map[string]*types.Value) {
	if rows == nil || len(rows) == 0 {
		return
	}
	for _, row := range rows {
		data := om.Obj
		om.Bind(data, row)
		result = append(result, data)
	}
}

// 根据原始的sql，查询一条数据
func (om *ObjectManager) FetchSimple(querySql string, params ...interface{}) (*types.Value, error) {
	// 查询
	executor := NewSqlExecutor(om.UseReadConn)
	result, err := executor.QueryRawScalar(om.GetDbName(), querySql, params...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 根据原始的sql，查询一条数据
func (om *ObjectManager) FetchOne(querySql string, params ...interface{}) (interface{}, error) {
	// 查询
	executor := NewSqlExecutor(om.UseReadConn)
	rowData, err := executor.QueryRawRow(om.GetDbName(), querySql, params...)
	if err != nil {
		return nil, err
	}
	result := om.Obj.New()
	// 绑定值
	om.Bind(result, rowData)
	return result, nil
}

// 根据原始的sql，查询所有数据
func (om *ObjectManager) FetchAll(querySql string, params ...interface{}) ([]interface{}, error) {
	// 查询
	executor := NewSqlExecutor(om.UseReadConn)
	rowsData, err := executor.QueryRawRows(om.GetDbName(), querySql, params...)
	if err != nil {
		return nil, err
	}
	if len(rowsData) == 0 {
		return nil, nil
	}
	var result = make([]interface{}, 0)
	for _, rowData := range rowsData {
		data := om.Obj.New()
		om.Bind(data, rowData)
		result = append(result, data)
	}
	return result, nil
}

// orm查询

func (om *ObjectManager) Match() *DbObject {
	return nil
}
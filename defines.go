package lizardDb

// 定义数据对象
type DbObject interface {
	// 创建一个新的对象
	New() DbObject
	// 获取数据库名称
	GetDbName() string
	// 获取数据表名称
	GetTableName() string
}

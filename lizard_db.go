package lizardDb

// 全局连接管理器
var DbManager ConnectionManager

// 提供注册方法，用于注册连接管理对象
func RegisterConnectionManager(mgr ConnectionManager) {
	DbManager = mgr
}


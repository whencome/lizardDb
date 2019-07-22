package orm

// 全局参数设置
var TimeFormat = "2006-01-02 15:04:05"
// 设置Tag名称
var TagName = "orm"


// 设置全局时间格式
func SetTimeFormat(format string) {
	TimeFormat = format
}

// 设置Tag名称
func SetTagName(tagName string) {
	TagName = tagName
}
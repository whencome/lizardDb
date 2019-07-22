package orm

import "time"

type User struct {
	Object
	TableName			TableName `orm:"users"`
	Id 			  		int `orm:"id"`
	Name 				string `orm:"name"`
	Password 			string `orm:"password"`
	NickName 			string `orm:"nick_name"`
	Email				string `orm:"email"`
	MobileNo			string `orm:"mobile_no"`
	RegisterTime		time.Time `orm:"register_time"`
	RegisterIP			string `orm:"reg_ip"`
	LastLoginTime 		time.Time `orm:"last_login_time"`
}

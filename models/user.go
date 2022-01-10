package models

// 数据库结构体
// db对应的是数据库的字段
type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Token    string	
}

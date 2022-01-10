package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"web1/models"
)

// 把每一层数据库操作封装成函数
// 等待logic层根据业务需求调用

const secret = "lin"

// CheckUserExist 检查指定的用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlstr := `select count(user_id) from user where username=?`
	// 当用户存在的时候count是大于0的
	var count int
	if err := db.Get(&count, sqlstr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 数据库不能存储明文的密码
	// 1、对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 2、执行SQL语句入库
	sqlstr := "insert into user(user_id,username,password) values(?,?,?)"
	_, err = db.Exec(sqlstr, user.UserID, user.Username, user.Password)
	return err
}

// encryptPassword 给密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login 向数据库中查询记录
func Login(user *models.User) (err error) {
	// 保存用户输入的密码，后面需要进行加密比较
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err != nil {
		// 查询数据库失败
		return ErrorUserNotExist
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserById 根据id获取用户信息
func GetUserById(userid int64) (user *models.User, err error) {
	user = new(models.User)
	// 只向外保留两个字段：user_id和user_name
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, userid)
	return
}

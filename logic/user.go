package logic

import (
	"web1/dao/mysql"
	"web1/models"
	"web1/pkg/jwt"
	"web1/pkg/snowflake"
)

// logic层是用来存放业务逻辑的代码
// SignUp 用来处理注册的业务
func SignUp(p *models.ParamSignUp) (err error) {
	// 1.判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询出错
		return err
	}
	// 2、生成UID
	userID := snowflake.GenID()
	// 3.1、构造
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3、保存进数据库
	return mysql.InsertUser(user)
}

// Login 用来处理登录的业务f
func Login(p *models.ParamLogin) (user *models.User, err error) {
	// 向数据库进行数据操作时，一定要传结构体tag为db的
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 返回mysql.Login中的错误，如果没有则为nil
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 由于mysql.Login(user)中传入的是指针，因此就可以拿到user.UserID
	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}

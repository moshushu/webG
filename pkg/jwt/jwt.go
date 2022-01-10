package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// jwt认证

// token过期时间，2小时

// 加掩
var MySercet = []byte("林小猪")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT Token
// 传入userID和username
// 生成Tkone(string)和一个error
func GenToken(userID int64, username string) (string, error) {
	// 创建一个自己的声明的数据
	c := MyClaims{
		userID,
		"username",
		jwt.StandardClaims{
			// time.Now()当前的时间，加上过期时间TkoenExpireDuration,再将它们转成Unix()
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), // 过期时间
			Issuer: "林小猪", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	// jwt.SigningMethodES256是一个加密算法，对c进行加密
	// 返回一个token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的sercet签名并获得完整的编码后的字符串token
	// 加掩操作
	return token.SignedString(MySercet)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(
		tokenString, // 需要解析的token
		mc,          // 将解析的token放到mc这个MyClaims结构体变量中
		// 生成JWT中有一个加掩的操作，这里是告诉怎么解掩
		func(token *jwt.Token) (i interface{}, err error) {
			return MySercet, nil
		},
	)
	// 解析失败
	if err != nil {
		return nil, err
	}
	// 解析成功
	if token.Valid {
		return mc, nil
	}
	// 无效token
	return nil, errors.New("invalid token")
}

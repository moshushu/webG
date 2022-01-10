package mysql

import (
	"fmt"
	"web1/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

// MySQL 初始化
func Init(cfg *settings.MySQLConfig) (err error) {
	// 1、“user:passwor@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True”
	// 传cfg是为了方便可以在配置文件中控制
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	fmt.Println(dsn)
	// 2、连接数据库。Connect的作用：(连接到数据库并使用ping进行验证。)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	// 3、设置最大连接数
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	// 4、设置最大闲置连接数
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

// 想在mysql包外，关闭db数据库连接，db则必须要大写，但是也可以封装一个包外可见的函数用来
// 关闭db数据库连接，不然则需要把db改为大写，并且会写成：mysql.DB.Close()。
// 没有mysql.Close()简洁
func Close() {
	_ = db.Close()
}

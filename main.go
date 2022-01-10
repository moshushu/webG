package main

import (
	"fmt"
	"web1/dao/mysql"
	"web1/logger"
	"web1/pkg/snowflake"
	"web1/routes"
	"web1/settings"

	"go.uber.org/zap"
)

//搭建比较通用的web脚手架模板

func main() {
	// 1：加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	// 2：初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync() //注册关闭全局logger
	zap.L().Debug("init logger success ....")
	// 3：初始化mysql连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	// 4：初始化Redis连接
	// m: 初始化snowflake，生成分布式ID
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 5：注册路由
	r := routes.Routes(settings.Conf.Mode)
	// 6：启动服务
	r.Run(":8081")
}

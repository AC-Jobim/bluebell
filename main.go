package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/settings"
	"fmt"
	"go.uber.org/zap"
)

// Go Web开发较通用的脚手架模板

func main() {
	// 1.加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("settings初始化配置文件失败了：%v", err)
	}
	// 2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("logger初始化配置文件失败了：%v", err)
	}
	defer zap.L().Sync()
	// 3.初始化MySQL
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("mysql初始化配置文件失败了：%v", err)
	}
	defer mysql.Close()
	// 4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("redis初始化配置文件失败了：%v", err)
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Println("Init() failed, err = ", err)
		return
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	// 5.注册路由
	r := router.SetupRouter(settings.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}

}

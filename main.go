package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

//Goweb开发脚手架通用模板

// @title Huangb's blog接口文档
// @version 1.0
// @description Huangb's blog
// @termsOfService http://swagger.io/terms/

// @contact.name Huangb
// @contact.url http://www.liwenzhou.com
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8081
// @BasePath /api/v1

func main() {
	if len(os.Args) < 2 {
		return
	}
	//1.加载配置文件
	if err := settings.Init(os.Args[1]); err != nil { //传入配置文件
		fmt.Printf("配置文件初始化错误,err:%v", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.AppConfig.Mode); err != nil {
		fmt.Printf("日志初始化错误,err:%v", err)
		return
	}
	defer zap.L().Sync() //用于将缓冲的日志条目写入其输出,确保在程序退出之前刷新日志
	zap.L().Debug("初始化日志成功")
	//3.初始化数据库连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("数据库初始化错误,err:%v", err)
		return
	}
	defer mysql.Close()
	//4.初始化redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("redis初始化错误,err:%v", err)
		return
	}
	defer redis.Close()

	//初始化雪花算法
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("雪花算法初始化错误,err:%v", err)
		return
	}

	//初始化gin框架内置校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("翻译器初始化错误,err:%v\n", err)
		return
	}
	//5.注册路由
	r := routes.Setup(settings.Conf.AppConfig.Mode)
	//6.启动服务（优雅关机）
	srv := &http.Server{
		Addr: fmt.Sprintf("%s",
			viper.GetString("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}

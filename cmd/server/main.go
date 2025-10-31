package main

import (
	"flag"
	"go_blog/internal"
	"go_blog/pkg/config"
	"go_blog/pkg/utils"
	"path/filepath"
	"time"

	"go.uber.org/fx"
)

func main() {
	// Load application configuration.
	var configPath string
	// 命令行参数 -c 指定配置文件路径 使用flag包
	flag.StringVar(&configPath, "c", filepath.Join(utils.GetRootDir(), "conf/config.yaml"), "文件配置地址")
	flag.Parse()
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	//防止onstop方法长时间阻塞 应用优雅退出的超时时间
	// Create a new application container with various components and configurations.
	modules := fx.Options(
		// Supply configuration values to the container.
		fx.Supply(conf),
		// Set a timeout for graceful shutdown of the application.
		fx.StopTimeout(conf.Server.GracefulShutdown+time.Second),
		internal.Module,
		// pkg.Module,
		// crontab.Module,
	)

	if err := fx.ValidateApp(modules); err != nil {
		panic(err)
	}
	app := fx.New(modules)
	// Run the application container.
	app.Run()
}

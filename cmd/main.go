package main

import (
	"os"

	"github.com/WillemCode/AliCloud_Resources/internal/api"
	"github.com/WillemCode/AliCloud_Resources/internal/services"
	"github.com/WillemCode/AliCloud_Resources/pkg/config"
	"github.com/WillemCode/AliCloud_Resources/pkg/database"
	"github.com/WillemCode/AliCloud_Resources/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化日志（使用默认级别，避免配置读取错误时无法记录日志）
	logger.Init("info")

	// 2. 加载配置文件（config.yaml）和环境变量
	cfg, err := config.LoadConfig("")
	if err != nil {
		// 配置加载失败，记录错误并退出
		logger.Log.Fatalf("加载配置时出错: %v", err)
	}
	// fmt.Printf("Config Loaded: %+v ", cfg)
	// 3. 根据配置调整日志级别（如果配置中指定了非默认级别）
	logger.Init(cfg.LogLevel)

	// 4. 初始化数据库连接（SQLite）
	err = database.Init(cfg.Database.Path)
	if err != nil {
		logger.Log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close() // 程序退出时关闭数据库

	// // 5. 遍历配置中的每个阿里云账户，调用相应服务同步数据
	for _, account := range cfg.AliyunAccounts {

		// 同步 ECS 信息
		if len(account.ECSRegionIds) > 0 {
			if err := services.SyncECSInfo(account.Name, account.ECSRegionIds, account.AccessKey, account.AccessSecret); err != nil {
				// 记录错误但不中断，继续处理其他服务
				logger.Log.Errorf("ECS 同步失败 (账户=%s): %v", account.Name, err)
			}
		}
		// 同步 RDS 信息
		if len(account.RDSRegionIds) > 0 {
			// if account.RDSRegionId != "" && account.RDSRegionId != "nil" {
			if err := services.SyncRDSInfo(account.Name, account.RDSRegionIds, account.AccessKey, account.AccessSecret); err != nil {
				logger.Log.Errorf("RDS 同步失败 (账户=%s): %v", account.Name, err)
			}
		}
		// 同步 SLB 信息
		if len(account.SLBRegionIds) > 0 {
			// if account.SLBRegionId != "" && account.SLBRegionId != "nil" {
			if err := services.SyncSLBInfo(account.Name, account.SLBRegionIds, account.AccessKey, account.AccessSecret); err != nil {
				logger.Log.Errorf("SLB 同步失败 (账户=%s): %v", account.Name, err)
			}
		}
		// 同步 PolarDB 信息
		if len(account.PolarDBRegionIds) > 0 {
			// if account.PolarDBRegionId != "" && account.PolarDBRegionId != "nil" {
			if err := services.SyncPolarDBInfo(account.Name, account.PolarDBRegionIds, account.AccessKey, account.AccessSecret); err != nil {
				logger.Log.Errorf("PolarDB 同步失败 (账户=%s): %v", account.Name, err)
			}
		}
	}

	// 6. 启动 Gin Web 服务，提供RESTful查询接口
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		router := gin.Default()

		// 设置API路由
		api.SetupRoutes(router)

		logger.Log.Info("启动 HTTP 服务，监听 :8080 ...")
		_ = router.Run(":8080")
	}
}

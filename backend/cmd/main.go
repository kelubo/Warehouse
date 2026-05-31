package main

import (
	"fmt"
	"os"

	"warehouse-management/backend/cmd/initialize"
	server "warehouse-management/backend/server"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			initialize.Run()
			return
		case "migrate":
			initialize.Migrate()
			return
		case "help":
			printHelp()
			return
		}
	}

	// 默认启动服务
	server.Run()
}

func printHelp() {
	fmt.Println("仓库管理系统")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  warehouse-management [command]")
	fmt.Println()
	fmt.Println("命令:")
	fmt.Println("  init     初始化数据库和配置")
	fmt.Println("  migrate  运行数据库迁移")
	fmt.Println("  help     显示帮助信息")
	fmt.Println("  (空)     启动服务")
	fmt.Println()
}

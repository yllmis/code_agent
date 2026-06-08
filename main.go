package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"go-agents/agent"
	"go-agents/config"
	"go-agents/svc"
	"go-agents/tools"
)

var configFile = flag.String("f", "./etc/config.yaml", "the config file")

func main() {
	// 加载配置
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Println("加载配置失败:", err)
		os.Exit(1)
	}

	// 创建 ServiceContext
	ctx := svc.NewServiceContext(cfg)

	// 注册工具
	toolList := []tools.Tool{
		&tools.SearchTool{},
		&tools.ReadTool{},
		&tools.ListTool{},
	}

	// CLI 交互循环
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("代码助手已启动，输入问题开始对话（输入 exit 退出）")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" {
			break
		}
		// 跳过空输入
		if strings.TrimSpace(input) == "" {
			continue
		}

		reply, err := agent.Run(ctx, toolList, input)
		if err != nil {
			fmt.Println("错误:", err)
			continue
		}
		fmt.Println(reply)
	}
}

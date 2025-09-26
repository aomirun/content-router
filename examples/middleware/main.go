package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aomirun/content-router"
	"github.com/aomirun/content-router/middleware"
)

func main() {
	// 创建路由器实例
	r := contentrouter.NewRouter()

	// 注册中间件
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// 注册路由处理器 - 匹配包含"trigger"的路径
	r.Match("trigger", func(ctx contentrouter.Context) error {
		// 模拟处理时间
		time.Sleep(100 * time.Millisecond)

		// 获取请求数据
		data := ctx.Buffer().Get()
		fmt.Printf("Processing data: %s\n", string(data))

		// 模拟一个可能的panic情况
		if string(data) == "trigger_panic" {
			panic("Simulated panic for demonstration")
		}

		// 模拟处理错误
		if string(data) == "trigger_error" {
			return fmt.Errorf("simulated processing error")
		}

		// 正常处理完成
		fmt.Println("Data processed successfully")
		return nil
	})

	// 创建测试数据
	testData := []string{
		"Hello, World!",
		"trigger_error",
		"trigger_panic",
	}

	// 测试不同的数据
	for i, data := range testData {
		fmt.Printf("\n=== Test Case %d ===\n", i+1)

		// 创建buffer
		buf := contentrouter.NewBuffer()
		buf.Write([]byte(data))

		// 处理数据
		_, err := r.Route(context.Background(), buf)
		if err != nil {
			fmt.Printf("Error processing data: %v\n", err)
		}
	}
}

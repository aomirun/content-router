package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aomirun/content-router"
)

func main() {
	fmt.Println("Router root package is working correctly!")

	// 创建路由器
	router := contentrouter.NewRouter()

	// 注册路由处理器
	router.Match("Hello", func(ctx contentrouter.Context) error {
		// 处理逻辑
		fmt.Printf("Handling route with buffer: %s\n", string(ctx.Buffer().Get()))
		return nil
	})

	// 创建缓冲区并写入数据
	buf := contentrouter.NewBuffer()
	buf.WriteString("Hello, World!")

	// 路由处理
	_, err := router.Route(context.Background(), buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Simple example completed successfully!")

	// 验证接口导入
	var _ contentrouter.Context = nil
	var _ contentrouter.Buffer = nil
	var _ contentrouter.Router = nil
	var _ contentrouter.Handler = nil
	var _ contentrouter.Matcher = nil
	var _ contentrouter.MiddlewareFunc = nil
	var _ contentrouter.Pipeline = nil

	fmt.Println("All interfaces are properly imported and accessible!")
}

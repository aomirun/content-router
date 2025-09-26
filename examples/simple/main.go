package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aomirun/content-router/api"
)

func main() {
	fmt.Println("Router API package is working correctly!")

	// 创建路由器
	router := api.NewRouter()

	// 注册路由处理器
	router.Match("Hello", func(ctx api.Context) error {
		// 处理逻辑
		fmt.Printf("Handling route with buffer: %s\n", string(ctx.Buffer().Get()))
		return nil
	})

	// 创建缓冲区并写入数据
	buf := api.NewBuffer()
	buf.WriteString("Hello, World!")

	// 路由处理
	_, err := router.Route(context.Background(), buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Simple example completed successfully!")

	// 验证接口导入
	var _ api.Context = nil
	var _ api.Buffer = nil
	var _ api.Router = nil
	var _ api.Handler = nil
	var _ api.Matcher = nil
	var _ api.MiddlewareFunc = nil
	var _ api.Pipeline = nil

	fmt.Println("All interfaces are properly imported and accessible!")
}

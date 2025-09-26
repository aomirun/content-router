package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aomirun/content-router/api"
)

// 自定义处理器示例
func customHandler(ctx api.Context) error {
	fmt.Printf("Custom handler processing buffer: %s\n", string(ctx.Buffer().Get()))
	return nil
}

// 自定义中间件示例
func loggingMiddleware(ctx api.Context, next api.HandlerFunc) error {
	fmt.Println("Before processing...")
	err := next(ctx)
	fmt.Println("After processing...")
	return err
}

func main() {
	fmt.Println("Fine-grained interfaces example")

	// 1. 直接使用Buffer接口
	buf := api.NewBuffer()
	buf.WriteString("Hello from fine-grained interfaces!")

	// 2. 直接使用Context接口
	_ = api.NewContext(context.Background(), buf)

	// 3. 使用Router接口
	router := api.NewRouter()

	// 4. 注册路由和中间件
	router.Use(loggingMiddleware)
	router.Match("Hello", customHandler)

	// 5. 路由处理
	_, err := router.Route(context.Background(), buf)
	if err != nil {
		log.Fatal(err)
	}

	// 6. 使用BufferManager
	bufferManager := router.BufferManager()
	newBuf := bufferManager.Acquire()
	newBuf.WriteString("Acquired from buffer manager!")
	fmt.Printf("Buffer manager test: %s\n", string(newBuf.Get()))
	bufferManager.Release(newBuf)

	fmt.Println("Fine-grained interfaces example completed successfully!")

	// 验证细粒度接口
	var _ api.RouteHandler = nil
	var _ api.RouteRegistrar = nil
	var _ api.MiddlewareHandler = nil
	var _ api.PipelineManager = nil
	var _ api.ContextCreator = nil
	var _ api.BufferManagerAccessor = nil

	// Context相关的细粒度接口
	var _ api.ValueStore = nil
	var _ api.BufferAccessor = nil

	// Buffer相关的细粒度接口
	var _ api.Readable = nil
	var _ api.Writable = nil
	var _ api.Mutable = nil
	var _ api.Sliceable = nil
	var _ api.Cloneable = nil

	fmt.Println("All fine-grained interfaces are properly imported and accessible!")
}

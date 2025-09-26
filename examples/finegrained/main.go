package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aomirun/content-router"
)

// 自定义处理器示例
func customHandler(ctx contentrouter.Context) error {
	fmt.Printf("Custom handler processing buffer: %s\n", string(ctx.Buffer().Get()))
	return nil
}

// 自定义中间件示例
func loggingMiddleware(ctx contentrouter.Context, next contentrouter.HandlerFunc) error {
	fmt.Println("Before processing...")
	err := next(ctx)
	fmt.Println("After processing...")
	return err
}

func main() {
	fmt.Println("Fine-grained interfaces example")

	// 1. 直接使用Buffer接口
	buf := contentrouter.NewBuffer()
	buf.WriteString("Hello from fine-grained interfaces!")

	// 2. 直接使用Context接口
	_ = contentrouter.NewContext(context.Background(), buf)

	// 3. 使用Router接口
	router := contentrouter.NewRouter()

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
	var _ contentrouter.RouteHandler = nil
	var _ contentrouter.RouteRegistrar = nil
	var _ contentrouter.MiddlewareHandler = nil
	var _ contentrouter.PipelineManager = nil
	var _ contentrouter.ContextCreator = nil
	var _ contentrouter.BufferManagerAccessor = nil

	// Context相关的细粒度接口
	var _ contentrouter.ValueStore = nil
	var _ contentrouter.BufferAccessor = nil

	// Buffer相关的细粒度接口
	var _ contentrouter.Readable = nil
	var _ contentrouter.Writable = nil
	var _ contentrouter.Mutable = nil
	var _ contentrouter.Sliceable = nil
	var _ contentrouter.Cloneable = nil

	fmt.Println("All fine-grained interfaces are properly imported and accessible!")
}

package router

import (
	router_context "github.com/aomirun/content-router/context"
)

// Pipeline 定义责任链管道接口
// 所有管道实现应该遵循此接口，提供一致的处理逻辑
//
// 实现此接口的类型应该确保:
// 1. 线程安全性
// 2. 中间件按注册顺序执行
// 3. 可以动态添加和移除中间件
//
// 命名规范:
// - 管道实例: xxxPipeline
// - 管道实现: xxxPipelineImpl
type Pipeline interface {
	// Use 添加中间件到管道
	//  - middleware: 中间件列表，用于在处理前后执行额外逻辑
	Use(middleware ...MiddlewareFunc)

	// Handle 处理内容，执行中间件链
	//  - ctx: 请求上下文
	// 返回: 可能的错误
	Handle(ctx router_context.Context) error
}

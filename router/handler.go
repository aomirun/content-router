package router

import (
	router_context "github.com/aomirun/content-router/context"
)

// Handler 定义处理函数接口
// 所有处理函数实现应该遵循此接口，提供一致的处理逻辑
// T 是处理函数处理的内容类型参数
//
// 实现此接口的类型应该确保:
// 1. 线程安全性
// 2. 处理逻辑的准确性
// 3. 合理的性能考虑
//
// 命名规范:
// - 处理方法: HandleXXX(ctx Context) error
// - 处理函数实例: xxxHandler
// - 处理函数实现: xxxHandlerImpl
type Handler interface {
	// Handle 处理内容
	//  - ctx: 上下文对象，用于传递请求信息和状态
	// 返回: 可能的错误
	Handle(ctx router_context.Context) error
}

// HandlerFunc 定义处理器函数类型
// 它是一个函数类型，用于处理内容
// ctx: 上下文信息
// 返回: 可能的错误
type HandlerFunc func(ctx router_context.Context) error

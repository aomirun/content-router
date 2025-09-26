package router

import (
	router_context "github.com/aomirun/content-router/context"
)

// MiddlewareFunc 定义中间件函数类型
// T: 泛型类型参数，用于指定处理器函数的参数类型
// 它是一个函数类型，用于在处理前后执行额外逻辑
// ctx: 请求上下文
// next: 下一个处理器函数
// 返回: 可能的错误
type MiddlewareFunc func(ctx router_context.Context, next HandlerFunc) error

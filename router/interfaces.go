package router

import (
	"context"

	"github.com/aomirun/content-router/buffer"
	router_context "github.com/aomirun/content-router/context"
	"github.com/aomirun/content-router/manage"
)

// RouteHandler 定义路由处理器接口
type RouteHandler interface {
	// Route 使用Buffer进行消息路由，减少数据复制
	//  - ctx: 上下文，用于传递请求范围的值和控制超时
	//  - buffer: 要路由的消息内容，以Buffer形式提供
	// 返回: 处理结果（可能是同一个Buffer）和可能的错误
	Route(ctx context.Context, buffer buffer.Buffer) (buffer.Buffer, error)
}

// RouteRegistrar 定义路由注册接口
type RouteRegistrar interface {
	// Register 注册新的路由规则
	//  - matcher: 内容匹配器，用于判断消息是否匹配
	//  - handler: 消息处理器，用于处理匹配的消息
	Register(matcher Matcher, handler HandlerFunc)

	// Match 注册基于字符串前缀的路由规则
	// pattern: 匹配模式
	// 支持的匹配模式:
	//  - "/regex/正则表达式": 符合正则表达式的消息
	//  - "/contains/特征值": 包含特征值的消息
	//  - "/prefix/前缀": 以指定前缀开头的消息
	//  - "/suffix/后缀": 以指定后缀结尾的消息
	// handler: 消息处理器，用于处理匹配的消息
	Match(pattern string, handler HandlerFunc)
}

// MiddlewareHandler 定义中间件处理接口
type MiddlewareHandler interface {
	// Use 添加中间件
	//  - middleware: 中间件列表，用于在处理前后执行额外逻辑
	Use(middleware ...MiddlewareFunc)
}

// PipelineManager 定义管道管理接口
type PipelineManager interface {
	// Pipeline 创建一个新的责任链管道，并与指定的匹配器关联
	//  - matcher: 内容匹配器，用于判断消息是否匹配
	// 返回: 新创建的管道
	Pipeline(matcher Matcher) Pipeline
}

// ContextCreator 定义上下文创建接口
type ContextCreator interface {
	// NewContext 创建一个新的增强上下文
	//  - parent: 父上下文
	//  - buffer: 关联的缓冲区
	// 返回: 新创建的上下文
	NewContext(parent context.Context, buffer buffer.Buffer) router_context.Context
}

// BufferManagerAccessor 定义缓冲区管理器访问接口
type BufferManagerAccessor interface {
	// BufferManager 获取BufferManager接口
	BufferManager() manage.BufferManager
}

// Router 定义路由器接口
// 它组合了所有路由器功能接口
type Router interface {
	RouteHandler
	RouteRegistrar
	MiddlewareHandler
	PipelineManager
	ContextCreator
	BufferManagerAccessor
}

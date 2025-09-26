package contentrouter

import (
	"context"

	"github.com/aomirun/content-router/buffer"
	router_context "github.com/aomirun/content-router/context"
	"github.com/aomirun/content-router/manage"
	"github.com/aomirun/content-router/router"
)

// Router 定义路由器接口
type Router = router.Router

// RouteHandler 定义路由处理器接口
type RouteHandler = router.RouteHandler

// RouteRegistrar 定义路由注册接口
type RouteRegistrar = router.RouteRegistrar

// MiddlewareHandler 定义中间件处理接口
type MiddlewareHandler = router.MiddlewareHandler

// PipelineManager 定义管道管理接口
type PipelineManager = router.PipelineManager

// ContextCreator 定义上下文创建接口
type ContextCreator = router.ContextCreator

// BufferManagerAccessor 定义缓冲区管理器访问接口
type BufferManagerAccessor = router.BufferManagerAccessor

// Context 定义增强的上下文接口
type Context = router_context.Context

// ValueStore 定义键值存储接口
type ValueStore = router_context.ValueStore

// BufferAccessor 定义缓冲区访问接口
type BufferAccessor = router_context.BufferAccessor

// Buffer 定义可重用的缓冲区接口
type Buffer = buffer.Buffer

// Readable 定义可读缓冲区接口
type Readable = buffer.Readable

// Writable 定义可写缓冲区接口
type Writable = buffer.Writable

// Mutable 定义可变缓冲区接口
type Mutable = buffer.Mutable

// Sliceable 定义可切片缓冲区接口
type Sliceable = buffer.Sliceable

// Cloneable 定义可克隆缓冲区接口
type Cloneable = buffer.Cloneable

// BufferManager 定义缓冲区管理接口
type BufferManager = manage.BufferManager

// Handler 定义处理函数接口
type Handler = router.Handler

// HandlerFunc 定义处理器函数类型
type HandlerFunc = router.HandlerFunc

// Matcher 定义内容匹配器接口
type Matcher = router.Matcher

// MatcherFunc 定义匹配器函数类型
type MatcherFunc = router.MatcherFunc

// MiddlewareFunc 定义中间件函数类型
type MiddlewareFunc = router.MiddlewareFunc

// Pipeline 定义责任链管道接口
type Pipeline = router.Pipeline

// ObjectPool 定义通用对象池接口
type ObjectPool[T any] = buffer.ObjectPool[T]

// NewRouter 创建一个新的路由器实例
func NewRouter() Router {
	return router.NewRouter()
}

// NewBuffer 创建一个新的缓冲区实例
func NewBuffer() Buffer {
	return buffer.NewBuffer()
}

// NewContext 创建一个新的上下文实例
func NewContext(parent context.Context, buf Buffer) Context {
	return router_context.NewContext(parent, buf)
}
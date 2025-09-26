package router

import (
	"bytes"
	"context"

	"github.com/aomirun/content-router/buffer"
	router_context "github.com/aomirun/content-router/context"
	"github.com/aomirun/content-router/manage"
)

// routerImpl 是Router接口的具体实现
type routerImpl struct {
	bufferManager manage.BufferManager
	routes        []routeEntry
	middlewares   []MiddlewareFunc
	pipelines     []pipelineEntry
	handlerChain  HandlerFunc
	dirty         bool // 标记路由或中间件是否发生变化
}

// routeEntry 定义路由条目
type routeEntry struct {
	matcher Matcher
	handler HandlerFunc
}

// pipelineEntry 定义管道条目
type pipelineEntry struct {
	matcher  Matcher
	pipeline Pipeline
}

// NewRouter 创建一个新的路由器实例
func NewRouter() Router {
	return &routerImpl{
		bufferManager: manage.NewBufferManager(),
		routes:        make([]routeEntry, 0),
		middlewares:   make([]MiddlewareFunc, 0),
		pipelines:     make([]pipelineEntry, 0),
	}
}

// Route 使用Buffer进行消息路由，减少数据复制
func (r *routerImpl) Route(ctx context.Context, buffer buffer.Buffer) (buffer.Buffer, error) {
	// 创建路由器上下文
	routerCtx := router_context.NewContext(ctx, buffer)

	// 应用全局中间件
	handler := r.buildHandlerChain()

	// 执行处理链
	err := handler(routerCtx)

	// 如果上下文实现了Reset方法，则重置它
	if resettable, ok := routerCtx.(interface{ Reset() }); ok {
		resettable.Reset()
	}

	return buffer, err
}

// buildHandlerChain 构建处理链
func (r *routerImpl) buildHandlerChain() HandlerFunc {
	// 如果处理链未变化，直接返回缓存的处理链
	if !r.dirty && r.handlerChain != nil {
		return r.handlerChain
	}

	// 基础处理器
	baseHandler := func(ctx router_context.Context) error {
		// 查找匹配的路由
		for _, entry := range r.routes {
			if entry.matcher.Match(ctx) {
				return entry.handler(ctx)
			}
		}
		return nil
	}

	// 如果没有中间件，直接返回基础处理器并缓存
	if len(r.middlewares) == 0 {
		r.handlerChain = baseHandler
		r.dirty = false
		return baseHandler
	}

	// 从后往前应用中间件（符合中间件链的常规做法）
	handler := baseHandler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		middleware := r.middlewares[i]
		next := handler
		handler = func(ctx router_context.Context) error {
			return middleware(ctx, next)
		}
	}

	// 缓存处理链并重置dirty标记
	r.handlerChain = handler
	r.dirty = false

	return handler
}

// Register 注册新的路由规则
func (r *routerImpl) Register(matcher Matcher, handler HandlerFunc) {
	r.routes = append(r.routes, routeEntry{
		matcher: matcher,
		handler: handler,
	})
	r.dirty = true
}

// Match 注册基于字符串前缀的路由规则
func (r *routerImpl) Match(pattern string, handler HandlerFunc) {
	// 简单实现：只支持前缀匹配
	patternBytes := []byte(pattern)
	matcher := MatcherFunc(func(ctx router_context.Context) bool {
		data := ctx.Buffer().Get()
		return len(data) >= len(patternBytes) && bytes.HasPrefix(data, patternBytes)
	})

	r.Register(matcher, handler)
}

// Use 添加中间件
func (r *routerImpl) Use(middleware ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware...)
	r.dirty = true
}

// Pipeline 创建一个新的责任链管道，并与指定的匹配器关联
func (r *routerImpl) Pipeline(matcher Matcher) Pipeline {
	// 简单实现：创建一个新的管道
	pipeline := &pipelineImpl{
		middlewares: make([]MiddlewareFunc, 0),
	}

	r.pipelines = append(r.pipelines, pipelineEntry{
		matcher:  matcher,
		pipeline: pipeline,
	})

	return pipeline
}

// NewContext 创建一个新的增强上下文
func (r *routerImpl) NewContext(parent context.Context, buffer buffer.Buffer) router_context.Context {
	return router_context.NewContext(parent, buffer)
}

// BufferManager 获取BufferManager接口
func (r *routerImpl) BufferManager() manage.BufferManager {
	return r.bufferManager
}

// pipelineImpl 是Pipeline接口的简单实现
type pipelineImpl struct {
	middlewares []MiddlewareFunc
}

// Use 添加中间件到管道
func (p *pipelineImpl) Use(middleware ...MiddlewareFunc) {
	p.middlewares = append(p.middlewares, middleware...)
}

// Handle 处理内容，执行中间件链
func (p *pipelineImpl) Handle(ctx router_context.Context) error {
	// 基础处理器
	baseHandler := func(ctx router_context.Context) error {
		// 管道的最终处理逻辑（这里简化处理）
		return nil
	}

	// 如果没有中间件，直接返回基础处理器
	if len(p.middlewares) == 0 {
		return baseHandler(ctx)
	}

	// 从后往前应用中间件
	handler := baseHandler
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		middleware := p.middlewares[i]
		next := handler
		handler = func(ctx router_context.Context) error {
			return middleware(ctx, next)
		}
	}

	return handler(ctx)
}

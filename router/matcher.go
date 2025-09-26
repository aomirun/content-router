package router

import (
	router_context "github.com/aomirun/content-router/context"
)

// Matcher 定义内容匹配器接口
// 所有匹配器实现应该遵循此接口，提供一致的匹配逻辑
//
// 实现此接口的类型应该确保:
// 1. 线程安全性
// 2. 匹配逻辑的准确性
// 3. 合理的性能考虑
//
// 命名规范:
// - 匹配方法: MatchXXX(ctx Context) bool
// - 匹配器实例: xxxMatcher
// - 匹配器实现: xxxMatcherImpl
type Matcher interface {
	// Match 检查内容是否匹配
	// ctx: 请求上下文
	// 返回: 如果匹配则返回true，否则返回false
	Match(ctx router_context.Context) bool
}

// MatcherFunc 定义匹配器函数类型
// 它是一个函数类型，用于检查内容是否匹配
// ctx: 请求上下文
// 返回: 如果匹配则返回true，否则返回false
type MatcherFunc func(ctx router_context.Context) bool

// Match 检查内容是否匹配
// ctx: 请求上下文
// 返回: 如果匹配则返回true，否则返回false
func (f MatcherFunc) Match(ctx router_context.Context) bool {
	return f(ctx)
}

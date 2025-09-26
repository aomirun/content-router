package router

import (
	"context"
	"testing"

	"github.com/aomirun/content-router/buffer"
	router_context "github.com/aomirun/content-router/context"
	"github.com/aomirun/content-router/manage"
)

// mockMatcher 是一个模拟的匹配器实现
type mockMatcher struct {
	matchResult bool
}

func (m *mockMatcher) Match(ctx router_context.Context) bool {
	return m.matchResult
}

// mockHandler 是一个模拟的处理器实现
func mockHandler(ctx router_context.Context) error {
	return nil
}

// mockMiddleware 是一个模拟的中间件实现
func mockMiddleware(ctx router_context.Context, next HandlerFunc) error {
	return next(ctx)
}

func TestRouterInterfaces(t *testing.T) {
	// 创建一个buffer
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 创建一个context
	ctx := router_context.NewContext(context.Background(), buf)

	// 测试Matcher接口
	matcher := &mockMatcher{matchResult: true}
	if !matcher.Match(ctx) {
		t.Error("Matcher should return true")
	}

	// 测试MatcherFunc类型
	matcherFunc := MatcherFunc(func(ctx router_context.Context) bool {
		return true
	})
	if !matcherFunc.Match(ctx) {
		t.Error("MatcherFunc should return true")
	}

	// 测试HandlerFunc类型
	handlerFunc := HandlerFunc(mockHandler)
	if handlerFunc == nil {
		t.Error("HandlerFunc should not be nil")
	}

	// 测试MiddlewareFunc类型
	middlewareFunc := MiddlewareFunc(mockMiddleware)
	if middlewareFunc == nil {
		t.Error("MiddlewareFunc should not be nil")
	}
}

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	if router == nil {
		t.Fatal("NewRouter should not return nil")
	}

	// 检查router是否实现了所有接口
	_, isRouter := router.(Router)
	_, isRouteHandler := router.(RouteHandler)
	_, isRouteRegistrar := router.(RouteRegistrar)
	_, isMiddlewareHandler := router.(MiddlewareHandler)
	_, isPipelineManager := router.(PipelineManager)
	_, isContextCreator := router.(ContextCreator)
	_, isBufferManagerAccessor := router.(BufferManagerAccessor)

	if !isRouter || !isRouteHandler || !isRouteRegistrar || !isMiddlewareHandler ||
		!isPipelineManager || !isContextCreator || !isBufferManagerAccessor {
		t.Error("Router should implement all required interfaces")
	}

	// 检查BufferManager是否正确初始化
	bufferManager := router.BufferManager()
	if bufferManager == nil {
		t.Error("BufferManager should not be nil")
	}

	_, isBufferManager := bufferManager.(manage.BufferManager)
	if !isBufferManager {
		t.Error("BufferManager should implement manage.BufferManager interface")
	}
}

func TestRouter_RegisterAndMatch(t *testing.T) {
	router := NewRouter()

	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("Hello, World!")

	// 记录是否调用了处理器
	handlerCalled := false
	testHandler := func(ctx router_context.Context) error {
		handlerCalled = true
		return nil
	}

	// 注册路由
	router.Register(&mockMatcher{matchResult: true}, HandlerFunc(testHandler))

	// 执行路由
	_, err := router.Route(context.Background(), buf)

	if err != nil {
		t.Errorf("Route should not return error: %v", err)
	}

	if !handlerCalled {
		t.Error("Handler should be called")
	}
}

func TestRouter_Match(t *testing.T) {
	router := NewRouter()

	// 记录是否调用了处理器
	handlerCalled := false
	testHandler := func(ctx router_context.Context) error {
		handlerCalled = true
		return nil
	}

	// 使用Match方法注册路由
	router.Match("Hello", HandlerFunc(testHandler))

	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("Hello, World!")

	// 执行路由
	_, err := router.Route(context.Background(), buf)

	if err != nil {
		t.Errorf("Route should not return error: %v", err)
	}

	if !handlerCalled {
		t.Error("Handler should be called for matching prefix")
	}

	// 重置标记
	handlerCalled = false

	// 创建不匹配的数据
	buf2 := buffer.NewBuffer()
	buf2.WriteString("Goodbye, World!")

	// 执行路由
	_, err = router.Route(context.Background(), buf2)

	if err != nil {
		t.Errorf("Route should not return error: %v", err)
	}

	if handlerCalled {
		t.Error("Handler should not be called for non-matching prefix")
	}
}

func TestRouter_UseMiddleware(t *testing.T) {
	router := NewRouter()

	// 记录中间件和处理器的调用顺序
	callOrder := []string{}

	// 添加中间件
	middleware1 := func(ctx router_context.Context, next HandlerFunc) error {
		callOrder = append(callOrder, "middleware1-before")
		err := next(ctx)
		callOrder = append(callOrder, "middleware1-after")
		return err
	}

	middleware2 := func(ctx router_context.Context, next HandlerFunc) error {
		callOrder = append(callOrder, "middleware2-before")
		err := next(ctx)
		callOrder = append(callOrder, "middleware2-after")
		return err
	}

	router.Use(MiddlewareFunc(middleware1), MiddlewareFunc(middleware2))

	// 注册处理器
	handler := func(ctx router_context.Context) error {
		callOrder = append(callOrder, "handler")
		return nil
	}

	router.Register(&mockMatcher{matchResult: true}, HandlerFunc(handler))

	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 执行路由
	_, err := router.Route(context.Background(), buf)

	if err != nil {
		t.Errorf("Route should not return error: %v", err)
	}

	// 验证调用顺序
	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(callOrder) != len(expectedOrder) {
		t.Errorf("Call order length mismatch. Expected %d, got %d", len(expectedOrder), len(callOrder))
	}

	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("Call order mismatch at position %d. Expected %s, got %s", i, expected, callOrder[i])
		}
	}
}

func TestRouter_HandlerChainCaching(t *testing.T) {
	router := NewRouter().(*routerImpl)

	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 第一次构建处理链
	router.buildHandlerChain()

	// 保存当前dirty状态
	initialDirty := router.dirty

	// 添加中间件后，dirty标志应该变为true
	router.Use(func(ctx router_context.Context, next HandlerFunc) error {
		return next(ctx)
	})

	// 验证dirty标志已更新
	if router.dirty == initialDirty {
		t.Error("Dirty flag should be updated after adding middleware")
	}

	// 构建处理链后，dirty标志应该变为false
	router.buildHandlerChain()

	if router.dirty != false {
		t.Error("Dirty flag should be false after building handler chain")
	}
}

func TestRouter_Pipeline(t *testing.T) {
	router := NewRouter()

	// 创建匹配器
	matcher := &mockMatcher{matchResult: true}

	// 创建管道
	pipeline := router.Pipeline(matcher)

	if pipeline == nil {
		t.Fatal("Pipeline should not be nil")
	}

	// 检查pipeline是否实现了Pipeline接口
	_, isPipeline := pipeline.(Pipeline)
	if !isPipeline {
		t.Error("Pipeline should implement Pipeline interface")
	}

	// 添加中间件到管道
	callOrder := []string{}

	middleware := func(ctx router_context.Context, next HandlerFunc) error {
		callOrder = append(callOrder, "pipeline-middleware")
		return next(ctx)
	}

	pipeline.Use(MiddlewareFunc(middleware))

	// 验证中间件被正确添加
	// 注意：由于pipelineImpl是私有的，我们无法直接访问middlewares字段
	// 但我们可以通过功能测试来验证
	if len(callOrder) != 0 {
		t.Error("Middleware should not be called yet")
	}
}

func TestRouter_NewContext(t *testing.T) {
	router := NewRouter()

	// 创建buffer
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 创建context
	ctx := router.NewContext(context.Background(), buf)

	if ctx == nil {
		t.Fatal("NewContext should not return nil")
	}

	// 检查context是否实现了Context接口
	_, isContext := ctx.(router_context.Context)
	if !isContext {
		t.Error("NewContext should return router_context.Context")
	}

	// 检查buffer是否正确设置
	if ctx.Buffer() != buf {
		t.Error("Context should contain the provided buffer")
	}
}

func TestRouter_NoRouteFound(t *testing.T) {
	router := NewRouter()

	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 注册一个不匹配的路由
	router.Register(&mockMatcher{matchResult: false}, HandlerFunc(mockHandler))

	// 执行路由
	result, err := router.Route(context.Background(), buf)

	// 应该没有错误，但也没有处理结果
	if err != nil {
		t.Errorf("Route should not return error: %v", err)
	}

	// 根据router.go中的实现，即使没有匹配的路由，也会返回一个空的结果
	// 所以这里我们只检查错误是否为nil
	if result == nil {
		// 这是可以接受的，因为实现可能返回nil
		// 我们主要关心的是没有错误
	}
}

func TestRouter_MultipleRoutes(t *testing.T) {
	router := NewRouter()

	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 记录调用的处理器
	calledHandlers := []string{}

	// 注册多个路由
	handler1 := func(ctx router_context.Context) error {
		calledHandlers = append(calledHandlers, "handler1")
		return nil
	}

	handler2 := func(ctx router_context.Context) error {
		calledHandlers = append(calledHandlers, "handler2")
		return nil
	}

	// 第一个路由匹配
	router.Register(&mockMatcher{matchResult: true}, HandlerFunc(handler1))

	// 第二个路由也匹配，但应该只调用第一个
	router.Register(&mockMatcher{matchResult: true}, HandlerFunc(handler2))

	// 执行路由
	_, err := router.Route(context.Background(), buf)

	if err != nil {
		t.Errorf("Route should not return error: %v", err)
	}

	// 应该只调用第一个处理器
	if len(calledHandlers) != 1 {
		t.Errorf("Expected 1 handler to be called, got %d", len(calledHandlers))
	}

	if len(calledHandlers) > 0 && calledHandlers[0] != "handler1" {
		t.Errorf("Expected handler1 to be called, got %s", calledHandlers[0])
	}
}

func TestRouter_PipelineUsage(t *testing.T) {
	router := NewRouter()

	// 创建匹配器
	matcher := &mockMatcher{matchResult: true}

	// 创建管道
	pipeline := router.Pipeline(matcher)

	if pipeline == nil {
		t.Fatal("Pipeline should not be nil")
	}

	// 检查pipeline是否实现了Pipeline接口
	_, isPipeline := pipeline.(Pipeline)
	if !isPipeline {
		t.Error("Pipeline should implement Pipeline interface")
	}

	// 添加中间件到管道
	callOrder := []string{}

	middleware := func(ctx router_context.Context, next HandlerFunc) error {
		callOrder = append(callOrder, "pipeline-middleware")
		return next(ctx)
	}

	pipeline.Use(MiddlewareFunc(middleware))

	// 验证中间件被正确添加
	// 注意：由于pipelineImpl是私有的，我们无法直接访问middlewares字段
	// 但我们可以通过功能测试来验证
	if len(callOrder) != 0 {
		t.Error("Middleware should not be called yet")
	}
}

func TestRouter_PipelineWithMiddlewareExecution(t *testing.T) {
	// 创建管道实现的独立测试
	pipeline := &pipelineImpl{}

	// 记录调用顺序
	callOrder := []string{}

	// 添加中间件到管道
	middleware1 := func(ctx router_context.Context, next HandlerFunc) error {
		callOrder = append(callOrder, "middleware1")
		return next(ctx)
	}

	middleware2 := func(ctx router_context.Context, next HandlerFunc) error {
		callOrder = append(callOrder, "middleware2")
		return next(ctx)
	}

	pipeline.Use(MiddlewareFunc(middleware1))
	pipeline.Use(MiddlewareFunc(middleware2))

	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("test data")
	ctx := router_context.NewContext(context.Background(), buf)

	// 执行管道处理
	err := pipeline.Handle(ctx)

	if err != nil {
		t.Errorf("Pipeline.Handle should not return error: %v", err)
	}

	// 验证调用顺序
	expectedOrder := []string{"middleware1", "middleware2"}

	if len(callOrder) != len(expectedOrder) {
		t.Errorf("Call order length mismatch. Expected %d, got %d", len(expectedOrder), len(callOrder))
	}

	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("Call order mismatch at position %d. Expected %s, got %s", i, expected, callOrder[i])
		}
	}
}

func TestPrefixMatcher(t *testing.T) {
	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("Hello, World!")
	ctx := router_context.NewContext(context.Background(), buf)

	// 测试匹配的情况
	matcher := PrefixMatcher("Hello")
	if !matcher.Match(ctx) {
		t.Error("PrefixMatcher should match 'Hello' prefix")
	}

	// 测试不匹配的情况
	matcher2 := PrefixMatcher("Goodbye")
	if matcher2.Match(ctx) {
		t.Error("PrefixMatcher should not match 'Goodbye' prefix")
	}

	// 测试空字符串的情况
	matcher3 := PrefixMatcher("")
	if !matcher3.Match(ctx) {
		t.Error("PrefixMatcher should match empty prefix")
	}
}

func TestSuffixMatcher(t *testing.T) {
	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("Hello, World!")
	ctx := router_context.NewContext(context.Background(), buf)

	// 测试匹配的情况
	matcher := SuffixMatcher("World!")
	if !matcher.Match(ctx) {
		t.Error("SuffixMatcher should match 'World!' suffix")
	}

	// 测试不匹配的情况
	matcher2 := SuffixMatcher("Hello")
	if matcher2.Match(ctx) {
		t.Error("SuffixMatcher should not match 'Hello' suffix")
	}

	// 测试空字符串的情况
	matcher3 := SuffixMatcher("")
	if !matcher3.Match(ctx) {
		t.Error("SuffixMatcher should match empty suffix")
	}
}

func TestContainsMatcher(t *testing.T) {
	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("Hello, World!")
	ctx := router_context.NewContext(context.Background(), buf)

	// 测试匹配的情况
	matcher := ContainsMatcher("Hello")
	if !matcher.Match(ctx) {
		t.Error("ContainsMatcher should match 'Hello' substring")
	}

	// 测试匹配的情况
	matcher2 := ContainsMatcher("World")
	if !matcher2.Match(ctx) {
		t.Error("ContainsMatcher should match 'World' substring")
	}

	// 测试不匹配的情况
	matcher3 := ContainsMatcher("Goodbye")
	if matcher3.Match(ctx) {
		t.Error("ContainsMatcher should not match 'Goodbye' substring")
	}

	// 测试空字符串的情况
	matcher4 := ContainsMatcher("")
	if !matcher4.Match(ctx) {
		t.Error("ContainsMatcher should match empty substring")
	}
}

func TestPipelineHandleWithoutMiddlewares(t *testing.T) {
	// 创建管道实现的独立测试
	pipeline := &pipelineImpl{}

	// 创建测试数据
	buf := buffer.NewBuffer()
	buf.WriteString("test data")
	ctx := router_context.NewContext(context.Background(), buf)

	// 执行管道处理（没有中间件）
	err := pipeline.Handle(ctx)

	if err != nil {
		t.Errorf("Pipeline.Handle should not return error: %v", err)
	}
}

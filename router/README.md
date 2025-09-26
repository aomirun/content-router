# Router 包

[English Version](README_en.md)

Router 包是 content-router 系统的核心组件，负责根据匹配规则路由内容，并通过中间件链和处理器处理内容。

## 功能特性

1. **路由**：基于buffer.Buffer实现消息路由，减少数据复制开销
2. **灵活路由规则**：支持多种匹配模式（前缀、后缀、包含、正则等）
3. **中间件支持**：提供全局中间件和管道中间件机制
4. **责任链模式**：通过Pipeline实现处理链的灵活组合
5. **上下文管理**：集成增强的上下文管理功能
6. **缓冲区管理**：内置缓冲区管理器，优化内存使用

## 核心接口

### Router接口
Router是核心接口，组合了所有路由器功能：

```go
type Router interface {
	RouteHandler
	RouteRegistrar
	MiddlewareHandler
	PipelineManager
	ContextCreator
	BufferManagerAccessor
}
```

### RouteHandler接口
定义消息路由处理功能：

```go
type RouteHandler interface {
	// Route 使用Buffer进行消息路由，减少数据复制
	Route(ctx context.Context, buffer buffer.Buffer) (buffer.Buffer, error)
}
```

### RouteRegistrar接口
定义路由注册功能：

```go
type RouteRegistrar interface {
	// Register 注册新的路由规则
	Register(matcher Matcher, handler HandlerFunc)
	
	// Match 注册基于字符串模式的路由规则
	Match(pattern string, handler HandlerFunc)
}
```

### MiddlewareHandler接口
定义中间件处理功能：

```go
type MiddlewareHandler interface {
	// Use 添加中间件
	Use(middleware ...MiddlewareFunc)
}
```

### PipelineManager接口
定义管道管理功能：

```go
type PipelineManager interface {
	// Pipeline 创建一个新的责任链管道
	Pipeline(matcher Matcher) Pipeline
}
```

### ContextCreator接口
定义上下文创建功能：

```go
type ContextCreator interface {
	// NewContext 创建一个新的增强上下文
	NewContext(parent context.Context, buffer buffer.Buffer) router_context.Context
}
```

### BufferManagerAccessor接口
定义缓冲区管理器访问功能：

```go
type BufferManagerAccessor interface {
	// BufferManager 获取BufferManager接口
	BufferManager() manage.BufferManager
}
```

## 核心组件

### Matcher（匹配器）
Matcher用于匹配消息内容是否符合路由规则：

```go
type Matcher interface {
	// Match 检查内容是否匹配
	Match(ctx router_context.Context) bool
}
```

提供了多种内置匹配器：
- PrefixMatcher：前缀匹配器
- SuffixMatcher：后缀匹配器
- ContainsMatcher：包含匹配器

### Middleware（中间件）
Middleware用于在处理前后执行额外逻辑：

```go
type MiddlewareFunc func(ctx router_context.Context, next HandlerFunc) error
```

### Pipeline（管道）
Pipeline实现了责任链模式，用于组织处理流程：

```go
type Pipeline interface {
	// Use 添加中间件到管道
	Use(middleware ...MiddlewareFunc)
	
	// Handle 处理内容，执行中间件链
	Handle(ctx router_context.Context) error
}
```

### Handler（处理器）
Handler定义了消息处理逻辑：

```go
type HandlerFunc func(ctx router_context.Context) error
```

## 实现细节

### routerImpl结构体
Router接口的具体实现，包含以下主要字段：
- `bufferManager`：缓冲区管理器
- `routes`：路由条目列表
- `middlewares`：全局中间件列表
- `pipelines`：管道条目列表
- `handlerChain`：处理链缓存
- `dirty`：路由或中间件变化标记

### 路由处理流程
1. 创建路由器上下文
2. 构建处理链（包含中间件和路由处理器）
3. 执行处理链
4. 重置上下文（如果支持）

### 处理链构建
- 缓存机制：避免重复构建相同的处理链
- 中间件应用：从后往前应用中间件（符合责任链模式）
- 路由匹配：按注册顺序查找匹配的路由

## 使用示例

### 基本路由注册和使用

```go
// 创建路由器实例
router := router.NewRouter()

// 注册路由规则
router.Match("/api/", func(ctx router_context.Context) error {
    // 处理以"/api/"开头的消息
    buf := ctx.Buffer()
    // 处理逻辑...
    return nil
})

// 创建缓冲区并写入数据
buf := router.BufferManager().Acquire()
buf.WriteString("Hello, World!")

// 执行路由
_, err := router.Route(context.Background(), buf)
if err != nil {
    // 处理错误
}

// 释放缓冲区
router.BufferManager().Release(buf)
```

### 中间件使用

```go
// 添加全局中间件
router.Use(func(ctx router_context.Context, next HandlerFunc) error {
    // 前置处理
    fmt.Println("Before processing")
    
    // 执行下一个处理器
    err := next(ctx)
    
    // 后置处理
    fmt.Println("After processing")
    
    return err
})
```

### 管道使用

```go
// 创建管道
matcher := router.PrefixMatcher("/api/")
pipeline := router.Pipeline(matcher)

// 为管道添加中间件
pipeline.Use(func(ctx router_context.Context, next HandlerFunc) error {
    // 管道特定的中间件逻辑
    return next(ctx)
})

// 注册使用管道的路由
router.Register(matcher, func(ctx router_context.Context) error {
    // 管道处理逻辑
    return pipeline.Handle(ctx)
})
```

## 线程安全性

### Router实例的线程安全性
Router实现不是完全线程安全的，在并发环境中需要注意：

1. **路由注册**：在并发环境中注册路由可能导致数据竞争
2. **中间件添加**：并发添加中间件可能导致数据竞争
3. **推荐使用模式**：
   - 在应用程序初始化阶段完成所有路由注册和中间件添加
   - 初始化完成后，Router实例可以在多个goroutine中并发使用进行消息路由

```go
// 推荐的使用方式
func main() {
    // 应用启动时初始化路由器
    router := router.NewRouter()
    
    // 注册所有路由和中间件
    router.Match("/api/", apiHandler)
    router.Use(loggingMiddleware)
    
    // 初始化完成后，并发使用路由功能
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // 并发执行路由（线程安全）
            _, _ = router.Route(context.Background(), buffer)
        }()
    }
    wg.Wait()
}
```

### 组件线程安全性
1. **Matcher**：实现应该是线程安全的
2. **Middleware**：实现应该是线程安全的
3. **Handler**：实现应该是线程安全的
4. **Pipeline**：实现应该是线程安全的

## 性能优化

1. **技术**：通过buffer.Buffer避免数据复制
2. **处理链缓存**：缓存构建好的处理链，避免重复构建
3. **对象池**：使用manage.BufferManager管理缓冲区
4. **延迟构建**：仅在需要时构建处理链

## 与其他组件的关系

1. **依赖buffer包**：使用buffer.Buffer进行数据传输
2. **依赖context包**：使用增强的上下文管理功能
3. **依赖manage包**：使用BufferManager进行缓冲区管理
4. **被api包使用**：api包通过类型别名导出Router接口
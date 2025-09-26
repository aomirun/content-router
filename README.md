# Content Router - 高性能路由器

Content Router是一个高性能、的路由器框架，专为Go语言设计。它专注于基于消息内容的路由处理，与Go 1.22的HTTP路径路由形成互补。通过细粒度接口设计和对象池技术，实现了高效的内存管理和路由处理。

## 目录

- [特性](#特性)
- [包结构](#包结构)
  - [API包说明](#api-package)
- [核心接口](#核心接口)
  - [Buffer相关接口](#buffer-interfaces)
  - [Context相关接口](#context-interfaces)
  - [Router相关接口](#router-interfaces)
- [使用示例](#usage-examples)
- [细粒度接口的优势](#fine-grained-advantages)
- [API包与子包的关系](#api-vs-subpackages)
- [性能优化](#performance)
- [测试](#testing)
- [安装](#installation)
- [许可证](#license)

## 特性

1. **设计**：通过Buffer接口和对象池技术，最大程度减少内存分配和复制
2. **细粒度接口**：将功能拆分为更小的接口，提供更好的接口隔离和组合性
3. **高性能**：使用对象池和高效的路由匹配算法
4. **易测试**：每个接口都可以独立测试和模拟
5. **向后兼容**：保持与旧版本的兼容性
6. **基于内容路由**：专注于基于消息内容的路由，与Go 1.22的HTTP路径路由形成互补

## 包结构

```
├── api              # 统一API入口
├── buffer           # 缓冲区管理
├── context          # 上下文管理
├── manage           # 资源管理
├── middleware       # 中间件
├── router           # 路由核心
└── examples         # 使用示例
    ├── simple       # 简单示例
    ├── finegrained  # 细粒度接口示例
    ├── middleware   # 中间件示例
    └── http         # HTTP服务器示例
```

### <a name="api-package"></a>API包说明

`api`包是项目的统一入口，它通过类型别名的方式导出所有核心接口和工厂函数，
为用户提供简化的使用方式。主要包含：

- 核心接口的类型别名（Router、Context、Buffer等）
- 工厂函数（NewRouter、NewBuffer、NewContext）

对于大多数用户来说，只需导入`api`包即可使用所有功能，无需关心内部包结构。
只有在需要细粒度控制或深入了解实现细节时，才需要直接导入子包。

Content Router专注于基于消息内容的路由处理，与Go 1.22的HTTP路径路由形成互补，为Go应用提供更全面的路由解决方案。

## <a name="core-interfaces"></a>核心接口

### <a name="buffer-interfaces"></a>Buffer相关接口
- `Readable` - 可读缓冲区接口
- `Writable` - 可写缓冲区接口
- `Mutable` - 可变缓冲区接口
- `Sliceable` - 可切片缓冲区接口
- `Cloneable` - 可克隆缓冲区接口
- `Buffer` - 组合所有缓冲区操作的接口

### <a name="context-interfaces"></a>Context相关接口
- `ValueStore` - 键值存储接口
- `BufferAccessor` - 缓冲区访问接口
- `Context` - 组合标准context.Context、ValueStore和BufferAccessor的增强上下文接口

### <a name="router-interfaces"></a>Router相关接口
- `RouteHandler` - 路由处理器接口
- `RouteRegistrar` - 路由注册接口
- `MiddlewareHandler` - 中间件处理接口
- `PipelineManager` - 管道管理接口
- `ContextCreator` - 上下文创建接口
- `BufferManagerAccessor` - 缓冲区管理器访问接口
- `Router` - 组合所有路由器功能的接口

## <a name="usage-examples"></a>使用示例

对于大多数用户，推荐使用`api`包作为入口，它提供了所有核心功能的统一访问接口。

### 简单使用
```go
// 创建路由器（通过api包）
router := api.NewRouter()

// 注册路由
router.Match("Hello", func(ctx api.Context) error {
    // 处理逻辑
    return nil
})

// 创建缓冲区并写入数据（通过api包）
buf := api.NewBuffer()
buf.WriteString("Hello, World!")

// 路由处理
router.Route(context.Background(), buf)
```

### HTTP服务器示例

这个示例展示了如何将Content Router集成到HTTP服务器中。虽然Go 1.22增强了HTTP路径路由功能，但Content Router专注于基于消息内容的路由处理，两者可以互补使用。

```go
func httpHandler(w http.ResponseWriter, r *http.Request) {
    // 创建缓冲区（通过api包）
    buf := api.NewBuffer()
    buf.WriteString("HTTP request: " + r.URL.Path)
    
    // 创建路由器（通过api包）
    router := api.NewRouter()
    
    // 注册基于内容的路由规则
    // 这里根据消息内容而不是URL路径进行路由
    router.Match("Hello", func(ctx api.Context) error {
        response := "Processed: " + string(ctx.Buffer().Get())
        fmt.Fprintf(w, "%s", response)
        return nil
    })
    
    // 处理请求
    router.Route(context.Background(), buf)
}
```

### 细粒度接口使用（高级用法）

对于需要更细粒度控制的场景，可以直接使用各子包：

```go
// 直接使用router包
import "github.com/aomirun/content-router/router"
import "github.com/aomirun/content-router/buffer"
import "github.com/aomirun/content-router/context"

// 创建路由器（直接使用router包）
router := router.NewRouter()

// 创建缓冲区（直接使用buffer包）
buf := buffer.NewBuffer()
buf.WriteString("Hello, World!")

// 创建上下文（直接使用context包）
ctx := context.NewContext(context.Background(), buf)

// 路由处理
router.Route(ctx, buf)
```

### 中间件使用示例

Content Router支持中间件功能，可以用于日志记录、错误恢复等：

```go
// 创建路由器
router := api.NewRouter()

// 添加中间件
router.Use(
    // 日志中间件
    middleware.Logging(),
    // 错误恢复中间件
    middleware.Recovery(),
)

// 注册路由
router.Match("Hello", func(ctx api.Context) error {
    fmt.Println("处理Hello消息")
    return nil
})
```

## <a name="fine-grained-advantages"></a>细粒度接口的优势

1. **更好的接口隔离** - 组件只依赖它们实际使用的功能
2. **更容易测试** - 可以只模拟特定的接口而不是整个大接口
3. **更好的可组合性** - 可以根据需要组合不同的接口
4. **更清晰的职责分离** - 每个接口都有明确的职责

### <a name="api-vs-subpackages"></a>API包与子包的关系

`api`包通过类型别名的方式导出各子包的接口和函数，提供了简化的使用方式：

```go
// api/api.go中的定义示例
type Router = router.Router
type Context = context.Context
type Buffer = buffer.Buffer

func NewRouter() Router {
    return router.NewRouter()
}

func NewBuffer() Buffer {
    return buffer.NewBuffer()
}
```

这种设计既保持了细粒度接口的优势，又为用户提供了便捷的使用方式。用户可以根据需要选择使用层次：
- **初学者或一般用途**：使用`api`包即可
- **高级用户或特殊需求**：直接导入相应的子包

## <a name="performance"></a>性能优化

1. **对象池**：使用sync.Pool实现Buffer对象池，减少内存分配
2. ****：通过Buffer接口和引用传递，避免不必要的数据复制
3. **高效路由匹配**：实现多种匹配器（前缀、后缀、包含等）

## <a name="testing"></a>测试

运行所有测试：
```bash
go test ./...
```

### 测试覆盖情况

项目具有全面的测试覆盖，确保各组件的稳定性和可靠性：

- `buffer` 包: 100% 测试覆盖率
- `context` 包: 100% 测试覆盖率
- `manage` 包: 100% 测试覆盖率
- `router` 包: 100% 测试覆盖率

各包都包含单元测试，验证核心功能的正确性。

### 基准测试

运行基准测试：
```bash
go test -bench=. ./benchmark/...

goos: linux
goarch: amd64
pkg: github.com/aomirun/content-router/benchmark
cpu: AMD Ryzen 5 3600 6-Core Processor              
BenchmarkRouter_Route-12                        22929966                51.64 ns/op            0 B/op          0 allocs/op
BenchmarkRouter_RouteWithMiddleware-12          20555922                57.85 ns/op            0 B/op          0 allocs/op
BenchmarkBuffer_AcquireRelease-12               39527199                30.23 ns/op            0 B/op          0 allocs/op
BenchmarkContext_ValueStore-12                  22999408                47.52 ns/op            0 B/op          0 allocs/op
BenchmarkContextPool_AcquireRelease-12          31431421                36.62 ns/op            0 B/op          0 allocs/op
BenchmarkContext_Methods-12                      9394586               125.8 ns/op             0 B/op          0 allocs/op
BenchmarkContext_Keys-12                         9657076               121.9 ns/op            96 B/op          1 allocs/op
BenchmarkContext_Methods_WithoutKeys-12          8851708               126.9 ns/op             0 B/op          0 allocs/op
BenchmarkMatcher_Prefix-12                      151812487                7.678 ns/op           0 B/op          0 allocs/op
BenchmarkMatcher_Suffix-12                      145771711                8.426 ns/op           0 B/op          0 allocs/op
BenchmarkMatcher_Contains-12                    91538996                12.91 ns/op            0 B/op          0 allocs/op
BenchmarkPipeline_WithMiddleware-12             18005524                63.48 ns/op           48 B/op          2 allocs/op
BenchmarkPipeline_WithoutMiddleware-12          567985728                2.127 ns/op           0 B/op          0 allocs/op
BenchmarkBufferPool_AcquireRelease-12           57168133                20.06 ns/op            0 B/op          0 allocs/op
PASS

```

## <a name="installation"></a>安装

```bash
go get github.com/aomirun/content-router
```

## <a name="license"></a>许可证

MIT
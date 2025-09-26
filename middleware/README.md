# 中间件包

该包包含了路由器的各种中间件实现。

## 概述

本系统中的中间件遵循"洋葱模型"，其中每个中间件都包装链中的下一个中间件。当请求进来时，它会按顺序通过每个中间件，而响应则以相反的顺序返回。

```
请求    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
────────►│中间件 1     │─►│中间件 2     │─►│中间件 3     │─► 处理器
响应    ◄─────────────┘  ◄─────────────┘  ◄─────────────┘
```

## 中间件列表

### 1. 恢复中间件
- **文件**: `recovery.go`
- **用途**: 捕获请求处理过程中的`panic`并记录错误信息
- **特性**:
  - 捕获并记录`panic`堆栈跟踪
  - 防止由于未处理的`panic`导致应用程序崩溃
  - 可扩展以向监控系统发送错误报告

### 2. 日志中间件
- **文件**: `logging.go`
- **用途**: 记录请求处理信息
- **特性**:
  - 记录请求开始时间和数据预览
  - 测量并记录处理持续时间
  - 记录处理结果（成功/失败）

## 使用方法

要使用这些中间件，请导入它们并向路由器注册：

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/aomirun/content-router"
    "github.com/aomirun/content-router/middleware"
)

func main() {
    // 创建路由器实例
    r := contentrouter.NewRouter()
    
    // 注册中间件
    r.Use(middleware.LoggingMiddleware())
    r.Use(middleware.RecoveryMiddleware())
    
    // 注册路由处理器
    r.Match("test", func(ctx contentrouter.Context) error {
        // 模拟处理时间
        time.Sleep(100 * time.Millisecond)
        
        // 获取请求数据
        data := ctx.Buffer().Get()
        fmt.Printf("正在处理数据: %s\n", string(data))
        
        // 模拟恐慌情况
        if string(data) == "trigger_panic" {
            panic("为演示而模拟的恐慌")
        }
        
        // 模拟处理错误
        if string(data) == "trigger_error" {
            return fmt.Errorf("模拟的处理错误")
        }
        
        // 正常处理完成
        fmt.Println("数据处理成功")
        return nil
    })
    
    // 创建测试数据
    buf := contentrouter.NewBuffer()
    buf.Write([]byte("test_data"))
    
    // 处理数据
    _, err := r.Route(context.Background(), buf)
    if err != nil {
        fmt.Printf("处理数据时出错: %v\n", err)
    }
}
```

## 创建自定义中间件

要创建自定义中间件，请遵循以下模式：

```go
func CustomMiddleware() func(ctx router_context.Context, next func(router_context.Context) error) error {
    return func(ctx router_context.Context, next func(router_context.Context) error) error {
        // 预处理逻辑
        
        // 调用下一个中间件/处理器
        err := next(ctx)
        
        // 后处理逻辑
        
        return err
    }
}
```

中间件函数接收：
1. `ctx` - 请求上下文
2. `next` - 链中的下一个处理器

要继续链，请调用 `next(ctx)`。要中止链，请返回错误而不调用 `next`。

## 测试

中间件包包含以下测试函数：

1. `TestLoggingMiddleware` - 测试正常情况下的日志记录中间件
2. `TestLoggingMiddlewareWithError` - 测试带有错误的日志记录中间件
3. `TestRecoveryMiddleware` - 测试错误恢复中间件
4. `TestRecoveryMiddlewareWithoutPanic` - 测试没有panic时的错误恢复中间件
5. `TestLoggingMiddlewareWithLongData` - 测试日志记录中间件处理长数据的情况

使用以下命令运行测试：

```bash
go test -v ./middleware/...
```

测试覆盖率达到 100%。
# Context Package

Context包提供了一个增强的上下文实现，专为消息路由框架设计。它扩展了标准库的context.Context，增加了键值存储和缓冲区访问功能。

## 功能特性

1. **标准库兼容**：完全兼容`context.Context`接口
2. **键值存储**：支持多种类型的键值对存储和获取
3. **缓冲区访问**：直接关联和访问消息缓冲区
4. **上下文复制**：支持创建上下文的副本（Fork）
5. **对象池优化**：内部使用对象池减少内存分配
6. **类型安全**：提供类型安全的值获取方法

## 核心接口

### Context接口
Context是核心接口，组合了标准库的context.Context、ValueStore和BufferAccessor：

```go
type Context interface {
    context.Context
    ValueStore
    BufferAccessor

    Fork() Context
    ForkWithBuffer(buffer buffer.Buffer) Context
}
```

### ValueStore接口
提供键值存储功能：

```go
type ValueStore interface {
    Set(key, value interface{})
    Get(key interface{}) interface{}
    GetString(key interface{}) (string, bool)
    GetInt(key interface{}) (int, bool)
    GetInt64(key interface{}) (int64, bool)
    GetBool(key interface{}) (bool, bool)
    GetFloat64(key interface{}) (float64, bool)
    GetBytes(key interface{}) ([]byte, bool)
    GetTime(key interface{}) (time.Time, bool)
    Delete(key interface{})
    Keys() []interface{}
}
```

### BufferAccessor接口
提供缓冲区访问功能：

```go
type BufferAccessor interface {
    Buffer() buffer.Buffer
}
```

## 实现细节

### contextImpl结构体
Context接口的具体实现，包含：
- 嵌入标准库的`context.Context`
- 关联的`buffer.Buffer`实例
- 键值对存储的`map[interface{}]interface{}`

### 对象池优化
使用`sync.Pool`管理contextImpl实例，减少内存分配：
- `NewContext()`从池中获取实例
- 上下文使用完毕后自动放回池中

## 使用示例

```go
// 创建新的上下文实例
buf := buffer.NewBuffer()
ctx := context.NewContext(context.Background(), buf)

// 设置键值对
ctx.Set("key", "value")
ctx.Set("number", 42)

// 获取值
if val, ok := ctx.GetString("key"); ok {
    fmt.Println("key:", val) // 输出: key: value
}

// 获取关联的缓冲区
buffer := ctx.Buffer()
buffer.WriteString("Hello, World!")

// Fork创建副本
childCtx := ctx.Fork()
childCtx.Set("child", "data")

// ForkWithBuffer创建副本并使用新缓冲区
newBuf := buffer.NewBuffer()
newCtx := ctx.ForkWithBuffer(newBuf)
```

## 线程安全性

### Context实例的线程安全性
Context实现不是线程安全的，这意味着多个goroutine同时访问同一个Context实例可能会导致数据竞争。如果需要在并发环境中使用Context，建议采用以下模式：

1. **单goroutine使用**：每个goroutine使用自己的Context实例
2. **通过Fork创建副本**：在不同goroutine中使用Fork创建的副本
3. **传递所有权**：Context在goroutine间传递而不是共享

### 并发使用示例

```go
// 在多个goroutine中安全地使用Context
var wg sync.WaitGroup
parentCtx := context.NewContext(context.Background(), buf)

// 启动多个goroutine
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        
        // 通过Fork创建副本在不同goroutine中使用
        ctx := parentCtx.Fork()
        ctx.Set("goroutine_id", id)
        
        // 使用ctx...
        processData(ctx)
    }(i)
}

wg.Wait()
```

## 性能优化

1. **对象池**：使用`sync.Pool`管理Context实例，减少GC压力
2. **内存重用**：Fork操作会复制values map但共享缓冲区
3. **类型安全方法**：提供类型安全的值获取方法，避免类型断言错误
4. **预分配**：在Fork操作中预分配map容量

## 与标准库context的差异

1. **增强功能**：增加了键值存储和缓冲区访问功能
2. **对象池**：内部使用对象池管理实例
3. **Fork操作**：支持创建上下文副本
4. **类型安全**：提供多种类型安全的值获取方法
# Buffer Package

Buffer包提供了一个高效、可重用的缓冲区实现，专为消息处理场景设计。它支持零拷贝操作，减少内存分配，提高性能。

## 功能特性

1. **接口化设计**：通过细粒度接口定义缓冲区的各种操作
2. **高性能实现**：底层使用字节切片，支持动态扩容
3. **对象池支持**：内置对象池实现，减少GC压力
4. **零拷贝操作**：支持切片操作而不复制数据
5. **标准库兼容**：实现`io.Writer`和`io.StringWriter`接口

## 核心接口

### Buffer接口
Buffer是核心接口，组合了所有缓冲区操作接口：

```go
type Buffer interface {
    Readable
    Writable
    Mutable
    Sliceable
    Cloneable
}
```

### 功能接口

1. **Readable** - 可读操作
   - `Get() []byte` - 获取底层字节数组引用
   - `Len() int` - 获取有效数据长度
   - `Cap() int` - 获取缓冲区容量

2. **Writable** - 可写操作
   - `Write(p []byte) (n int, err error)` - 写入字节数据
   - `WriteString(s string) (n int, err error)` - 写入字符串数据

3. **Mutable** - 可变操作
   - `Reset()` - 重置缓冲区内容
   - `Truncate(n int)` - 截断到指定长度

4. **Sliceable** - 切片操作
   - `Slice(start, end int) Buffer` - 创建子切片（不复制数据）

5. **Cloneable** - 克隆操作
   - `Clone() Buffer` - 创建深拷贝

## 对象池

### ObjectPool接口
通用对象池接口，支持任何类型：

```go
type ObjectPool[T any] interface {
    Acquire() T
    Release(obj T)
    Size() int
}
```

### BufferPool
专门针对Buffer类型的对象池实现，自动重置归还的对象。

## 使用示例

```go
// 创建新的Buffer实例
buf := buffer.NewBuffer()

// 写入数据
buf.WriteString("Hello, World!")

// 读取数据
data := buf.Get()
fmt.Println(string(data)) // 输出: Hello, World!

// 使用对象池
pool := buffer.NewPool()
buf := pool.Acquire()
// 使用buf...
pool.Release(buf)
```

## 性能优化

1. **对象池**：通过`NewPool()`创建对象池，减少内存分配
2. **零拷贝**：`Slice()`方法创建子切片时不复制数据
3. **自动重置**：对象池自动重置归还的对象
4. **容量预分配**：默认初始容量1024字节，减少扩容次数

## 线程安全性

### Buffer实例的线程安全性
Buffer实现不是线程安全的，这意味着多个goroutine同时访问同一个Buffer实例可能会导致数据竞争。如果需要在并发环境中使用Buffer，建议采用以下模式：

1. **单goroutine使用**：每个goroutine使用自己的Buffer实例
2. **通过对象池管理**：从池中获取Buffer，在单个goroutine中使用，使用完后归还
3. **传递所有权**：Buffer在goroutine间传递而不是共享

### ObjectPool的线程安全性
ObjectPool实现是线程安全的，基于`sync.Pool`实现。以下操作可以在多个goroutine中并发调用：
- `Acquire()` - 从池中获取Buffer实例
- `Release(obj)` - 将Buffer实例归还给池

### 并发使用示例

```go
// 线程安全：在多个goroutine中使用对象池
var wg sync.WaitGroup
pool := buffer.NewPool()

// 启动多个goroutine
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        
        // 线程安全：从池中获取Buffer
        buf := pool.Acquire()
        
        // 注意：获取到的Buffer应在单个goroutine中使用
        buf.WriteString(fmt.Sprintf("Goroutine %d", id))
        processData(buf.Get())
        
        // 线程安全：将Buffer归还给池
        pool.Release(buf)
    }(i)
}

wg.Wait()
```

### 设计原则
这种设计遵循了Go语言的并发哲学：
1. **性能优先**：避免不必要的锁开销
2. **职责分离**：Buffer专注于数据操作，线程安全通过其他机制实现
3. **推荐模式**：通过对象池和goroutine隔离实现并发安全
# Buffer Package

[中文版本](README.md)

The Buffer package provides an efficient, reusable buffer implementation designed specifically for message processing scenarios. It supports zero-copy operations to reduce memory allocation and improve performance.

## Features

1. **Interface-based Design**: Fine-grained interfaces define various buffer operations
2. **High-performance Implementation**: Uses byte slices at the core with dynamic expansion support
3. **Object Pool Support**: Built-in object pool implementation reduces GC pressure
4. **Zero-copy Operations**: Supports slice operations without data copying
5. **Standard Library Compatibility**: Implements `io.Writer` and `io.StringWriter` interfaces

## Core Interfaces

### Buffer Interface
Buffer is the core interface that combines all buffer operation interfaces:

```go
type Buffer interface {
    Readable
    Writable
    Mutable
    Sliceable
    Cloneable
}
```

### Functional Interfaces

1. **Readable** - Read operations
   - `Get() []byte` - Get reference to underlying byte array
   - `Len() int` - Get length of valid data
   - `Cap() int` - Get buffer capacity

2. **Writable** - Write operations
   - `Write(p []byte) (n int, err error)` - Write byte data
   - `WriteString(s string) (n int, err error)` - Write string data

3. **Mutable** - Mutable operations
   - `Reset()` - Reset buffer content
   - `Truncate(n int)` - Truncate to specified length

4. **Sliceable** - Slice operations
   - `Slice(start, end int) Buffer` - Create sub-slice (without copying data)

5. **Cloneable** - Clone operations
   - `Clone() Buffer` - Create deep copy

## Object Pool

### ObjectPool Interface
Generic object pool interface supporting any type:

```go
type ObjectPool[T any] interface {
    Acquire() T
    Release(obj T)
    Size() int
}
```

### BufferPool
Specialized object pool implementation for Buffer types that automatically resets returned objects.

## Usage Example

```go
// Create a new Buffer instance
buf := buffer.NewBuffer()

// Write data
buf.WriteString("Hello, World!")

// Read data
data := buf.Get()
fmt.Println(string(data)) // Output: Hello, World!

// Use object pool
pool := buffer.NewPool()
buf := pool.Acquire()
// Use buf...
pool.Release(buf)
```

## Performance Optimization

1. **Object Pool**: Create object pools with `NewPool()` to reduce memory allocation
2. **Zero-copy**: The `Slice()` method creates sub-slices without copying data
3. **Automatic Reset**: Object pools automatically reset returned objects
4. **Capacity Pre-allocation**: Default initial capacity of 1024 bytes reduces expansion frequency

## Thread Safety

### Thread Safety of Buffer Instances
Buffer implementations are not thread-safe, which means that multiple goroutines accessing the same Buffer instance simultaneously may cause data races. If you need to use Buffer in concurrent environments, it is recommended to adopt the following patterns:

1. **Single-goroutine Usage**: Each goroutine uses its own Buffer instance
2. **Object Pool Management**: Acquire Buffer from pool, use in single goroutine, and release when done
3. **Ownership Transfer**: Pass Buffer between goroutines rather than sharing

### Thread Safety of ObjectPool
ObjectPool implementation is thread-safe, based on `sync.Pool`. The following operations can be called concurrently from multiple goroutines:
- `Acquire()` - Acquire Buffer instance from pool
- `Release(obj)` - Return Buffer instance to pool

### Concurrent Usage Example

```go
// Thread-safe: Use object pool in multiple goroutines
var wg sync.WaitGroup
pool := buffer.NewPool()

// Start multiple goroutines
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        
        // Thread-safe: Acquire Buffer from pool
        buf := pool.Acquire()
        
        // Note: The acquired Buffer should be used in a single goroutine
        buf.WriteString(fmt.Sprintf("Goroutine %d", id))
        processData(buf.Get())
        
        // Thread-safe: Return Buffer to pool
        pool.Release(buf)
    }(i)
}

wg.Wait()
```

### Design Principles
This design follows Go's concurrency philosophy:
1. **Performance First**: Avoid unnecessary lock overhead
2. **Separation of Concerns**: Buffer focuses on data operations, thread safety is achieved through other mechanisms
3. **Recommended Patterns**: Achieve concurrency safety through object pools and goroutine isolation
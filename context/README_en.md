# Context Package

The Context package provides an enhanced context implementation designed specifically for the message routing framework. It extends the standard library's context.Context with key-value storage and buffer access functionality.

## Features

1. **Standard Library Compatibility**: Fully compatible with the `context.Context` interface
2. **Key-Value Storage**: Supports storage and retrieval of multiple types of key-value pairs
3. **Buffer Access**: Direct association and access to message buffers
4. **Context Copying**: Supports creating copies of contexts (Fork)
5. **Object Pool Optimization**: Uses internal object pools to reduce memory allocation
6. **Type Safety**: Provides type-safe value retrieval methods

## Core Interfaces

### Context Interface
Context is the core interface that combines the standard library's context.Context, ValueStore, and BufferAccessor:

```go
type Context interface {
    context.Context
    ValueStore
    BufferAccessor

    Fork() Context
    ForkWithBuffer(buffer buffer.Buffer) Context
}
```

### ValueStore Interface
Provides key-value storage functionality:

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

### BufferAccessor Interface
Provides buffer access functionality:

```go
type BufferAccessor interface {
    Buffer() buffer.Buffer
}
```

## Implementation Details

### contextImpl Struct
The concrete implementation of the Context interface, containing:
- Embedded standard library `context.Context`
- Associated `buffer.Buffer` instance
- Key-value storage `map[interface{}]interface{}`

### Object Pool Optimization
Uses `sync.Pool` to manage contextImpl instances and reduce memory allocation:
- `NewContext()` acquires instances from the pool
- Contexts are automatically returned to the pool when finished

## Usage Example

```go
// Create a new context instance
buf := buffer.NewBuffer()
ctx := context.NewContext(context.Background(), buf)

// Set key-value pairs
ctx.Set("key", "value")
ctx.Set("number", 42)

// Get values
if val, ok := ctx.GetString("key"); ok {
    fmt.Println("key:", val) // Output: key: value
}

// Get associated buffer
buffer := ctx.Buffer()
buffer.WriteString("Hello, World!")

// Fork to create a copy
childCtx := ctx.Fork()
childCtx.Set("child", "data")

// ForkWithBuffer to create a copy with a new buffer
newBuf := buffer.NewBuffer()
newCtx := ctx.ForkWithBuffer(newBuf)
```

## Thread Safety

### Thread Safety of Context Instances
Context implementations are not thread-safe, which means that multiple goroutines accessing the same Context instance simultaneously may cause data races. If you need to use Context in concurrent environments, it is recommended to adopt the following patterns:

1. **Single-goroutine Usage**: Each goroutine uses its own Context instance
2. **Copy Creation with Fork**: Use Fork-created copies in different goroutines
3. **Ownership Transfer**: Pass Context between goroutines rather than sharing

### Concurrent Usage Example

```go
// Safely use Context in multiple goroutines
var wg sync.WaitGroup
parentCtx := context.NewContext(context.Background(), buf)

// Start multiple goroutines
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        
        // Create a copy using Fork for use in different goroutines
        ctx := parentCtx.Fork()
        ctx.Set("goroutine_id", id)
        
        // Use ctx...
        processData(ctx)
    }(i)
}

wg.Wait()
```

## Performance Optimization

1. **Object Pool**: Uses `sync.Pool` to manage Context instances, reducing GC pressure
2. **Memory Reuse**: Fork operations copy the values map but share the buffer
3. **Type-safe Methods**: Provides type-safe value retrieval methods to avoid type assertion errors
4. **Pre-allocation**: Pre-allocates map capacity in Fork operations

## Differences from Standard Library Context

1. **Enhanced Functionality**: Adds key-value storage and buffer access functionality
2. **Object Pool**: Uses internal object pools to manage instances
3. **Fork Operation**: Supports creating context copies
4. **Type Safety**: Provides multiple type-safe value retrieval methods
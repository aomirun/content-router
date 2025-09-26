# Manage Package

[中文版本](README.md)

The Manage package provides a buffer manager implementation for efficiently managing buffer resources in message processing. It is based on the object pool mechanism of the buffer package and provides unified interfaces for acquiring and releasing buffers.

## Features

1. **Buffer Pool Management**: Manages buffer instances based on the object pool mechanism
2. **Resource Efficient Utilization**: Reduces memory allocation and garbage collection pressure
3. **Unified Interface**: Provides simple and easy-to-use acquire and release interfaces
4. **Automatic Reset**: Automatically resets buffer state when released

## Core Interfaces

### BufferManager Interface
BufferManager is the core interface that defines basic buffer management operations:

```go
type BufferManager interface {
    // Acquire a buffer from the pool
    Acquire() buffer.Buffer
    
    // Release the buffer back to the pool
    Release(buf buffer.Buffer)
}
```

## Implementation Details

### bufferManagerImpl Struct
The concrete implementation of the BufferManager interface, containing:
- `pool` field: Points to the object pool instance in the buffer package

### Automatic Reset Mechanism
Buffers are automatically reset during the Release operation:
```go
func (bm *bufferManagerImpl) Release(buf buffer.Buffer) {
    // Reset the buffer before returning it to the pool
    buf.Reset()
    bm.pool.Release(buf)
}
```

## Usage Example

```go
// Create a BufferManager instance
manager := manage.NewBufferManager()

// Acquire a buffer
buf := manager.Acquire()

// Use the buffer
buf.WriteString("Hello, World!")
data := buf.Get()

// Release the buffer after processing
manager.Release(buf)
```

## Thread Safety

### Thread Safety of BufferManager Instances
BufferManager implementation is thread-safe because:

1. **Thread-safe Underlying Object Pool**: Based on the ObjectPool implementation in the buffer package, whose Acquire and Release operations are thread-safe
2. **Stateless Design**: BufferManager itself does not maintain state; all state is managed by the underlying object pool

Therefore, the same BufferManager instance can be used concurrently in multiple goroutines:

```go
// Safely use BufferManager in multiple goroutines
var wg sync.WaitGroup
manager := manage.NewBufferManager()

// Start multiple goroutines
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        
        // Thread-safe: Concurrently acquire and release buffers
        buf := manager.Acquire()
        buf.WriteString(fmt.Sprintf("Goroutine %d", id))
        processData(buf.Get())
        manager.Release(buf)
    }(i)
}

wg.Wait()
```

## Performance Optimization

1. **Object Pool**: Based on the object pool mechanism of the buffer package to reduce memory allocation
2. **Automatic Reset**: Automatically resets when releasing buffers to ensure acquired buffers are in initial state
3. **Zero-copy**: Maintains consistency with the zero-copy feature of the buffer package

## Relationship with Other Components

1. **Depends on buffer package**: Uses buffer.ObjectPool for actual buffer management
2. **Used by router package**: The router package accesses buffer management functionality through the BufferManagerAccessor interface
3. **Unified Resource Management**: Provides a unified buffer management interface for the entire message routing framework
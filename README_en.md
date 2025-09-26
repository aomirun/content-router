# Content Router - High-Performance Router

Content Router is a high-performance router framework designed for Go language. It focuses on content-based routing processing, complementing Go 1.22's HTTP path routing. Through fine-grained interface design and object pool technology, it achieves efficient memory management and routing processing.

## Table of Contents

- [Features](#features)
- [Package Structure](#package-structure)
- [Core Interfaces](#core-interfaces)
  - [Buffer Related Interfaces](#buffer-related-interfaces)
  - [Context Related Interfaces](#context-related-interfaces)
  - [Router Related Interfaces](#router-related-interfaces)
- [Usage Examples](#usage-examples)
- [Advantages of Fine-Grained Interfaces](#advantages-of-fine-grained-interfaces)
- [Performance Optimization](#performance-optimization)
- [Testing](#testing)
- [Installation](#installation)
- [License](#license)

## Features

1. **Design**: Minimizes memory allocation and copying through Buffer interface and object pool technology
2. **Fine-Grained Interfaces**: Breaks functionality into smaller interfaces for better interface isolation and composability
3. **High Performance**: Uses object pools and efficient routing matching algorithms
4. **Easy Testing**: Each interface can be independently tested and mocked
5. **Backward Compatibility**: Maintains compatibility with older versions
6. **Content-Based Routing**: Focuses on content-based routing, complementing Go 1.22's HTTP path routing

## Package Structure

```
├── buffer           # Buffer management
├── context          # Context management
├── manage           # Resource management
├── middleware       # Middleware
├── router           # Router core
└── examples         # Usage examples
    ├── simple       # Simple example
    ├── finegrained  # Fine-grained interface example
    ├── middleware   # Middleware example
    └── http         # HTTP server example
```

Content Router focuses on content-based routing processing, complementing Go 1.22's HTTP path routing to provide a more comprehensive routing solution for Go applications.

## Core Interfaces

### Buffer Related Interfaces
- `Readable` - Readable buffer interface
- `Writable` - Writable buffer interface
- `Mutable` - Mutable buffer interface
- `Sliceable` - Sliceable buffer interface
- `Cloneable` - Cloneable buffer interface
- `Buffer` - Interface combining all buffer operations

### Context Related Interfaces
- `ValueStore` - Key-value storage interface
- `BufferAccessor` - Buffer access interface
- `Context` - Enhanced context interface combining standard context.Context, ValueStore, and BufferAccessor

### Router Related Interfaces
- `RouteHandler` - Route handler interface
- `RouteRegistrar` - Route registration interface
- `MiddlewareHandler` - Middleware handling interface
- `PipelineManager` - Pipeline management interface
- `ContextCreator` - Context creation interface
- `BufferManagerAccessor` - Buffer manager access interface
- `Router` - Interface combining all router functionality

## Usage Examples

For most users, it's recommended to use the root directory convenience package as the entry point, which provides unified access to all core functionality.

### Simple Usage
```go
import "github.com/aomirun/content-router"

// Create router
router := contentrouter.NewRouter()

// Register route
router.Match("Hello", func(ctx contentrouter.Context) error {
    // Processing logic
    return nil
})

// Create buffer and write data
buf := contentrouter.NewBuffer()
buf.WriteString("Hello, World!")

// Route processing
router.Route(context.Background(), buf)
```

### HTTP Server Example

This example shows how to integrate Content Router into an HTTP server. While Go 1.22 enhanced HTTP path routing capabilities, Content Router focuses on content-based routing processing, and the two can complement each other.

```go
import "github.com/aomirun/content-router"

func httpHandler(w http.ResponseWriter, r *http.Request) {
    // Create buffer
    buf := contentrouter.NewBuffer()
    buf.WriteString("HTTP request: " + r.URL.Path)
    
    // Create router
    router := contentrouter.NewRouter()
    
    // Register content-based routing rules
    // Here we route based on message content rather than URL path
    router.Match("Hello", func(ctx contentrouter.Context) error {
        response := "Processed: " + string(ctx.Buffer().Get())
        fmt.Fprintf(w, "%s", response)
        return nil
    })
    
    // Process request
    router.Route(context.Background(), buf)
}
```

### Fine-Grained Interface Usage (Advanced)

For scenarios requiring more granular control, you can directly use the sub-packages:

```go
// Directly use router package
import "github.com/aomirun/content-router/router"
import "github.com/aomirun/content-router/buffer"
import "github.com/aomirun/content-router/context"

// Create router (directly using router package)
router := router.NewRouter()

// Create buffer (directly using buffer package)
buf := buffer.NewBuffer()
buf.WriteString("Hello, World!")

// Create context (directly using context package)
ctx := context.NewContext(context.Background(), buf)

// Route processing
router.Route(ctx, buf)
```

### Middleware Usage Example

Content Router supports middleware functionality, which can be used for logging, error recovery, etc.:

```go
import "github.com/aomirun/content-router"

// Create router
router := contentrouter.NewRouter()

// Add middleware (note: middleware still needs to be imported from middleware package)
router.Use(
    // Logging middleware
    middleware.Logging(),
    // Error recovery middleware
    middleware.Recovery(),
)

// Register route
router.Match("Hello", func(ctx contentrouter.Context) error {
    fmt.Println("Processing Hello message")
    return nil
})
```

## Advantages of Fine-Grained Interfaces

1. **Better Interface Isolation** - Components only depend on the functionality they actually use
2. **Easier Testing** - Specific interfaces can be mocked rather than entire large interfaces
3. **Better Composability** - Different interfaces can be combined as needed
4. **Clearer Separation of Responsibilities** - Each interface has a clear responsibility

## Performance Optimization

1. **Object Pooling**: Uses sync.Pool to implement Buffer object pools, reducing memory allocation
2. **Avoid Data Copying**: Through Buffer interface and reference passing, avoids unnecessary data copying
3. **Efficient Route Matching**: Implements multiple matchers (prefix, suffix, contains, etc.)

## Testing

Run all tests:
```bash
go test ./...
```

### Test Coverage

The project has comprehensive test coverage to ensure the stability and reliability of components:

- `buffer` package: 100% test coverage
- `context` package: 100% test coverage
- `manage` package: 100% test coverage
- `router` package: 100% test coverage

Each package includes unit tests to verify the correctness of core functionality.

### Benchmark Testing

Run benchmark tests:
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

## Installation

```bash
go get github.com/aomirun/content-router
```

## License

MIT
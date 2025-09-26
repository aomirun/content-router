# Router Package

The router package is a core component of the content-router system, responsible for routing content based on matching rules and processing it through a chain of middleware and handlers.

## Overview

The router package implements a flexible routing mechanism that allows content to be directed to different handlers based on matching criteria. It follows a middleware pattern that enables processing content before and after the main handler execution.

### Features

- **Flexible Routing**: Route content based on custom matching rules
- **Middleware Support**: Process content through a chain of middleware functions
- **Pipeline Management**: Create isolated processing pipelines for specific routes
- **Thread Safety**: Safe for concurrent use after initialization
- **Performance Optimized**: Uses caching and object pooling for efficient processing

## Core Interfaces

### Router
The main interface that combines all routing functionality:
```go
type Router interface {
    RouteRegistrar
    MiddlewareHandler
    PipelineManager
    ContextCreator
    RouteHandler
    BufferManagerAccessor
}
```

### RouteHandler
Handles the routing process:
```go
type RouteHandler interface {
    Route(ctx context.Context, buf buffer.Buffer) (interface{}, error)
}
```

### RouteRegistrar
Manages route registration and matching:
```go
type RouteRegistrar interface {
    Register(matcher Matcher, handler Handler)
    Match(ctx router_context.Context) (Handler, bool)
}
```

### MiddlewareHandler
Manages global middleware:
```go
type MiddlewareHandler interface {
    Use(middleware MiddlewareFunc)
}
```

### PipelineManager
Creates isolated processing pipelines:
```go
type PipelineManager interface {
    Pipeline(matcher Matcher) Pipeline
}
```

### ContextCreator
Creates routing contexts:
```go
type ContextCreator interface {
    NewContext(ctx context.Context, buf buffer.Buffer) router_context.Context
}
```

### BufferManagerAccessor
Provides access to the buffer manager:
```go
type BufferManagerAccessor interface {
    GetBufferManager() manage.BufferManager
}
```

## Core Components

### Matcher
Matchers are used to determine if a route should handle a specific content. The package provides several built-in matchers:

- **PrefixMatcher**: Matches content that starts with a specific prefix
- **SuffixMatcher**: Matches content that ends with a specific suffix
- **ContainsMatcher**: Matches content that contains a specific substring

You can also create custom matchers by implementing the Matcher interface:
```go
type Matcher interface {
    Match(ctx router_context.Context) bool
}
```

### Middleware
Middleware functions allow you to process content before and after the main handler. They follow the onion model where each middleware can execute code before and after the next handler in the chain.

```go
type MiddlewareFunc func(ctx router_context.Context, next HandlerFunc) error
```

### Pipeline
Pipelines provide isolated processing chains for specific routes. They allow you to add middleware that only applies to certain routes.

```go
type Pipeline interface {
    Use(middleware MiddlewareFunc)
    Handle(ctx router_context.Context) error
}
```

### Handler
Handlers are the final destination for routed content. They perform the actual processing of the content.

```go
type Handler interface {
    Handle(ctx router_context.Context) error
}
```

## Implementation Details

### routerImpl
The main implementation of the Router interface. It maintains:
- A list of registered routes
- Global middleware functions
- A buffer manager for efficient buffer handling
- A dirty flag for cache invalidation

### Routing Process
1. When Route is called, the router iterates through registered routes
2. For each route, it checks if the matcher matches the content
3. If a match is found, it builds a handler chain with global middleware and the matched handler
4. The handler chain is executed with the content

### Handler Chain Building
The router uses a caching mechanism to avoid rebuilding the same handler chain repeatedly:
1. When middleware is added, a dirty flag is set
2. Before routing, if the dirty flag is set, the handler chain is rebuilt
3. The rebuilt chain is cached for future use

## Usage Examples

### Basic Route Registration
```go
// Create a new router
router := NewRouter()

// Register a route with a prefix matcher
router.Register(PrefixMatcher("ERROR"), HandlerFunc(func(ctx router_context.Context) error {
    // Handle error logs
    return nil
}))

// Route content
buf := buffer.NewBuffer()
buf.WriteString("ERROR: Something went wrong")
result, err := router.Route(context.Background(), buf)
```

### Using Middleware
```go
// Add global middleware
router.Use(func(ctx router_context.Context, next HandlerFunc) error {
    // Pre-processing
    fmt.Println("Before handler")
    
    // Call next middleware/handler
    err := next(ctx)
    
    // Post-processing
    fmt.Println("After handler")
    
    return err
})
```

### Using Pipelines
```go
// Create a pipeline for specific routes
matcher := PrefixMatcher("API")
pipeline := router.Pipeline(matcher)

// Add middleware to the pipeline
pipeline.Use(func(ctx router_context.Context, next HandlerFunc) error {
    // Pipeline-specific middleware
    return next(ctx)
})

// Register a handler for the pipeline
pipeline.Use(HandlerFunc(func(ctx router_context.Context) error {
    // Handle API requests
    return nil
}))
```

## Thread Safety

- Route registration and middleware addition are NOT thread-safe and should only be done during initialization
- Once initialized, the router can be safely used concurrently for routing operations
- The handler chain caching mechanism is designed to be thread-safe

## Performance Considerations

- Handler chains are cached to avoid rebuilding them for each request
- Object pooling is used for buffer management
- Lazy initialization is used where possible to defer expensive operations

## Testing

The package includes comprehensive tests covering:
- Interface implementation verification
- Route registration and matching
- Middleware execution order
- Pipeline functionality
- Handler chain caching
- Built-in matcher implementations

## Dependencies

- `buffer` package: For content buffering
- `context` package: For context management
- `manage` package: For buffer management
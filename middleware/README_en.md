# Middleware Package

[中文版本](README.md)

The middleware package provides common middleware components for the content router, implementing the onion model architecture. Middleware components can execute additional logic before and after request processing, such as logging, error recovery, etc.

## Features

- **Onion Model Architecture**: Middleware follows the onion model, where each middleware wraps the next one in the chain
- **Recovery Middleware**: Captures panics during handler execution and logs error information
- **Logging Middleware**: Records request processing time and related information
- **Easy Integration**: Simple API for registering middleware with the router
- **Custom Middleware Support**: Easy to create custom middleware following a standard pattern

## Core Interfaces

### MiddlewareFunc

The `MiddlewareFunc` type defines the middleware function signature:

```go
type MiddlewareFunc func(ctx router_context.Context, next HandlerFunc) error
```

Parameters:
- `ctx`: Request context
- `next`: The next handler in the chain

## Middleware Components

### Recovery Middleware

The `RecoveryMiddleware` captures panics during handler execution and logs error information with stack traces.

Key Features:
- Captures and logs panics with stack trace information
- Prevents application crashes due to unhandled panics
- Extensible for additional error handling logic (monitoring, logging, etc.)

Usage:
```go
r.Use(middleware.RecoveryMiddleware())
```

### Logging Middleware

The `LoggingMiddleware` records request processing time and related information.

Key Features:
- Records processing start time and duration
- Logs data preview (first 50 bytes to avoid sensitive information)
- Distinguishes between successful processing and errors
- Measures processing time for performance monitoring

Usage:
```go
r.Use(middleware.LoggingMiddleware())
```

## Usage Example

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/aomirun/content-router"
    "github.com/aomirun/content-router/middleware"
)

func main() {
    // Create router instance
    r := contentrouter.NewRouter()
    
    // Register middleware
    r.Use(middleware.LoggingMiddleware())
    r.Use(middleware.RecoveryMiddleware())
    
    // Register route handler
    r.Match("test", func(ctx contentrouter.Context) error {
        // Simulate processing time
        time.Sleep(100 * time.Millisecond)
        
        // Get request data
        data := ctx.Buffer().Get()
        fmt.Printf("Processing data: %s\n", string(data))
        
        // Simulate panic scenario
        if string(data) == "trigger_panic" {
            panic("Simulated panic for demonstration")
        }
        
        // Simulate processing error
        if string(data) == "trigger_error" {
            return fmt.Errorf("simulated processing error")
        }
        
        // Normal processing completion
        fmt.Println("Data processed successfully")
        return nil
    })
    
    // Create test data
    buf := contentrouter.NewBuffer()
    buf.Write([]byte("test_data"))
    
    // Process data
    _, err := r.Route(context.Background(), buf)
    if err != nil {
        fmt.Printf("Error processing data: %v\n", err)
    }
}
```

## Creating Custom Middleware

To create custom middleware, follow this pattern:

```go
func CustomMiddleware() router.MiddlewareFunc {
    return func(ctx router_context.Context, next router.HandlerFunc) error {
        // Pre-processing logic
        
        // Call the next middleware/handler
        err := next(ctx)
        
        // Post-processing logic
        
        return err
    }
}
```

Middleware functions receive:
1. `ctx` - Request context
2. `next` - The next handler in the chain

To continue the chain, call `next(ctx)`. To abort the chain, return an error without calling `next`.

## Thread Safety

All middleware components are designed to be thread-safe and can be used concurrently across multiple goroutines.

## Performance Optimization

- Minimal overhead: Middleware adds minimal performance overhead
- Efficient logging: Only logs first 50 bytes of data to avoid excessive output
- Stack trace capture: Efficiently captures stack traces only when panics occur

## Testing

The middleware package includes the following test functions:

1. `TestLoggingMiddleware` - Tests logging middleware under normal conditions
2. `TestLoggingMiddlewareWithError` - Tests logging middleware with errors
3. `TestRecoveryMiddleware` - Tests error recovery middleware
4. `TestRecoveryMiddlewareWithoutPanic` - Tests error recovery middleware without panics
5. `TestLoggingMiddlewareWithLongData` - Tests logging middleware handling long data

Run tests with the following command:

```bash
go test -v ./middleware/...
```

Test coverage reaches 100%.
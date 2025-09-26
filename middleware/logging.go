package middleware

import (
	"fmt"
	"time"

	router_context "github.com/aomirun/content-router/context"
	"github.com/aomirun/content-router/router"
)

// LoggingMiddleware 创建一个日志记录中间件
// 该中间件会记录请求的处理时间和其他相关信息
func LoggingMiddleware() router.MiddlewareFunc {
	return func(ctx router_context.Context, next router.HandlerFunc) error {
		// 记录开始时间
		start := time.Now()

		// 获取请求数据的前几个字节用于日志记录（避免记录敏感信息）
		data := ctx.Buffer().Get()
		dataPreview := ""
		if len(data) > 50 {
			dataPreview = string(data[:50]) + "..."
		} else {
			dataPreview = string(data)
		}

		// 记录请求开始
		fmt.Printf("Starting processing at %v, data preview: %s\n", start, dataPreview)

		// 执行下一个处理器
		err := next(ctx)

		// 记录结束时间和处理结果
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("Processing failed after %v, error: %v\n", duration, err)
		} else {
			fmt.Printf("Processing completed in %v\n", duration)
		}

		return err
	}
}

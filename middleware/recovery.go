package middleware

import (
	"fmt"
	"runtime"

	router_context "github.com/aomirun/content-router/context"
	"github.com/aomirun/content-router/router"
)

// RecoveryMiddleware 创建一个错误恢复中间件
// 该中间件会捕获处理器执行过程中的panic，并记录错误信息
func RecoveryMiddleware() router.MiddlewareFunc {
	return func(ctx router_context.Context, next router.HandlerFunc) error {
		defer func() {
			if err := recover(); err != nil {
				// 获取panic时的堆栈信息
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, false)]

				// 记录错误信息和堆栈
				fmt.Printf("Recovery middleware caught panic: %v\nStack: %s\n", err, stack)

				// 可以在这里添加更多的错误处理逻辑，比如：
				// 1. 发送错误报告到监控系统
				// 2. 记录到日志文件
				// 3. 返回统一的错误响应格式
			}
		}()

		// 执行下一个处理器
		return next(ctx)
	}
}

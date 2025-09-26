package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/aomirun/content-router/api"
	"github.com/aomirun/content-router/buffer"
	"github.com/aomirun/content-router/router"
)

func BenchmarkRouter_Route(b *testing.B) {
	// 创建路由器
	router := api.NewRouter()

	// 注册路由
	router.Match("Hello", func(ctx api.Context) error {
		// 简单处理逻辑
		return nil
	})

	// 创建缓冲区
	buf := api.NewBuffer()
	buf.WriteString("Hello, World!")

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := router.Route(context.Background(), buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRouter_RouteWithMiddleware(b *testing.B) {
	// 创建路由器
	router := api.NewRouter()

	// 添加中间件
	router.Use(func(ctx api.Context, next api.HandlerFunc) error {
		// 简单中间件逻辑
		return next(ctx)
	})

	// 注册路由
	router.Match("Hello", func(ctx api.Context) error {
		// 简单处理逻辑
		return nil
	})

	// 创建缓冲区
	buf := api.NewBuffer()
	buf.WriteString("Hello, World!")

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := router.Route(context.Background(), buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuffer_AcquireRelease(b *testing.B) {
	// 创建路由器以获取BufferManager
	router := api.NewRouter()
	bufferManager := router.BufferManager()

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		buf := bufferManager.Acquire()
		bufferManager.Release(buf)
	}
}

func BenchmarkContext_ValueStore(b *testing.B) {
	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	ctx := api.NewContext(context.Background(), buf)

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		ctx.Set("key", "value")
		_ = ctx.Get("key")
	}
}

func BenchmarkContextPool_AcquireRelease(b *testing.B) {
	// 创建缓冲区
	buf := api.NewBuffer()

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		ctx := api.NewContext(context.Background(), buf)
		// 释放上下文到池中
		if c, ok := ctx.(interface{ Reset() }); ok {
			c.Reset()
		}
	}
}

func BenchmarkContext_Methods(b *testing.B) {
	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	ctx := api.NewContext(context.Background(), buf)

	// 设置测试数据
	ctx.Set("string_key", "test_string")
	ctx.Set("int_key", 42)
	ctx.Set("bool_key", true)
	ctx.Set("float_key", 3.14)
	ctx.Set("bytes_key", []byte("test_bytes"))
	ctx.Set("time_key", time.Now())

	// 重置计时器
	b.ResetTimer()

	// 测试各种Context方法
	for i := 0; i < b.N; i++ {
		// 测试Get方法
		_ = ctx.Get("string_key")

		// 测试类型特定的Get方法
		_, _ = ctx.GetString("string_key")
		_, _ = ctx.GetInt("int_key")
		_, _ = ctx.GetBool("bool_key")
		_, _ = ctx.GetFloat64("float_key")
		_, _ = ctx.GetBytes("bytes_key")
		_, _ = ctx.GetTime("time_key")
	}
}

func BenchmarkContext_Keys(b *testing.B) {
	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	ctx := api.NewContext(context.Background(), buf)

	// 设置测试数据
	ctx.Set("string_key", "test_string")
	ctx.Set("int_key", 42)
	ctx.Set("bool_key", true)
	ctx.Set("float_key", 3.14)
	ctx.Set("bytes_key", []byte("test_bytes"))
	ctx.Set("time_key", time.Now())

	// 重置计时器
	b.ResetTimer()

	// 测试Keys方法
	for i := 0; i < b.N; i++ {
		_ = ctx.Keys()
	}
}

func BenchmarkContext_Methods_WithoutKeys(b *testing.B) {
	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	ctx := api.NewContext(context.Background(), buf)

	// 设置测试数据
	ctx.Set("string_key", "test_string")
	ctx.Set("int_key", 42)
	ctx.Set("bool_key", true)
	ctx.Set("float_key", 3.14)
	ctx.Set("bytes_key", []byte("test_bytes"))
	ctx.Set("time_key", time.Now())

	// 重置计时器
	b.ResetTimer()

	// 测试各种Context方法（不包括Keys）
	for i := 0; i < b.N; i++ {
		// 测试Get方法
		_ = ctx.Get("string_key")

		// 测试类型特定的Get方法
		_, _ = ctx.GetString("string_key")
		_, _ = ctx.GetInt("int_key")
		_, _ = ctx.GetBool("bool_key")
		_, _ = ctx.GetFloat64("float_key")
		_, _ = ctx.GetBytes("bytes_key")
		_, _ = ctx.GetTime("time_key")
	}
}

func BenchmarkMatcher_Prefix(b *testing.B) {
	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	buf.WriteString("Hello, World!")
	ctx := api.NewContext(context.Background(), buf)

	// 创建匹配器
	matcher := router.PrefixMatcher("Hello")

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_ = matcher.Match(ctx)
	}
}

func BenchmarkMatcher_Suffix(b *testing.B) {
	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	buf.WriteString("Hello, World!")
	ctx := api.NewContext(context.Background(), buf)

	// 创建匹配器
	matcher := router.SuffixMatcher("World!")

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_ = matcher.Match(ctx)
	}
}

func BenchmarkMatcher_Contains(b *testing.B) {
	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	buf.WriteString("Hello, World!")
	ctx := api.NewContext(context.Background(), buf)

	// 创建匹配器
	matcher := router.ContainsMatcher("World")

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_ = matcher.Match(ctx)
	}
}

func BenchmarkPipeline_WithMiddleware(b *testing.B) {
	// 创建路由器
	r := api.NewRouter()

	// 创建匹配器
	matcher := router.PrefixMatcher("test")

	// 创建管道
	pipeline := r.Pipeline(matcher)

	// 添加中间件
	pipeline.Use(func(ctx api.Context, next api.HandlerFunc) error {
		return next(ctx)
	})

	pipeline.Use(func(ctx api.Context, next api.HandlerFunc) error {
		return next(ctx)
	})

	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	buf.WriteString("test data")
	ctx := api.NewContext(context.Background(), buf)

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_ = pipeline.Handle(ctx)
	}
}

func BenchmarkPipeline_WithoutMiddleware(b *testing.B) {
	// 创建路由器
	r := api.NewRouter()

	// 创建匹配器
	matcher := router.PrefixMatcher("test")

	// 创建管道
	pipeline := r.Pipeline(matcher)

	// 创建缓冲区和上下文
	buf := api.NewBuffer()
	buf.WriteString("test data")
	ctx := api.NewContext(context.Background(), buf)

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_ = pipeline.Handle(ctx)
	}
}

func BenchmarkBufferPool_AcquireRelease(b *testing.B) {
	// 创建缓冲区池
	pool := buffer.NewPool()

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		buf := pool.Acquire()
		pool.Release(buf)
	}
}

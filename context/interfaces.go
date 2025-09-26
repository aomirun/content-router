package context

import (
	"context"
	"time"

	"github.com/aomirun/content-router/buffer"
)

// ValueStore 定义键值存储接口
type ValueStore interface {
	// Set 设置键值对
	Set(key, value interface{})

	// Get 获取值
	Get(key interface{}) interface{}

	// GetString 获取字符串值
	GetString(key interface{}) (string, bool)

	// GetInt 获取整数值
	GetInt(key interface{}) (int, bool)

	// GetInt64 获取64位整数值
	GetInt64(key interface{}) (int64, bool)

	// GetBool 获取布尔值
	GetBool(key interface{}) (bool, bool)

	// GetFloat64 获取浮点数值
	GetFloat64(key interface{}) (float64, bool)

	// GetBytes 获取字节数组
	GetBytes(key interface{}) ([]byte, bool)

	// GetTime 获取时间值
	GetTime(key interface{}) (time.Time, bool)

	// Delete 删除键值对
	Delete(key interface{})

	// Keys 获取所有键
	Keys() []interface{}
}

// BufferAccessor 定义缓冲区访问接口
type BufferAccessor interface {
	// Buffer 获取与上下文关联的缓冲区
	Buffer() buffer.Buffer
}

// Context 定义增强的上下文接口
// 它组合了标准context.Context、ValueStore和BufferAccessor接口
type Context interface {
	context.Context
	ValueStore
	BufferAccessor

	// Fork 创建上下文的副本，但共享相同的缓冲区
	Fork() Context

	// ForkWithBuffer 创建上下文的副本，并使用新的缓冲区
	ForkWithBuffer(buffer buffer.Buffer) Context
}

package context

import (
	"context"
	"sync"
	"time"

	"github.com/aomirun/content-router/buffer"
)

// contextImpl 是Context接口的具体实现
type contextImpl struct {
	context.Context
	buffer buffer.Buffer
	values map[interface{}]interface{}
}

// contextPool 是contextImpl的对象池
var contextPool = sync.Pool{
	New: func() interface{} {
		return &contextImpl{
			values: make(map[interface{}]interface{}),
		}
	},
}

// NewContext 创建一个新的上下文实例
func NewContext(parent context.Context, buf buffer.Buffer) Context {
	// 如果父上下文为nil，则使用Background上下文
	if parent == nil {
		parent = context.Background()
	}

	// 从对象池获取contextImpl
	ctx := contextPool.Get().(*contextImpl)
	ctx.Context = parent
	ctx.buffer = buf
	// 清空values map
	for k := range ctx.values {
		delete(ctx.values, k)
	}

	return ctx
}

// Reset 重置上下文，将其放回对象池
func (c *contextImpl) Reset() {
	// 清空values map
	for k := range c.values {
		delete(c.values, k)
	}
	c.buffer = nil
	c.Context = nil
	contextPool.Put(c)
}

// Set 设置键值对
func (c *contextImpl) Set(key, value interface{}) {
	c.values[key] = value
}

// Get 获取值
func (c *contextImpl) Get(key interface{}) interface{} {
	return c.values[key]
}

// GetString 获取字符串值
func (c *contextImpl) GetString(key interface{}) (string, bool) {
	val := c.Get(key)
	if str, ok := val.(string); ok {
		return str, true
	}
	return "", false
}

// GetInt 获取整数值
func (c *contextImpl) GetInt(key interface{}) (int, bool) {
	val := c.Get(key)
	if i, ok := val.(int); ok {
		return i, true
	}
	return 0, false
}

// GetInt64 获取64位整数值
func (c *contextImpl) GetInt64(key interface{}) (int64, bool) {
	val := c.Get(key)
	if i, ok := val.(int64); ok {
		return i, true
	}
	return 0, false
}

// GetBool 获取布尔值
func (c *contextImpl) GetBool(key interface{}) (bool, bool) {
	val := c.Get(key)
	if b, ok := val.(bool); ok {
		return b, true
	}
	return false, false
}

// GetFloat64 获取浮点数值
func (c *contextImpl) GetFloat64(key interface{}) (float64, bool) {
	val := c.Get(key)
	if f, ok := val.(float64); ok {
		return f, true
	}
	return 0.0, false
}

// GetBytes 获取字节数组
func (c *contextImpl) GetBytes(key interface{}) ([]byte, bool) {
	val := c.Get(key)
	if b, ok := val.([]byte); ok {
		return b, true
	}
	return nil, false
}

// GetTime 获取时间值
func (c *contextImpl) GetTime(key interface{}) (time.Time, bool) {
	val := c.Get(key)
	if t, ok := val.(time.Time); ok {
		return t, true
	}
	return time.Time{}, false
}

// Delete 删除键值对
func (c *contextImpl) Delete(key interface{}) {
	delete(c.values, key)
}

// Keys 获取所有键
func (c *contextImpl) Keys() []interface{} {
	keys := make([]interface{}, 0, len(c.values))
	for k := range c.values {
		keys = append(keys, k)
	}
	return keys
}

// Buffer 获取与上下文关联的缓冲区
func (c *contextImpl) Buffer() buffer.Buffer {
	return c.buffer
}

// Fork 创建上下文的副本，但共享相同的缓冲区
func (c *contextImpl) Fork() Context {
	// 复制values map
	values := make(map[interface{}]interface{}, len(c.values))
	for k, v := range c.values {
		values[k] = v
	}

	return &contextImpl{
		Context: c.Context,
		buffer:  c.buffer,
		values:  values,
	}
}

// ForkWithBuffer 创建上下文的副本，并使用新的缓冲区
func (c *contextImpl) ForkWithBuffer(buf buffer.Buffer) Context {
	// 复制values map
	values := make(map[interface{}]interface{}, len(c.values))
	for k, v := range c.values {
		values[k] = v
	}

	return &contextImpl{
		Context: c.Context,
		buffer:  buf,
		values:  values,
	}
}

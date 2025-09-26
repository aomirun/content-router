package context

import (
	"context"
	"testing"
	"time"

	"github.com/aomirun/content-router/buffer"
)

func TestContextImpl(t *testing.T) {
	// 创建一个buffer
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 创建一个context
	ctx := NewContext(context.Background(), buf)

	// 测试Buffer方法
	if ctx.Buffer() == nil {
		t.Error("Buffer() returned nil")
	}

	// 测试ValueStore功能
	ctx.Set("key1", "value1")
	ctx.Set("key2", 42)
	ctx.Set("key3", true)
	ctx.Set("key4", int64(1234567890))
	ctx.Set("key5", float64(3.14159))
	ctx.Set("key6", []byte("byte data"))
	ctx.Set("key7", time.Now())

	// 测试Get方法
	if val := ctx.Get("key1"); val != "value1" {
		t.Errorf("Get(key1) returned %v, expected value1", val)
	}

	// 测试GetString方法
	if val, ok := ctx.GetString("key1"); !ok || val != "value1" {
		t.Errorf("GetString(key1) returned %v, %v, expected value1, true", val, ok)
	}

	// 测试GetString方法 - 错误类型
	if val, ok := ctx.GetString("key2"); ok || val != "" {
		t.Errorf("GetString(key2) returned %v, %v, expected '', false", val, ok)
	}

	// 测试GetInt方法
	if val, ok := ctx.GetInt("key2"); !ok || val != 42 {
		t.Errorf("GetInt(key2) returned %v, %v, expected 42, true", val, ok)
	}

	// 测试GetInt方法 - 错误类型
	if val, ok := ctx.GetInt("key1"); ok || val != 0 {
		t.Errorf("GetInt(key1) returned %v, %v, expected 0, false", val, ok)
	}

	// 测试GetBool方法
	if val, ok := ctx.GetBool("key3"); !ok || val != true {
		t.Errorf("GetBool(key3) returned %v, %v, expected true, true", val, ok)
	}

	// 测试GetBool方法 - 错误类型
	if val, ok := ctx.GetBool("key1"); ok || val != false {
		t.Errorf("GetBool(key1) returned %v, %v, expected false, false", val, ok)
	}

	// 测试GetInt64方法
	if val, ok := ctx.GetInt64("key4"); !ok || val != 1234567890 {
		t.Errorf("GetInt64(key4) returned %v, %v, expected 1234567890, true", val, ok)
	}

	// 测试GetInt64方法 - 错误类型
	if val, ok := ctx.GetInt64("key1"); ok || val != 0 {
		t.Errorf("GetInt64(key1) returned %v, %v, expected 0, false", val, ok)
	}

	// 测试GetFloat64方法
	if val, ok := ctx.GetFloat64("key5"); !ok || val != 3.14159 {
		t.Errorf("GetFloat64(key5) returned %v, %v, expected 3.14159, true", val, ok)
	}

	// 测试GetFloat64方法 - 错误类型
	if val, ok := ctx.GetFloat64("key1"); ok || val != 0.0 {
		t.Errorf("GetFloat64(key1) returned %v, %v, expected 0.0, false", val, ok)
	}

	// 测试GetBytes方法
	if val, ok := ctx.GetBytes("key6"); !ok || string(val) != "byte data" {
		t.Errorf("GetBytes(key6) returned %v, %v, expected byte data, true", val, ok)
	}

	// 测试GetBytes方法 - 错误类型
	if val, ok := ctx.GetBytes("key1"); ok || val != nil {
		t.Errorf("GetBytes(key1) returned %v, %v, expected nil, false", val, ok)
	}

	// 测试GetTime方法
	now := time.Now()
	ctx.Set("timeKey", now)
	if val, ok := ctx.GetTime("timeKey"); !ok || val != now {
		t.Errorf("GetTime(timeKey) returned %v, %v, expected %v, true", val, ok, now)
	}

	// 测试GetTime方法 - 错误类型
	if val, ok := ctx.GetTime("key1"); ok || !val.IsZero() {
		t.Errorf("GetTime(key1) returned %v, %v, expected zero time, false", val, ok)
	}

	// 测试Delete方法
	ctx.Delete("key1")
	if val := ctx.Get("key1"); val != nil {
		t.Errorf("Get(key1) after Delete returned %v, expected nil", val)
	}

	// 测试Keys方法
	keys := ctx.Keys()
	if len(keys) != 7 {
		t.Errorf("Keys() returned %d keys, expected 7", len(keys))
	}
}

func TestContextFork(t *testing.T) {
	// 创建一个buffer
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 创建一个context
	parent := NewContext(context.Background(), buf)
	parent.Set("parentKey", "parentValue")

	// 测试Fork
	child := parent.Fork()

	// 子context应该有父context的值
	if val, ok := child.GetString("parentKey"); !ok || val != "parentValue" {
		t.Errorf("Forked context missing parent value: %v, %v", val, ok)
	}

	// 修改子context不应该影响父context
	child.Set("childKey", "childValue")
	if val := parent.Get("childKey"); val != nil {
		t.Errorf("Parent context was affected by child modification: %v", val)
	}

	// 测试ForkWithBuffer
	newBuf := buffer.NewBuffer()
	newBuf.WriteString("new data")
	newChild := parent.ForkWithBuffer(newBuf)

	// 新的context应该有新的buffer
	if newChild.Buffer().Len() != newBuf.Len() {
		t.Errorf("ForkWithBuffer did not use the provided buffer")
	}

	// 新的context应该有父context的值
	if val, ok := newChild.GetString("parentKey"); !ok || val != "parentValue" {
		t.Errorf("ForkWithBuffer context missing parent value: %v, %v", val, ok)
	}
}

func TestContextReset(t *testing.T) {
	// 创建一个buffer
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 创建一个context
	ctx := NewContext(context.Background(), buf)
	ctx.Set("key1", "value1")

	// 验证context中有数据
	if val := ctx.Get("key1"); val != "value1" {
		t.Errorf("Get(key1) returned %v, expected value1", val)
	}

	// 重置context（通过类型断言调用Reset方法）
	if c, ok := ctx.(*contextImpl); ok {
		c.Reset()
	} else {
		t.Error("Failed to type assert context to *contextImpl")
	}

	// 注意：Reset会将context放回池中，所以我们不能再使用它进行测试
	// 我们只需要确保Reset方法能够执行而不报错
}

func TestNewContextWithNilParent(t *testing.T) {
	// 创建一个buffer
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 使用nil作为父context创建context
	ctx := NewContext(nil, buf)

	// 验证context不为nil
	if ctx == nil {
		t.Error("NewContext with nil parent returned nil")
	}

	// 验证context有默认的background context
	if ctx.Value("nonexistent") != nil {
		t.Error("Context should have background context as parent")
	}
}

func TestNewContextWithExistingContext(t *testing.T) {
	// 创建一个buffer
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 创建第一个context并重置它，使其回到池中
	ctx1 := NewContext(context.Background(), buf)
	ctx1.Set("key1", "value1")

	// 通过类型断言调用Reset方法
	if c, ok := ctx1.(*contextImpl); ok {
		c.Reset()
	} else {
		t.Error("Failed to type assert context to *contextImpl")
	}

	// 创建第二个context，应该从池中获取
	parent := context.WithValue(context.Background(), "parentKey", "parentValue")
	ctx2 := NewContext(parent, buf)

	// 验证context不为nil
	if ctx2 == nil {
		t.Error("NewContext returned nil")
	}

	// 验证父context被正确设置
	if val := ctx2.Value("parentKey"); val != "parentValue" {
		t.Errorf("Parent context value not set correctly: got %v, expected parentValue", val)
	}

	// 验证values map是空的（从池中获取的context应该被清空）
	if len(ctx2.Keys()) != 0 {
		t.Errorf("Context from pool should have empty values map, got %d keys", len(ctx2.Keys()))
	}
}

func TestNewContextPoolGetAndClear(t *testing.T) {
	// 创建一个buffer
	buf := buffer.NewBuffer()
	buf.WriteString("test data")

	// 创建一个context并设置一些值
	ctx1 := NewContext(context.Background(), buf)
	ctx1.Set("key1", "value1")
	ctx1.Set("key2", 42)

	// 验证context中有值
	if len(ctx1.Keys()) != 2 {
		t.Errorf("Context should have 2 keys, got %d", len(ctx1.Keys()))
	}

	// 重置context，将其放回池中
	if c, ok := ctx1.(*contextImpl); ok {
		c.Reset()
	} else {
		t.Error("Failed to type assert context to *contextImpl")
	}

	// 从池中获取一个新的context
	ctx2 := NewContext(context.Background(), buf)

	// 验证新context的values map是空的
	if len(ctx2.Keys()) != 0 {
		t.Errorf("New context from pool should have empty values map, got %d keys", len(ctx2.Keys()))
	}

	// 验证新context的buffer是正确的
	if ctx2.Buffer() != buf {
		t.Error("New context should have the correct buffer")
	}
}

func TestNewContextClearValuesMap(t *testing.T) {
	// 创建一个buffer
	buf := buffer.NewBuffer()

	// 直接创建一个contextImpl并设置一些值
	ctxImpl := &contextImpl{
		Context: context.Background(),
		buffer:  buf,
		values:  make(map[interface{}]interface{}),
	}
	ctxImpl.values["key1"] = "value1"
	ctxImpl.values["key2"] = 42

	// 将contextImpl放回池中
	contextPool.Put(ctxImpl)

	// 从池中获取context，这应该触发清空values map的代码
	ctx := NewContext(context.Background(), buf)

	// 验证values map是空的
	if len(ctx.Keys()) != 0 {
		t.Errorf("New context from pool should have empty values map, got %d keys", len(ctx.Keys()))
	}
}

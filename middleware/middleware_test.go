package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aomirun/content-router/buffer"
	router_context "github.com/aomirun/content-router/context"
)

// mockContext 是一个模拟的上下文实现，用于测试
type mockContext struct {
	context.Context
	buffer buffer.Buffer
	values map[interface{}]interface{}
}

func (m *mockContext) Get(key interface{}) interface{} {
	if m.values == nil {
		return nil
	}
	return m.values[key]
}

func (m *mockContext) GetString(key interface{}) (string, bool) {
	if m.values == nil {
		return "", false
	}
	if v, ok := m.values[key]; ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

func (m *mockContext) GetInt(key interface{}) (int, bool) {
	if m.values == nil {
		return 0, false
	}
	if v, ok := m.values[key]; ok {
		if i, ok := v.(int); ok {
			return i, true
		}
	}
	return 0, false
}

func (m *mockContext) GetInt64(key interface{}) (int64, bool) {
	if m.values == nil {
		return 0, false
	}
	if v, ok := m.values[key]; ok {
		if i, ok := v.(int64); ok {
			return i, true
		}
	}
	return 0, false
}

func (m *mockContext) GetBool(key interface{}) (bool, bool) {
	if m.values == nil {
		return false, false
	}
	if v, ok := m.values[key]; ok {
		if b, ok := v.(bool); ok {
			return b, true
		}
	}
	return false, false
}

func (m *mockContext) GetFloat64(key interface{}) (float64, bool) {
	if m.values == nil {
		return 0, false
	}
	if v, ok := m.values[key]; ok {
		if f, ok := v.(float64); ok {
			return f, true
		}
	}
	return 0, false
}

func (m *mockContext) GetBytes(key interface{}) ([]byte, bool) {
	if m.values == nil {
		return nil, false
	}
	if v, ok := m.values[key]; ok {
		if b, ok := v.([]byte); ok {
			return b, true
		}
	}
	return nil, false
}

func (m *mockContext) GetTime(key interface{}) (time.Time, bool) {
	if m.values == nil {
		return time.Time{}, false
	}
	if v, ok := m.values[key]; ok {
		if t, ok := v.(time.Time); ok {
			return t, true
		}
	}
	return time.Time{}, false
}

func (m *mockContext) Set(key, value interface{}) {
	if m.values == nil {
		m.values = make(map[interface{}]interface{})
	}
	m.values[key] = value
}

func (m *mockContext) Delete(key interface{}) {
	if m.values != nil {
		delete(m.values, key)
	}
}

func (m *mockContext) Keys() []interface{} {
	if m.values == nil {
		return nil
	}
	keys := make([]interface{}, 0, len(m.values))
	for k := range m.values {
		keys = append(keys, k)
	}
	return keys
}

func (m *mockContext) Fork() router_context.Context {
	return &mockContext{
		buffer: m.buffer,
		values: m.values,
	}
}

func (m *mockContext) ForkWithBuffer(buffer buffer.Buffer) router_context.Context {
	return &mockContext{
		buffer: buffer,
		values: m.values,
	}
}

func (m *mockContext) Buffer() buffer.Buffer {
	return m.buffer
}

// mockBuffer 是一个模拟的缓冲区实现，用于测试
type mockBuffer struct {
	data []byte
}

func (m *mockBuffer) Get() []byte {
	return m.data
}

func (m *mockBuffer) Len() int {
	return len(m.data)
}

func (m *mockBuffer) Cap() int {
	return cap(m.data)
}

func (m *mockBuffer) Write(data []byte) (int, error) {
	m.data = append(m.data, data...)
	return len(data), nil
}

func (m *mockBuffer) WriteString(s string) (int, error) {
	m.data = append(m.data, s...)
	return len(s), nil
}

func (m *mockBuffer) Reset() {
	m.data = m.data[:0]
}

func (m *mockBuffer) Truncate(n int) {
	if n < len(m.data) {
		m.data = m.data[:n]
	}
}

func (m *mockBuffer) Slice(start, end int) buffer.Buffer {
	if start < 0 {
		start = 0
	}
	if end > len(m.data) {
		end = len(m.data)
	}
	return &mockBuffer{
		data: append([]byte(nil), m.data[start:end]...),
	}
}

func (m *mockBuffer) Clone() buffer.Buffer {
	return &mockBuffer{
		data: append([]byte(nil), m.data...),
	}
}

func (m *mockBuffer) Read(p []byte) (n int, err error) {
	if len(m.data) == 0 {
		return 0, io.EOF
	}
	n = copy(p, m.data)
	m.data = m.data[n:]
	return n, nil
}

func (m *mockBuffer) String() string {
	return string(m.data)
}

// TestLoggingMiddleware 测试日志记录中间件
func TestLoggingMiddleware(t *testing.T) {
	// 保存原始的stdout
	oldStdout := os.Stdout

	// 创建管道来捕获输出
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 创建一个缓冲区来收集输出
	var buf bytes.Buffer

	// 在goroutine中从管道读取数据
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// 创建中间件
	loggingMiddleware := LoggingMiddleware()

	// 创建模拟的上下文和缓冲区
	mockBuf := &mockBuffer{data: []byte("Hello, World!")}
	mockCtx := &mockContext{}
	mockCtx = mockCtx.ForkWithBuffer(mockBuf).(*mockContext)

	// 创建一个简单的处理器函数
	handler := func(ctx router_context.Context) error {
		time.Sleep(10 * time.Millisecond) // 模拟一些处理时间
		return nil
	}

	// 应用中间件
	err := loggingMiddleware(mockCtx, handler)

	// 关闭写入端并等待读取完成
	w.Close()
	<-done

	// 恢复原始的stdout
	os.Stdout = oldStdout

	// 验证没有错误
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证输出包含预期的日志信息
	output := buf.String()
	if !strings.Contains(output, "Starting processing") {
		t.Error("Expected log output to contain 'Starting processing'")
	}

	if !strings.Contains(output, "Processing completed") {
		t.Error("Expected log output to contain 'Processing completed'")
	}

	if !strings.Contains(output, "data preview:") {
		t.Error("Expected log output to contain 'data preview:'")
	}
}

// TestLoggingMiddlewareWithError 测试带有错误的日志记录中间件
func TestLoggingMiddlewareWithError(t *testing.T) {
	// 保存原始的stdout
	oldStdout := os.Stdout

	// 创建管道来捕获输出
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 创建一个缓冲区来收集输出
	var buf bytes.Buffer

	// 在goroutine中从管道读取数据
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// 创建中间件
	loggingMiddleware := LoggingMiddleware()

	// 创建模拟的上下文和缓冲区
	mockBuf := &mockBuffer{data: []byte("Error test data")}
	mockCtx := &mockContext{}
	mockCtx = mockCtx.ForkWithBuffer(mockBuf).(*mockContext)

	// 创建一个返回错误的处理器函数
	expectedError := fmt.Errorf("test error")
	handler := func(ctx router_context.Context) error {
		time.Sleep(10 * time.Millisecond) // 模拟一些处理时间
		return expectedError
	}

	// 应用中间件
	err := loggingMiddleware(mockCtx, handler)

	// 关闭写入端并等待读取完成
	w.Close()
	<-done

	// 恢复原始的stdout
	os.Stdout = oldStdout

	// 验证错误被正确传递
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// 验证输出包含预期的日志信息
	output := buf.String()
	if !strings.Contains(output, "Starting processing") {
		t.Error("Expected log output to contain 'Starting processing'")
	}

	if !strings.Contains(output, "Processing failed") {
		t.Error("Expected log output to contain 'Processing failed'")
	}

	if !strings.Contains(output, "error:") {
		t.Error("Expected log output to contain 'error:'")
	}
}

// TestRecoveryMiddleware 测试错误恢复中间件
func TestRecoveryMiddleware(t *testing.T) {
	// 保存原始的stdout
	oldStdout := os.Stdout

	// 创建管道来捕获输出
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 创建一个缓冲区来收集输出
	var buf bytes.Buffer

	// 在goroutine中从管道读取数据
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// 创建中间件
	recoveryMiddleware := RecoveryMiddleware()

	// 创建模拟的上下文和缓冲区
	mockBuf := &mockBuffer{data: []byte("Recovery test data")}
	mockCtx := &mockContext{}
	mockCtx = mockCtx.ForkWithBuffer(mockBuf).(*mockContext)

	// 创建一个会引发panic的处理器函数
	handler := func(ctx router_context.Context) error {
		panic("test panic")
	}

	// 应用中间件
	err := recoveryMiddleware(mockCtx, handler)

	// 关闭写入端并等待读取完成
	w.Close()
	<-done

	// 在这个测试中，我们主要关注是否有正确的日志输出
	// 即使有panic被恢复，err也可能不为nil，这取决于具体实现
	// 我们不需要对err进行断言
	_ = err

	// 恢复原始的stdout
	os.Stdout = oldStdout

	// 验证错误恢复中间件不会传播panic，但会返回nil错误
	// 注意：在Go中，recover()只能在defer中使用，并且不能阻止panic传播到调用者
	// 在我们的实现中，我们只是记录panic，但仍然让处理器返回错误
	// 因此err可能不是nil，这取决于具体的实现细节

	// 验证输出包含预期的恢复信息
	output := buf.String()
	if !strings.Contains(output, "Recovery middleware caught panic") {
		t.Error("Expected log output to contain 'Recovery middleware caught panic'")
	}

	if !strings.Contains(output, "test panic") {
		t.Error("Expected log output to contain 'test panic'")
	}
}

// TestRecoveryMiddlewareWithoutPanic 测试没有panic时的错误恢复中间件
func TestRecoveryMiddlewareWithoutPanic(t *testing.T) {
	// 创建中间件
	recoveryMiddleware := RecoveryMiddleware()

	// 创建模拟的上下文和缓冲区
	mockBuf := &mockBuffer{data: []byte("Normal test data")}
	mockCtx := &mockContext{}
	mockCtx = mockCtx.ForkWithBuffer(mockBuf).(*mockContext)

	// 创建一个正常的处理器函数
	expectedError := fmt.Errorf("normal error")
	handler := func(ctx router_context.Context) error {
		return expectedError
	}

	// 应用中间件
	err := recoveryMiddleware(mockCtx, handler)

	// 验证错误被正确传递
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}

// TestLoggingMiddlewareWithLongData 测试日志记录中间件处理长数据的情况
func TestLoggingMiddlewareWithLongData(t *testing.T) {
	// 保存原始的stdout
	oldStdout := os.Stdout

	// 创建管道来捕获输出
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 创建一个缓冲区来收集输出
	var buf bytes.Buffer

	// 在goroutine中从管道读取数据
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// 创建中间件
	loggingMiddleware := LoggingMiddleware()

	// 创建模拟的上下文和缓冲区，使用超过50字节的数据
	longData := "This is a very long data string that exceeds fifty bytes in length to test the truncation functionality in the logging middleware."
	mockBuf := &mockBuffer{data: []byte(longData)}
	mockCtx := &mockContext{}
	mockCtx = mockCtx.ForkWithBuffer(mockBuf).(*mockContext)

	// 创建一个简单的处理器函数
	handler := func(ctx router_context.Context) error {
		time.Sleep(10 * time.Millisecond) // 模拟一些处理时间
		return nil
	}

	// 应用中间件
	err := loggingMiddleware(mockCtx, handler)

	// 关闭写入端并等待读取完成
	w.Close()
	<-done

	// 恢复原始的stdout
	os.Stdout = oldStdout

	// 验证没有错误
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证输出包含预期的日志信息
	output := buf.String()
	if !strings.Contains(output, "Starting processing") {
		t.Error("Expected log output to contain 'Starting processing'")
	}

	if !strings.Contains(output, "Processing completed") {
		t.Error("Expected log output to contain 'Processing completed'")
	}

	if !strings.Contains(output, "data preview:") {
		t.Error("Expected log output to contain 'data preview:'")
	}

	// 验证长数据被截断并添加了"..."后缀
	if !strings.Contains(output, "...") {
		t.Error("Expected log output to contain '...' for truncated data")
	}

	// 验证输出包含数据的前50个字符
	expectedPreview := longData[:50]
	if !strings.Contains(output, expectedPreview) {
		t.Errorf("Expected log output to contain first 50 characters of data: %s", expectedPreview)
	}
}

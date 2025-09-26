package manage

import (
	"testing"

	"github.com/aomirun/content-router/buffer"
)

// mockBuffer 是一个模拟的Buffer实现，用于测试
type mockBuffer struct {
	data        []byte
	resetCalled bool
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

func (m *mockBuffer) Write(p []byte) (n int, err error) {
	m.data = append(m.data, p...)
	return len(p), nil
}

func (m *mockBuffer) WriteString(s string) (n int, err error) {
	m.data = append(m.data, s...)
	return len(s), nil
}

func (m *mockBuffer) Reset() {
	m.resetCalled = true
	m.data = m.data[:0]
}

func (m *mockBuffer) Truncate(n int) {
	if n < len(m.data) {
		m.data = m.data[:n]
	}
}

func (m *mockBuffer) Slice(start, end int) buffer.Buffer {
	return &mockBuffer{
		data: m.data[start:end],
	}
}

func (m *mockBuffer) Clone() buffer.Buffer {
	clone := make([]byte, len(m.data))
	copy(clone, m.data)
	return &mockBuffer{
		data: clone,
	}
}

// mockObjectPool 是一个模拟的ObjectPool实现，用于测试
type mockObjectPool struct {
	acquireCalled bool
	releaseCalled bool
	resetCalled   bool
	buffer        buffer.Buffer
}

func (m *mockObjectPool) Acquire() buffer.Buffer {
	m.acquireCalled = true
	if m.buffer == nil {
		m.buffer = &mockBuffer{data: make([]byte, 0, 1024)}
	}
	return m.buffer
}

func (m *mockObjectPool) Release(buf buffer.Buffer) {
	m.releaseCalled = true
	// 检查buffer是否被重置
	if mockBuf, ok := buf.(*mockBuffer); ok {
		m.resetCalled = mockBuf.resetCalled
	}
}

func (m *mockObjectPool) Size() int {
	return 0
}

func TestNewBufferManager(t *testing.T) {
	manager := NewBufferManager()

	if manager == nil {
		t.Fatal("NewBufferManager should not return nil")
	}

	// 检查返回的manager是否实现了BufferManager接口
	if _, ok := manager.(BufferManager); !ok {
		t.Error("NewBufferManager should return a BufferManager implementation")
	}
}

func TestBufferManager_Acquire(t *testing.T) {
	manager := NewBufferManager()

	// 获取一个缓冲区
	buf := manager.Acquire()

	if buf == nil {
		t.Fatal("Acquire should not return nil")
	}

	// 检查返回的对象是否实现了Buffer接口
	if _, ok := buf.(buffer.Buffer); !ok {
		t.Error("Acquire should return a buffer.Buffer implementation")
	}

	// 测试获取的缓冲区是否可用
	testData := []byte("test data")
	n, err := buf.Write(testData)

	if err != nil {
		t.Errorf("Buffer should be writable: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	if buf.Len() != len(testData) {
		t.Errorf("Expected buffer length %d, got %d", len(testData), buf.Len())
	}
}

func TestBufferManager_Release(t *testing.T) {
	manager := NewBufferManager()

	// 获取一个缓冲区
	buf := manager.Acquire()

	// 向缓冲区写入一些数据
	testData := []byte("test data")
	buf.Write(testData)

	if buf.Len() == 0 {
		t.Error("Buffer should contain data after writing")
	}

	// 释放缓冲区
	manager.Release(buf)

	// 验证缓冲区已被重置（通过再次获取来检查）
	buf2 := manager.Acquire()

	if buf2.Len() != 0 {
		t.Error("Buffer should be reset after release")
	}

	// 检查是否获取到了相同的缓冲区实例（池化效果）
	// 注意：由于sync.Pool的特性，这在并发环境下可能不总是成立
	// 但我们可以通过功能测试来验证池化是否正常工作
}

func TestBufferManager_AcquireReleaseCycle(t *testing.T) {
	manager := NewBufferManager()

	// 执行多次获取和释放操作
	for i := 0; i < 10; i++ {
		buf := manager.Acquire()

		if buf == nil {
			t.Fatalf("Acquire should not return nil on iteration %d", i)
		}

		// 写入数据
		testData := []byte("test data")
		n, err := buf.Write(testData)

		if err != nil {
			t.Errorf("Buffer should be writable on iteration %d: %v", i, err)
		}

		if n != len(testData) {
			t.Errorf("Expected to write %d bytes on iteration %d, wrote %d", len(testData), i, n)
		}

		// 释放缓冲区
		manager.Release(buf)
	}
}

func TestBufferManager_ConcurrentAccess(t *testing.T) {
	manager := NewBufferManager()

	// 并发测试：多个goroutine同时获取和释放缓冲区
	// 这主要测试线程安全性
	done := make(chan bool)

	// 启动多个goroutine
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				buf := manager.Acquire()

				if buf == nil {
					t.Errorf("Acquire should not return nil in goroutine %d, iteration %d", id, j)
					done <- false
					return
				}

				// 写入一些数据
				data := []byte("test")
				buf.Write(data)

				// 释放缓冲区
				manager.Release(buf)
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 5; i++ {
		<-done
	}
}

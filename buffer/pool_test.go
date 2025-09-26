package buffer

import (
	"sync"
	"testing"
)

func TestObjectPoolAcquire(t *testing.T) {
	// 创建一个新的对象池
	pool := NewPool()
	
	// 从池中获取一个对象
	buf := pool.Acquire()
	
	// 验证返回的对象不为nil
	if buf == nil {
		t.Fatal("Acquire should not return nil")
	}
	
	// 验证返回的对象实现了Buffer接口
	if _, ok := buf.(Buffer); !ok {
		t.Error("Acquired object should implement Buffer interface")
	}
	
	// 验证获取的对象是空的
	if buf.Len() != 0 {
		t.Errorf("Acquired buffer should be empty, got length %d", buf.Len())
	}
	
	// 测试写入数据
	testData := []byte("test data")
	n, err := buf.Write(testData)
	if err != nil {
		t.Errorf("Buffer should be writable: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}
	
	// 验证数据已写入
	if buf.Len() != len(testData) {
		t.Errorf("Buffer length should be %d, got %d", len(testData), buf.Len())
	}
}

func TestObjectPoolRelease(t *testing.T) {
	// 创建一个新的对象池
	pool := NewPool()
	
	// 获取一个对象并写入数据
	buf := pool.Acquire()
	testData := []byte("test data")
	buf.Write(testData)
	
	// 验证缓冲区包含数据
	if buf.Len() == 0 {
		t.Error("Buffer should contain data after writing")
	}
	
	// 释放对象回池中
	pool.Release(buf)
	
	// 再次获取对象
	buf2 := pool.Acquire()
	
	// 验证缓冲区已被重置
	if buf2.Len() != 0 {
		t.Errorf("Buffer should be reset after release, got length %d", buf2.Len())
	}
}

func TestObjectPoolAcquireReleaseCycle(t *testing.T) {
	// 创建一个新的对象池
	pool := NewPool()
	
	// 执行多次获取和释放操作
	for i := 0; i < 10; i++ {
		buf := pool.Acquire()
		
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
		
		// 释放对象
		pool.Release(buf)
	}
}

func TestObjectPoolSize(t *testing.T) {
	// 创建一个新的对象池
	pool := NewPool()
	
	// 检查初始大小
	size := pool.Size()
	// 注意：由于sync.Pool的实现，Size()方法可能总是返回0
	// 这里我们只是验证方法可以被调用而不panic
	
	// 获取一个对象
	buf := pool.Acquire()
	
	// 释放对象
	pool.Release(buf)
	
	// 再次检查大小
	newSize := pool.Size()
	
	// 我们不验证具体的大小值，因为sync.Pool的行为可能因Go版本而异
	// 只要方法能正常调用即可
	_ = size
	_ = newSize
}

func TestObjectPoolMultipleAcquire(t *testing.T) {
	// 创建一个新的对象池
	pool := NewPool()
	
	// 连续获取多个对象
	buf1 := pool.Acquire()
	buf2 := pool.Acquire()
	buf3 := pool.Acquire()
	
	// 验证所有对象都不为nil
	if buf1 == nil || buf2 == nil || buf3 == nil {
		t.Fatal("All acquired buffers should not be nil")
	}
	
	// 验证对象是独立的
	data1 := []byte("data1")
	data2 := []byte("data2")
	data3 := []byte("data3")
	
	buf1.Write(data1)
	buf2.Write(data2)
	buf3.Write(data3)
	
	if buf1.Len() != len(data1) {
		t.Errorf("buf1 length should be %d, got %d", len(data1), buf1.Len())
	}
	
	if buf2.Len() != len(data2) {
		t.Errorf("buf2 length should be %d, got %d", len(data2), buf2.Len())
	}
	
	if buf3.Len() != len(data3) {
		t.Errorf("buf3 length should be %d, got %d", len(data3), buf3.Len())
	}
	
	// 释放所有对象
	pool.Release(buf1)
	pool.Release(buf2)
	pool.Release(buf3)
	
	// 再次获取对象并验证它们已被重置
	buf4 := pool.Acquire()
	buf5 := pool.Acquire()
	buf6 := pool.Acquire()
	
	if buf4.Len() != 0 {
		t.Errorf("buf4 should be reset, got length %d", buf4.Len())
	}
	
	if buf5.Len() != 0 {
		t.Errorf("buf5 should be reset, got length %d", buf5.Len())
	}
	
	if buf6.Len() != 0 {
		t.Errorf("buf6 should be reset, got length %d", buf6.Len())
	}
}

func TestObjectPoolReleaseNil(t *testing.T) {
	// 创建一个新的对象池
	pool := NewPool()
	
	// 测试释放nil对象不会panic
	pool.Release(nil)
}

func TestObjectPoolAcquireZeroValue(t *testing.T) {
	// 创建一个自定义的poolImpl，其sync.Pool没有New函数
	// 这样当池为空时，Acquire会返回零值
	pool := &poolImpl[Buffer]{
		pool: sync.Pool{},
	}
	
	// 从空池中获取对象，应该返回零值
	buf := pool.Acquire()
	
	// 验证返回的是零值（对于interface{}来说是nil）
	if buf != nil {
		t.Error("Acquire from empty pool should return nil/zero value")
	}
}

func TestObjectPoolGeneric(t *testing.T) {
	// 测试创建泛型对象池
	pool := &poolImpl[Buffer]{
		pool: sync.Pool{
			New: func() interface{} {
				return NewBuffer()
			},
		},
	}
	
	// 测试获取对象
	buf := pool.Acquire()
	if buf == nil {
		t.Error("Acquire should not return nil")
	}
	
	// 测试释放对象
	pool.Release(buf)
	
	// 测试大小方法
	size := pool.Size()
	_ = size // 只是确保方法可以调用
}
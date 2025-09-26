package buffer

import (
	"testing"
)

func TestBufferImpl(t *testing.T) {
	// 创建一个新的Buffer实例
	buf := NewBuffer()
	
	// 测试写入数据
	data := []byte("Hello, World!")
	n, err := buf.Write(data)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Write returned %d, expected %d", n, len(data))
	}
	
	// 测试读取数据
	readData := buf.Get()
	if len(readData) != len(data) {
		t.Errorf("Get returned %d bytes, expected %d", len(readData), len(data))
	}
	
	// 测试长度
	if buf.Len() != len(data) {
		t.Errorf("Len returned %d, expected %d", buf.Len(), len(data))
	}
	
	// 测试容量
	if buf.Cap() < len(data) {
		t.Errorf("Cap returned %d, expected at least %d", buf.Cap(), len(data))
	}
	
	// 测试写入字符串
	str := " Hello, Go!"
	n, err = buf.WriteString(str)
	if err != nil {
		t.Errorf("WriteString failed: %v", err)
	}
	if n != len(str) {
		t.Errorf("WriteString returned %d, expected %d", n, len(str))
	}
	
	// 测试重置
	buf.Reset()
	if buf.Len() != 0 {
		t.Errorf("Reset failed, Len returned %d, expected 0", buf.Len())
	}
	
	// 测试截断
	buf.Write([]byte("Hello, World!"))
	buf.Truncate(5)
	if buf.Len() != 5 {
		t.Errorf("Truncate failed, Len returned %d, expected 5", buf.Len())
	}
	
	// 测试克隆
	original := NewBuffer()
	original.Write([]byte("Clone test"))
	cloned := original.Clone()
	if cloned.Len() != original.Len() {
		t.Errorf("Clone failed, lengths don't match: %d vs %d", cloned.Len(), original.Len())
	}
	
	// 修改原始缓冲区，确保克隆是独立的
	original.Write([]byte(" modified"))
	if cloned.Len() == original.Len() {
		t.Errorf("Clone is not independent, lengths match: %d", cloned.Len())
	}
}

func TestBufferSlice(t *testing.T) {
	// 创建一个新的Buffer实例并写入数据
	buf := NewBuffer()
	data := []byte("Hello, World! This is a test for slice.")
	buf.Write(data)
	
	// 测试切片功能
	start, end := 7, 12
	sliced := buf.Slice(start, end)
	
	// 验证切片内容
	expected := data[start:end]
	actual := sliced.Get()
	
	if len(actual) != len(expected) {
		t.Errorf("Slice length mismatch: got %d, expected %d", len(actual), len(expected))
	}
	
	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("Slice content mismatch at index %d: got %c, expected %c", i, actual[i], expected[i])
		}
	}
	
	// 验证切片的独立性
	// 修改原缓冲区不应该影响切片
	buf.Write([]byte(" additional"))
	if len(sliced.Get()) != len(expected) {
		t.Error("Slice is not independent from original buffer")
	}
	
	// 测试边界情况
	// 1. 切片到缓冲区末尾
	endSlice := buf.Slice(7, buf.Len())
	if endSlice.Len() != buf.Len()-7 {
		t.Errorf("End slice length mismatch: got %d, expected %d", endSlice.Len(), buf.Len()-7)
	}
	
	// 2. 切片整个缓冲区
	fullSlice := buf.Slice(0, buf.Len())
	if fullSlice.Len() != buf.Len() {
		t.Errorf("Full slice length mismatch: got %d, expected %d", fullSlice.Len(), buf.Len())
	}
}

func TestBufferTruncateEdgeCases(t *testing.T) {
	buf := NewBuffer()
	
	// 测试截断空缓冲区
	buf.Truncate(0)
	if buf.Len() != 0 {
		t.Error("Truncate on empty buffer should not change length")
	}
	
	// 测试截断到更大索引
	buf.Write([]byte("test"))
	buf.Truncate(10) // 大于当前长度
	if buf.Len() != 4 { // 应该保持不变
		t.Errorf("Truncate to larger index should not change length, got %d, expected 4", buf.Len())
	}
	
	// 测试截断到0
	buf.Truncate(0)
	if buf.Len() != 0 {
		t.Errorf("Truncate to 0 should result in empty buffer, got length %d", buf.Len())
	}
}

func TestBufferWriteEdgeCases(t *testing.T) {
	buf := NewBuffer()
	
	// 测试写入空数据
	n, err := buf.Write([]byte{})
	if err != nil {
		t.Errorf("Write empty slice should not fail: %v", err)
	}
	if n != 0 {
		t.Errorf("Write empty slice should return 0 bytes written, got %d", n)
	}
	
	// 测试写入字符串
	n, err = buf.WriteString("")
	if err != nil {
		t.Errorf("WriteString empty string should not fail: %v", err)
	}
	if n != 0 {
		t.Errorf("WriteString empty string should return 0 bytes written, got %d", n)
	}
}
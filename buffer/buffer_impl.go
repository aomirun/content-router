package buffer

// bufferImpl 是Buffer接口的具体实现
type bufferImpl struct {
	data []byte
}

// Get 获取底层字节数组的引用
func (b *bufferImpl) Get() []byte {
	return b.data
}

// Len 获取当前有效数据长度
func (b *bufferImpl) Len() int {
	return len(b.data)
}

// Cap 获取缓冲区容量
func (b *bufferImpl) Cap() int {
	return cap(b.data)
}

// Write 写入数据到缓冲区，必要时会扩容
// 与标准库io.Writer接口兼容
func (b *bufferImpl) Write(p []byte) (n int, err error) {
	// 实现写入逻辑
	b.data = append(b.data, p...)
	return len(p), nil
}

// WriteString 写入字符串到缓冲区，必要时会扩容
// 与标准库io.StringWriter接口兼容
func (b *bufferImpl) WriteString(s string) (n int, err error) {
	// 实现写入字符串逻辑
	b.data = append(b.data, s...)
	return len(s), nil
}

// Reset 重置缓冲区，保留底层数组但清空内容
func (b *bufferImpl) Reset() {
	// 实现重置逻辑
	b.data = b.data[:0]
}

// Truncate 将缓冲区截断到指定长度
func (b *bufferImpl) Truncate(n int) {
	// 实现截断逻辑
	if n < len(b.data) {
		b.data = b.data[:n]
	}
}

// Slice 创建子切片但不复制数据
func (b *bufferImpl) Slice(start, end int) Buffer {
	// 实现切片逻辑
	return &bufferImpl{
		data: b.data[start:end],
	}
}

// Clone 创建缓冲区的深拷贝
func (b *bufferImpl) Clone() Buffer {
	// 实现克隆逻辑
	clone := make([]byte, len(b.data))
	copy(clone, b.data)
	return &bufferImpl{
		data: clone,
	}
}

// NewBuffer 创建一个新的Buffer实例
func NewBuffer() Buffer {
	return &bufferImpl{
		data: make([]byte, 0, 1024), // 初始容量1024字节
	}
}

package buffer

// Readable 定义可读缓冲区接口
type Readable interface {
	// Get 获取底层字节数组的引用
	Get() []byte

	// Len 获取当前有效数据长度
	Len() int

	// Cap 获取缓冲区容量
	Cap() int
}

// Writable 定义可写缓冲区接口
type Writable interface {
	// Write 写入数据到缓冲区，必要时会扩容
	// 与标准库io.Writer接口兼容
	Write(p []byte) (n int, err error)

	// WriteString 写入字符串到缓冲区，必要时会扩容
	// 与标准库io.StringWriter接口兼容
	WriteString(s string) (n int, err error)
}

// Mutable 定义可变缓冲区接口
type Mutable interface {
	// Reset 重置缓冲区，保留底层数组但清空内容
	Reset()

	// Truncate 将缓冲区截断到指定长度
	Truncate(n int)
}

// Sliceable 定义可切片缓冲区接口
type Sliceable interface {
	// Slice 创建子切片但不复制数据
	Slice(start, end int) Buffer
}

// Cloneable 定义可克隆缓冲区接口
type Cloneable interface {
	// Clone 创建缓冲区的深拷贝
	Clone() Buffer
}

// Buffer 定义可重用的缓冲区接口
// 它组合了所有缓冲区操作接口
type Buffer interface {
	Readable
	Writable
	Mutable
	Sliceable
	Cloneable
}
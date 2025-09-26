package manage

import (
	"github.com/aomirun/content-router/buffer"
)

// bufferManagerImpl 是BufferManager接口的实现
type bufferManagerImpl struct {
	pool buffer.ObjectPool[buffer.Buffer]
}

// NewBufferManager 创建一个新的BufferManager实例
func NewBufferManager() BufferManager {
	return &bufferManagerImpl{
		pool: buffer.NewPool(),
	}
}

// Acquire 从池中获取一个缓冲区
func (bm *bufferManagerImpl) Acquire() buffer.Buffer {
	return bm.pool.Acquire()
}

// Release 将缓冲区释放回池中
func (bm *bufferManagerImpl) Release(buf buffer.Buffer) {
	// 重置缓冲区后再放回池中
	buf.Reset()
	bm.pool.Release(buf)
}

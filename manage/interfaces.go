package manage

import (
	"github.com/aomirun/content-router/buffer"
)

// BufferManager 定义缓冲区管理器接口
// 该接口提供了缓冲区的获取和释放功能，用于高效管理缓冲区资源
//
// 主要功能:
// 1. 缓冲区池化管理
// 2. 减少内存分配和垃圾回收压力
// 3. 提供统一的缓冲区获取和释放接口
type BufferManager interface {
	// Acquire 从池中获取一个缓冲区
	// 返回: 可用的缓冲区实例
	Acquire() buffer.Buffer

	// Release 将缓冲区释放回池中
	// buf: 需要释放的缓冲区实例
	Release(buf buffer.Buffer)
}

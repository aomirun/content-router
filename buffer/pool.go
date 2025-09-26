package buffer

import "sync"

// ObjectPool 定义通用对象池接口
// 所有对象池实现应该遵循此接口，提供一致的获取和释放方法
// T 是池化对象的类型参数
//
// 实现此接口的类型应该确保:
// 1. 线程安全性
// 2. 对象重置的完整性
// 3. 合理的预分配策略
// 4. 避免内存泄漏
//
// 命名规范:
// - 获取对象: AcquireXXX()
// - 释放对象: ReleaseXXX(obj XXX)
// - 对象池实例: xxxPool
// - 池实现: xxxPoolImpl
type ObjectPool[T any] interface {
	// Acquire 从池中获取一个对象实例
	// 如果池为空，会创建一个新的对象实例
	// 返回的对象应该已经重置为初始状态
	Acquire() T

	// Release 将对象实例归还池中
	// 归还前应确保对象已适当重置
	// nil对象应当被安全处理
	Release(obj T)

	// Size 返回池中当前可用对象数量的估计值
	// 注意：此方法主要用于监控和调试，不保证精确性
	// 在高并发环境下，返回值可能不准确
	Size() int
}

// poolImpl 是ObjectPool接口的具体实现
type poolImpl[T any] struct {
	pool sync.Pool
}

// Acquire 从池中获取一个对象实例
func (p *poolImpl[T]) Acquire() T {
	obj := p.pool.Get()
	if obj == nil {
		var zero T
		return zero
	}
	return obj.(T)
}

// Release 将对象实例归还池中
func (p *poolImpl[T]) Release(obj T) {
	// 如果对象实现了Mutable接口，重置它
	if mutable, ok := interface{}(obj).(Mutable); ok {
		mutable.Reset()
	}
	p.pool.Put(obj)
}

// Size 返回池中当前可用对象数量的估计值
func (p *poolImpl[T]) Size() int {
	// sync.Pool没有提供获取大小的方法，这里返回0
	return 0
}

// NewPool 创建一个新的对象池
func NewPool() ObjectPool[Buffer] {
	return &poolImpl[Buffer]{
		pool: sync.Pool{
			New: func() interface{} {
				return NewBuffer()
			},
		},
	}
}

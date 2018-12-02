package zero

import "sync"

type BufferPool struct {
	pool       *sync.Pool
	bufferSize int
}

var DefaultBufferPool *BufferPool = nil

func init() {
	DefaultBufferPool = NewBufferPool(20)
}

func NewBufferPool(bufferSize int) *BufferPool {
	pool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, bufferSize)
		},
	}
	return &BufferPool{
		bufferSize: bufferSize,
		pool:       pool,
	}
}

func (b *BufferPool) Get() []byte {
	if val, ok := b.pool.Get().([]byte); ok {
		return val
	}
	return make([]byte, b.bufferSize)
}

func (b *BufferPool) Put(buffer []byte) {
	if len(buffer) != b.bufferSize {
		return
	}
	b.pool.Put(buffer)
}

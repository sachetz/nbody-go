package barrier

import (
	"sync"
)

type Barrier struct {
	counter int
	size    int
	mutex   *sync.Mutex
	cond    *sync.Cond
}

func NewBarrier(size int) *Barrier {
	b := &Barrier{}
	b.counter = 0
	b.size = size
	b.mutex = &sync.Mutex{}
	b.cond = sync.NewCond(b.mutex)
	return b
}

func (b *Barrier) Wait() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.counter += 1
	if b.counter < b.size {
		b.cond.Wait()
	} else {
		b.cond.Broadcast()
	}
}

func (b *Barrier) Reset() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.counter = 0
}

package queue

import (
	"sync/atomic"
	"unsafe"
)

type BoundedDeque struct {
	start  int
	end    int
	bottom int32
	top    unsafe.Pointer
}

func NewBoundedDeque(begin int, end int) *BoundedDeque {
	topPtr := AtomicStampedReference{id: begin, stamp: 0}
	q := BoundedDeque{top: unsafe.Pointer(&topPtr), bottom: int32(end), start: begin, end: end}
	return &q
}

func (q *BoundedDeque) PopTop() int {
	oldTop := loadASRFromPointer(&q.top)
	newTop := AtomicStampedReference{id: oldTop.id + 1, stamp: oldTop.stamp + 1}
	if int(q.bottom) <= oldTop.id {
		return -1
	}
	if compareAndSwap(&q.top, oldTop, &newTop) {
		return oldTop.id
	}
	return -1
}

func (q *BoundedDeque) PopBottom() int {
	bottom := int(q.bottom)
	if bottom == q.start {
		return -1
	}
	atomic.AddInt32(&q.bottom, -1)
	bottom -= 1
	particleIdx := bottom
	oldTop := loadASRFromPointer(&q.top)
	newTop := AtomicStampedReference{id: q.end - 1, stamp: oldTop.stamp + 1}
	if bottom > oldTop.id {
		return particleIdx
	}
	if bottom == oldTop.id {
		atomic.StoreInt32(&q.bottom, int32(q.start))
		if compareAndSwap(&q.top, oldTop, &newTop) {
			return particleIdx
		}
	}
	atomic.StorePointer(&q.top, unsafe.Pointer(&newTop))
	return -1
}

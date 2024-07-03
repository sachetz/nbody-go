package queue

import (
	"sync/atomic"
	"unsafe"
)

type AtomicStampedReference struct {
	id    int
	stamp int
}

func loadASRFromPointer(p *unsafe.Pointer) *AtomicStampedReference {
	if p == nil {
		return nil
	}
	return (*AtomicStampedReference)(atomic.LoadPointer(p))
}

func compareAndSwap(ptr *unsafe.Pointer, old *AtomicStampedReference, new *AtomicStampedReference) bool {
	var oldPtr unsafe.Pointer
	oldPtr = nil
	if old != nil {
		oldPtr = unsafe.Pointer(old)
	}
	return atomic.CompareAndSwapPointer(ptr, oldPtr, unsafe.Pointer(new))
}
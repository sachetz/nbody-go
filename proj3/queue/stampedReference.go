package queue

import "sync/atomic"

type AtomicStampedReference struct {
	val atomic.Int64
}

func combine(ref int32, stamp int32) int64 {
	return (int64(ref) << 32) | (int64(stamp) & 0xFFFFFFFF)
}

func split(val int64) (int32, int32) {
	return int32(val >> 32), int32(val & 0xFFFFFFFF)
}

func NewAtomicStampedReference(initialRef int32, initialStamp int32) *AtomicStampedReference {
	stampedRef := &AtomicStampedReference{}
	stampedRef.val.Store(combine(initialRef, initialStamp))
	return stampedRef
}

func (asr *AtomicStampedReference) Get(stampHolder *int32) int32 {
	ref, stamp := split(asr.val.Load())
	*stampHolder = stamp
	return ref
}

func (asr *AtomicStampedReference) GetReference() int32 {
	ref, _ := split(asr.val.Load())
	return ref
}

func (asr *AtomicStampedReference) Set(newRef int32, newStamp int32) {
	asr.val.Store(combine(newRef, newStamp))
}

func (asr *AtomicStampedReference) CompareAndSet(expectedRef, newRef, expectedStamp, newStamp int32) bool {
	oldVal := combine(expectedRef, expectedStamp)
	newVal := combine(newRef, newStamp)
	return asr.val.CompareAndSwap(oldVal, newVal)
}

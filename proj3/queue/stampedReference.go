package queue

import (
	"sync"
)

type AtomicStampedReference struct {
	reference int32
	stamp     int32
	lock      sync.Mutex
}

func NewAtomicStampedReference(initialRef int32, initialStamp int32) *AtomicStampedReference {
	stampedRef := &AtomicStampedReference{}
	stampedRef.reference = initialRef
	stampedRef.stamp = initialStamp
	return stampedRef
}

func (asr *AtomicStampedReference) Get(stampHolder *int32) int32 {
	asr.lock.Lock()
	defer asr.lock.Unlock()
	*stampHolder = asr.stamp
	return asr.reference
}

func (asr *AtomicStampedReference) GetReference() int32 {
	asr.lock.Lock()
	defer asr.lock.Unlock()
	return asr.reference
}

func (asr *AtomicStampedReference) Set(newRef int32, newStamp int32) {
	asr.lock.Lock()
	defer asr.lock.Unlock()
	asr.reference = newRef
	asr.stamp = newStamp
}

func (asr *AtomicStampedReference) CompareAndSet(expectedRef, newRef, expectedStamp, newStamp int32) bool {
	asr.lock.Lock()
	defer asr.lock.Unlock()
	if asr.reference == expectedRef && asr.stamp == expectedStamp {
		asr.stamp = newStamp
		asr.reference = newRef
		return true
	}
	return false
}

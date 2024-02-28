package queue

import "sync/atomic"

type BoundedDeque struct {
	tasks  []func()
	bottom atomic.Int32
	top    *AtomicStampedReference
}

func NewBDEQueue(capacity int) *BoundedDeque {
	queue := &BoundedDeque{}
	queue.tasks = make([]func(), capacity)
	queue.bottom.Store(0)
	queue.top = &AtomicStampedReference{}
	return queue
}

// PushBottom is only used by the owner thread - concurrency not required
func (q *BoundedDeque) PushBottom(task func()) {
	q.tasks[q.bottom.Load()] = task
	q.bottom.Add(1)
}

// Try to atomically pop from the top and increment the top reference and stamp
func (q *BoundedDeque) PopTop() func() {
	var oldStamp int32
	oldTop := q.top.Get(&oldStamp)
	newTop := oldTop + 1
	newStamp := oldStamp + 1
	// Queue is empty, nothing to pop
	if q.bottom.Load() <= oldTop {
		return nil
	}
	// Try to atomically update the top reference and the stamp
	f := q.tasks[oldTop]
	if q.top.CompareAndSet(oldTop, newTop, oldStamp, newStamp) {
		return f
	}
	return nil
}

func (q *BoundedDeque) PopBottom() func() {
	if q.bottom.Load() == 0 { // If queue is empty, return nil
		return nil
	}

	// Otherwise decrement bottom and claim a task
	q.bottom.Add(-1)
	f := q.tasks[q.bottom.Load()]

	var oldStamp int32
	oldTop := q.top.Get(&oldStamp)
	var newTop int32 = 0
	newStamp := oldStamp + 1

	// Return the function - bottom is not equal to top and concurrency not required
	if q.bottom.Load() > oldTop {
		return f
	}
	// Bottom = top, atomically update the top reference and stamp - other threads might be interacting
	if q.bottom.Load() == oldTop {
		q.bottom.Store(0)
		if q.top.CompareAndSet(oldTop, newTop, oldStamp, newStamp) {
			return f
		}
	}
	// Failed to fetch, i.e. thief succeed but top not updated, update and return
	q.top.Set(newTop, newStamp)
	return nil
}

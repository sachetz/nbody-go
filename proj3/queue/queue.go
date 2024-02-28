package queue

type BoundedDeque struct {
	tasks  []func()
	bottom int32
	top    *AtomicStampedReference
}

func NewBDEQueue(capacity int) *BoundedDeque {
	queue := &BoundedDeque{}
	queue.tasks = make([]func(), capacity)
	queue.bottom = 0
	queue.top = &AtomicStampedReference{}
	return queue
}

func (q *BoundedDeque) PushBottom(task func()) {
	q.tasks[q.bottom] = task
	q.bottom += 1
}

func (q *BoundedDeque) IsEmpty() bool {
	return q.top.GetReference() < q.bottom
}

func (q *BoundedDeque) PopTop() func() {
	var oldStamp int32
	oldTop := q.top.Get(&oldStamp)
	newTop := oldTop + 1
	newStamp := oldStamp + 1
	if q.bottom <= oldTop {
		return nil
	}
	f := q.tasks[oldTop]
	if q.top.CompareAndSet(oldTop, newTop, oldStamp, newStamp) {
		return f
	}
	return nil
}

func (q *BoundedDeque) PopBottom() func() {
	if q.bottom == 0 {
		return nil
	}
	q.bottom -= 1
	f := q.tasks[q.bottom]
	var oldStamp int32
	oldTop := q.top.Get(&oldStamp)
	var newTop int32 = 0
	newStamp := oldStamp + 1
	if q.bottom > oldTop {
		return f
	}
	if q.bottom == oldTop {
		q.bottom = 0
		if q.top.CompareAndSet(oldTop, newTop, oldStamp, newStamp) {
			return f
		}
	}
	q.top.Set(newTop, newStamp)
	return nil
}

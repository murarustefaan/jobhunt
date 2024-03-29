package queue

import "fmt"

//type queueNode[W interface{}] struct {
//	value W
//	next  *queueNode[W]
//}

type MemoryQueue[D interface{}] struct {
	data    []*D
	queue   chan *D
	dequeue chan *D
}

func NewMemoryQueue[D interface{}]() *MemoryQueue[D] {
	mq := &MemoryQueue[D]{
		data:    make([]*D, 0, 16),
		queue:   make(chan *D),
		dequeue: make(chan *D),
	}

	go mq.start()

	return mq
}

func (q *MemoryQueue[D]) start() {
	for {
		select {
		case item := <-q.queue:
			q.data = append(q.data, item)

		case q.dequeue <- q.data[0]:
			q.data = q.data[1:]
		}
	}
}

func (q *MemoryQueue[D]) Push(item *D) error {
	q.queue <- item
	return nil
}

func (q *MemoryQueue[D]) Pop() (*D, error) {
	select {
	case item := <-q.dequeue:
		return item, nil
	default:
		return nil, fmt.Errorf("pop from empty queue")
	}
}

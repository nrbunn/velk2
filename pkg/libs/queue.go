package libs

import "errors"

type Queue struct {
	capacity int
	q        chan string
}

func (q *Queue) Insert(item string) error {
	if len(q.q) < int(q.capacity) {
		q.q <- item
		return nil
	}
	return errors.New("Queue is full")
}

func (q *Queue) Remove() (string, error) {
	if len(q.q) > 0 {
		item := <-q.q
		return item, nil
	}
	return "", errors.New("Queue is empty")
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		capacity: capacity,
		q:        make(chan string, capacity),
	}
}

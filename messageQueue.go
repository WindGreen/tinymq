package tinymq

import (
	"github.com/pkg/errors"
	"time"
)

var (
	ErrBlock    = errors.New("Message Queue is block")
	ErrNotFound = errors.New("Topic not found")
	ErrEmpty    = errors.New("Queue is empty")
	ErrTimeout  = errors.New("Operation timeout")
)

type MQ struct {
	Size     int
	RTimeout int
	WTimeout int
	topics   map[string]*Queue
}

type Queue struct {
	list chan Item
}

type Item struct {
	Payload interface{}
}

func NewMQ(size, timeout int) *MQ {
	mq := MQ{
		Size:     size,
		RTimeout: timeout,
		WTimeout: timeout,
	}
	return &mq
}

func NewQueue(size int) *Queue {
	q := Queue{}
	q.list = make(chan Item, size)
	return &q
}

func (m *MQ) Run() {
	m.topics["default"] = NewQueue(m.Size)
}

func (m *MQ) Remove(name string) {

}

func (m *MQ) Publish(name string) {

}

func (m *MQ) Subscribe(name string) {

}

func (m *MQ) Unsubscribe(name string) {

}

func (m *MQ) ProduceWithTimeout(topic string, content interface{}) error {
	if t, ok := m.topics[topic]; ok {
		select {
		case t.list <- Item{
			Payload: content,
		}:
			return nil
		case time.After(time.Duration(m.WTimeout) * time.Millisecond):
			return ErrTimeout
		}
	}
	return ErrNotFound
}

func (m *MQ) Produce(topic string, content interface{}) error {
	if t, ok := m.topics[topic]; ok {
		t.list <- Item{
			Payload: content,
		}
		return nil
	}
	return ErrNotFound
}

func (m *MQ) ConsumeWithTimeout(topic string, wait bool) (error, interface{}) {
	if t, ok := m.topics[topic]; ok {
		if !wait && len(t.list) == 0 {
			return ErrEmpty, nil
		}
		select {
		case item := <-t.list:
			return nil, item
		case time.After(time.Duration(m.RTimeout) * time.Millisecond):
			return ErrTimeout, nil
		}
	}
	return ErrNotFound, nil
}

func (m *MQ) Consume(topic string, wait bool) (error, interface{}) {
	if t, ok := m.topics[topic]; ok {
		return nil, <-t.list
	}
	return ErrNotFound, nil
}

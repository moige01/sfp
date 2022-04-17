package pubsub

import (
	"sync"

	"sugud0r.dev/sfp/internal/loggin"
)

type Pubsub struct {
	mu     *sync.Mutex
	subs   map[int][]chan string
	closed bool
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		mu:   &sync.Mutex{},
		subs: make(map[int][]chan string),
	}
}

func (ps *Pubsub) Subscribe(topic int) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)

	return ch
}

func (ps *Pubsub) Publish(topic int, msg string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		loggin.Warning.Print("Trying to publishing on closed pubsub.")

		return
	}

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		loggin.Warning.Print("Trying to close an already closed pubsub.")

		return
	}

	ps.closed = true

	for _, subs := range ps.subs {
		for _, ch := range subs {
			close(ch)
		}
	}
}

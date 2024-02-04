package subscriptions

import (
	"context"
	"sync"
)

type INotifier interface {
	Subscribe(ctx context.Context, observer Listener)
	Send(onListen bool) error
	Close()
}

type Notifier struct {
	mtx       sync.Mutex
	observers []Listener
	onListen  chan bool
}

func NewNotifier() *Notifier {
	return &Notifier{
		mtx:       sync.Mutex{},
		observers: make([]Listener, 0),
		onListen:  make(chan bool, 1),
	}
}

func (r *Notifier) Subscribe(ctx context.Context, observer Listener) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.observers = append(r.observers, observer)

	go func(ctx context.Context, signal chan bool, observer Listener) {
		for running := range signal {
			if !running {
				continue
			}
			break
		}
		observer.OnNotify(ctx)
	}(ctx, r.onListen, observer)
}

func (r *Notifier) Close() {
	close(r.onListen)
}

func (r *Notifier) Send(onListen bool) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	for i := 0; i < len(r.observers); i++ {
		go func(onListen bool) {
			r.onListen <- onListen
		}(onListen)
	}
	return nil
}

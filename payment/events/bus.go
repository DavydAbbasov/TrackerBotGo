package events

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type Bus interface {
	Publish(ctx context.Context, e Event) error
	Close()
	Chan() <-chan Event
}

// top-levl errors
var (
	ErrClosed      = errors.New("events bus is closed")
	ErrPublishBusy = errors.New("publish time out (queue busy)")
)

// поверх buffer канала.
type bus struct {
	ch     chan Event
	once   sync.Once   // Close однократное закрытие
	closed atomic.Bool 
}

// поверх chan Event с  buffer
func NewBus(buffer int) Bus { //Технично интерфейсное значение = пара: (динамический тип, динамическое значение).
	if buffer <= 0 {
		buffer = 1
	}
	return &bus{
		ch: make(chan Event, buffer),
	}
}

// Производитель (producer) — пишет в канал: ch <- v
// Потребители (consumers) — читают из канала: <-ch / for range ch (у тебя это воркеры).
func (b *bus) Publish(ctx context.Context, e Event) error {
	if b.closed.Load() {
		return ErrClosed
	}
	select {
	case b.ch <- e:
		return nil
	case <-ctx.Done():
		return ErrPublishBusy
	}
}

func (b *bus) Close() {
	b.once.Do(func() {
		b.closed.Store(true)
		close(b.ch)
	},
	)
}

func (b *bus) Chan() <-chan Event {
	return b.ch
}

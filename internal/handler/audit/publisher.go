package audit

import (
	"context"
	"sync"

	"github.com/liebeSonne/shortlink/internal/logger"
)

// Publisher - интерфейс издателя событий.
type Publisher interface {
	// Subscribe - подписать наблюдателя.
	Subscribe(observer Observer)
	// Notify - уведомить о событии.
	Notify(event Event)
}

// NewPublisher - создание экземпляра издателя событий.
func NewPublisher(
	ctx context.Context,
	chanSize, countWorkers uint,
	logger logger.Logger,
) Publisher {
	p := &publisher{
		ch:        make(chan Event, chanSize),
		observers: make([]Observer, 0),
		logger:    logger,
	}

	for i := 0; i < int(countWorkers); i++ {
		go p.runWorker(ctx)
	}

	return p
}

type publisher struct {
	mu        sync.RWMutex
	observers []Observer
	ch        chan Event
	logger    logger.Logger
}

func (p *publisher) Subscribe(observer Observer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observers = append(p.observers, observer)
}

func (p *publisher) Notify(event Event) {
	p.ch <- event
}

func (p *publisher) handleEvent(event Event) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, observer := range p.observers {
		observer.Update(event)
	}
}

func (p *publisher) runWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.logger.Debugf("publisher worker: context closed")
			return
		case event := <-p.ch:
			p.handleEvent(event)
		}
	}
}

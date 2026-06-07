package audit

import "sync"

// Publisher - интерфейс издателя событий.
type Publisher interface {
	// Subscribe - подписать наблюдателя.
	Subscribe(observer Observer)
	// Notify - уведомить о событии.
	Notify(event Event)
}

// NewPublisher - создание экземпляра издателя событий.
func NewPublisher() Publisher {
	return &publisher{}
}

type publisher struct {
	mu        sync.RWMutex
	observers []Observer
}

func (p *publisher) Subscribe(observer Observer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observers = append(p.observers, observer)
}

func (p *publisher) Notify(event Event) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, observer := range p.observers {
		observer.Update(event)
	}
}

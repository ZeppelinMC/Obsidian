package broadcast

import "sync"

func New[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		players: make(map[string]T),
	}
}

type Broadcaster[T any] struct {
	players map[string]T
	mu      sync.RWMutex
}

func (b *Broadcaster[T]) Set(username string, p T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.players[username] = p
}

func (b *Broadcaster[T]) Count() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.players)
}

func (b *Broadcaster[T]) Get(username string) T {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.players[username]
}

func (b *Broadcaster[T]) Remove(username string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.players, username)
}

func (b *Broadcaster[T]) Range(f func(T) bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, p := range b.players {
		if !f(p) {
			break
		}
	}
}

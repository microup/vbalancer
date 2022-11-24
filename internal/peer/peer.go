package peer

import "sync"

type Peer struct {
	Name  string `yaml:"Name"`
	Proto string `yaml:"Proto"`
	Uri   string `yaml:"URI"`
	Alive bool
	mux   sync.RWMutex
}

func (b *Peer) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive

}

func (b *Peer) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return true
}

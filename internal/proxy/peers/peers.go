package peers

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/types"

	cache "github.com/microup/vcache"
)

// Peers define a struct that contains a list of peers.
type Peers struct {
	currentPeerIndex            *uint64
	blackListNotResponsePeers   *cache.VCache
	TimeToEvictNotResponsePeers time.Duration `yaml:"timeToEvictNotResponsePeers"`
	List                        []peer.Peer   `yaml:"list" json:"list"`
}

// Initialize Peers struct with a slice of Peer objects,
// copy peers from input to the new slice, set CurrentPeerIndex to 0
// to track the selected peer's index in the slice.
func (p *Peers) Init(ctx context.Context, peers []peer.Peer) error {
	p.blackListNotResponsePeers = cache.New(time.Second, p.TimeToEvictNotResponsePeers)

	err := p.blackListNotResponsePeers.StartEvict(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var startIndexInListPeer uint64
	p.currentPeerIndex = &startIndexInListPeer

	p.List = make([]peer.Peer, len(peers))

	for index, cfgPeer := range peers {
		peerCopy := cfgPeer
		p.List[index] = peerCopy
	}

	return nil
}

// GetNextPeer returns the next peer in the list.
func (p *Peers) GetNextPeer() (*peer.Peer, types.ResultCode) {
	var next int

	if *p.currentPeerIndex >= uint64(len(p.List)) {
		atomic.StoreUint64(p.currentPeerIndex, uint64(0))
	} else {
		next = p.nextIndex()
	}

	l := len(p.List) + next
	for i := next; i < l; i++ {
		idx := i % len(p.List)
		atomic.StoreUint64(p.currentPeerIndex, uint64(idx))

		peer := p.List[idx]

		_, found := p.blackListNotResponsePeers.Get(peer.URI)
		if found {
			continue
		}

		return &peer, types.ResultOK
	}

	return nil, types.ErrCantFindActivePeers
}

// AddToCacheBadPeer.
func (p *Peers) AddToCacheBadPeer(uri string) {
	_ = p.blackListNotResponsePeers.Add(uri, true)
}

// nextIndex returns the next index in a list of peers.
func (p *Peers) nextIndex() int {
	return int(atomic.AddUint64(p.currentPeerIndex, uint64(1)) % uint64(len(p.List)))
}

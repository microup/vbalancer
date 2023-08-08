package peers

import (
	"sync/atomic"

	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/types"
)

// Peers define a struct that contains a list of peers.
type Peers struct {
	List             []peer.IPeer
	CurrentPeerIndex *uint64
}

// newPeerList is the function that creates a list of peers for the balancer.
func New() *Peers {
	var startIndexInListPeer uint64

	return &Peers{
		List:             []peer.IPeer{},
		CurrentPeerIndex: &startIndexInListPeer,
	}
}

func (p *Peers) Init(peers []peer.Peer) error {
	p.List = make([]peer.IPeer, len(peers))

	for index, cfgPeer := range peers {
		peerCopy := cfgPeer
		p.List[index] = &peerCopy
	}

	return nil
}

// GetNextPeer returns the next peer in the list.
func (p *Peers) GetNextPeer() (*peer.Peer, types.ResultCode) {
	var next int

	if *p.CurrentPeerIndex >= uint64(len(p.List)) {
		atomic.StoreUint64(p.CurrentPeerIndex, uint64(0))
	} else {
		next = p.nextIndex()
	}

	l := len(p.List) + next
	for i := next; i < l; i++ {
		idx := i % len(p.List)
		atomic.StoreUint64(p.CurrentPeerIndex, uint64(idx))

		peerValue, ok := p.List[idx].(*peer.Peer)
		if !ok {
			continue
		}

		return peerValue, types.ResultOK
	}

	return nil, types.ErrCantFindActivePeers
}

// nextIndex returns the next index in a list of peers.
func (p *Peers) nextIndex() int {
	return int(atomic.AddUint64(p.CurrentPeerIndex, uint64(1)) % uint64(len(p.List)))
}

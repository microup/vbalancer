package peers

import (
	"sync/atomic"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/types"
)

// Peers define a struct that contains a list of peers and a pointer to the current peer index.
type Peers struct {
	List             []peer.IPeer
	CurrentPeerIndex *uint64
}

// New creates a new instance of Peers.
func New(list []peer.IPeer) *Peers {
	var startIndexInListPeer uint64

	return &Peers{
		List:             list,
		CurrentPeerIndex: &startIndexInListPeer,
	}
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
		peerValue, _ := p.List[idx].(*peer.Peer)

		return peerValue, types.ResultOK
	}

	// nextIndex returns the next index in the list
	return nil, types.ErrCantFindActivePeers
}

// nextIndex returns the next index in a list of peers.
func (p *Peers) nextIndex() int {
	return int(atomic.AddUint64(p.CurrentPeerIndex, uint64(1)) % uint64(len(p.List)))
}

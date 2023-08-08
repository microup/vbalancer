package peers

import (
	"sync/atomic"

	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/types"
)

// Peers define a struct that contains a list of peers.
type Peers struct {
	CurrentPeerIndex *uint64
	List             []peer.Peer `yaml:"list" json:"list"` 
}

// Initialize Peers struct with a slice of Peer objects,
// copy peers from input to the new slice, set CurrentPeerIndex to 0
// to track the selected peer's index in the slice.
func (p *Peers) Init(peers []peer.Peer) error {
	var startIndexInListPeer uint64
	p.CurrentPeerIndex = &startIndexInListPeer

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

	if *p.CurrentPeerIndex >= uint64(len(p.List)) {
		atomic.StoreUint64(p.CurrentPeerIndex, uint64(0))
	} else {
		next = p.nextIndex()
	}

	l := len(p.List) + next
	for i := next; i < l; i++ {
		idx := i % len(p.List)
		atomic.StoreUint64(p.CurrentPeerIndex, uint64(idx))

		return &p.List[idx], types.ResultOK
	}

	return nil, types.ErrCantFindActivePeers
}

// nextIndex returns the next index in a list of peers.
func (p *Peers) nextIndex() int {
	return int(atomic.AddUint64(p.CurrentPeerIndex, uint64(1)) % uint64(len(p.List)))
}

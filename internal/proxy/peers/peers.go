package peers

import (
	"sync/atomic"
	"vbalancer/internal/types"
	"vbalancer/internal/proxy/peer"	
)

type Peers struct {
	List            []*peer.Peer
	CurrentPeerIndex *uint64
}

func New(list []*peer.Peer) *Peers {
	var startIndexInListPeer uint64 

	return &Peers{
		List: list,
		CurrentPeerIndex: &startIndexInListPeer,
	}
}

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
		isAlive := p.List[idx].IsAlive()

		if isAlive {
			atomic.StoreUint64(p.CurrentPeerIndex, uint64(idx))

			return p.List[idx], types.ResultOK
		}
	}

	return nil, types.ErrCantFindActivePeers
}

func (p *Peers) nextIndex() int {
	return int(atomic.AddUint64(p.CurrentPeerIndex, uint64(1)) % uint64(len(p.List)))
}

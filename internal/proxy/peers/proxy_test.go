package peers_test

import (
	"testing"
	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
)

//nolint:funlen
func Test_API_Get_Next_Peer(t *testing.T) {
	t.Parallel()

	cases := []struct {
		nameTest          string
		peers             []*peer.Peer
		currentPeerIndex  uint64
		wantNextPeerIndex uint64
	}{
		{
			nameTest:          "test_9",
			peers:             []*peer.Peer{},
			currentPeerIndex:  0,
			wantNextPeerIndex: 0,
		},
		{
			nameTest: "test_8",
			peers: []*peer.Peer{
				{Alive: false}, {Alive: true}, {Alive: false},
			},
			currentPeerIndex:  1,
			wantNextPeerIndex: 1,
		},
		{
			nameTest: "test_7",
			peers: []*peer.Peer{
				{Alive: true}, {Alive: false}, {Alive: false},
			},
			currentPeerIndex:  0,
			wantNextPeerIndex: 0,
		},
		{
			nameTest: "test_6",
			peers: []*peer.Peer{
				{Alive: false}, {Alive: false}, {Alive: false},
			},
			currentPeerIndex:  0,
			wantNextPeerIndex: 1,
		},
		{
			nameTest: "test_5",
			peers: []*peer.Peer{
				{Alive: false}, {Alive: false}, {Alive: true},
			},
			currentPeerIndex:  0,
			wantNextPeerIndex: 2,
		},
		{
			nameTest: "test_4",
			peers: []*peer.Peer{
				{Alive: true}, {Alive: false}, {Alive: true},
			},
			currentPeerIndex:  0,
			wantNextPeerIndex: 2,
		},
		{
			nameTest: "test_3",
			peers: []*peer.Peer{
				{Alive: true}, {Alive: true}, {Alive: true},
			},
			currentPeerIndex:  0,
			wantNextPeerIndex: 1,
		},
		{
			nameTest: "test_2",
			peers: []*peer.Peer{
				{Alive: true}, {Alive: true}, {Alive: true},
			},
			currentPeerIndex:  2,
			wantNextPeerIndex: 0,
		},
		{
			nameTest: "test_1",
			peers: []*peer.Peer{
				{Alive: true}, {Alive: true}, {Alive: true},
			},
			currentPeerIndex:  3,
			wantNextPeerIndex: 0,
		},
	}

	peers := peers.New(nil)

	//nolint:varnamelen
	for _, c := range cases {
		peers.List = c.peers
		//nolint:exportloopref
		peers.CurrentPeerIndex = &c.currentPeerIndex

		_, _ = peers.GetNextPeer()

		if *peers.CurrentPeerIndex != c.wantNextPeerIndex {
			t.Errorf("Test: %s | Result failed. got %d, want: %d", c.nameTest, *peers.CurrentPeerIndex, c.wantNextPeerIndex)
		}
	}
}

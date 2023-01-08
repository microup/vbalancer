package peers_test

import (
	"sync"
	"testing"

	"vbalancer/internal/proxy/peer"
	"vbalancer/internal/proxy/peers"
)

//nolint:funlen
func Test_API_Get_Next_Peer(t *testing.T) {
	t.Parallel()

	cases := []struct {
		nameTest          string
		isAlive           []bool
		currentPeerIndex  uint64
		wantNextPeerIndex uint64
	}{
		{
			nameTest: "test_9",
			isAlive:  []bool{},
			//testListPeers:             []peer.Peer{},
			currentPeerIndex:  0,
			wantNextPeerIndex: 0,
		},
		{
			nameTest:          "test_8",
			isAlive:           []bool{false, true, false},
			currentPeerIndex:  1,
			wantNextPeerIndex: 1,
		},
		{
			nameTest:          "test_7",
			isAlive:           []bool{true, false, false},
			currentPeerIndex:  0,
			wantNextPeerIndex: 0,
		},
		{
			nameTest:          "test_6",
			isAlive:           []bool{false, false, false},
			currentPeerIndex:  0,
			wantNextPeerIndex: 1,
		},
		{
			nameTest:          "test_5",
			isAlive:           []bool{false, false, true},
			currentPeerIndex:  0,
			wantNextPeerIndex: 2,
		},
		{
			nameTest:          "test_4",
			isAlive:           []bool{true, false, true},
			currentPeerIndex:  0,
			wantNextPeerIndex: 2,
		},
		{
			nameTest:          "test_3",
			isAlive:           []bool{true, true, true},
			currentPeerIndex:  0,
			wantNextPeerIndex: 1,
		},
		{
			nameTest:          "test_2",
			isAlive:           []bool{true, true, true},
			currentPeerIndex:  2,
			wantNextPeerIndex: 0,
		},
		{
			nameTest:          "test_1",
			isAlive:           []bool{true, true, true},
			currentPeerIndex:  3,
			wantNextPeerIndex: 0,
		},
	}

	testListPeers := peers.New(nil)

	//nolint:varnamelen
	for _, c := range cases {
		var listPeers []peer.IPeer

		for _, valIsAlive := range c.isAlive {
			pPeer := &peer.Peer{
				Name:  "",
				Proto: "",
				URI:   "",
				Mu:    &sync.RWMutex{},
			}
			pPeer.Mu = &sync.RWMutex{}
			pPeer.SetAlive(valIsAlive)
			listPeers = append(listPeers, pPeer)
		}

		testListPeers.List = listPeers
		//nolint:exportloopref
		testListPeers.CurrentPeerIndex = &c.currentPeerIndex

		_, _ = testListPeers.GetNextPeer()

		if *testListPeers.CurrentPeerIndex != c.wantNextPeerIndex {
			t.Errorf("Test: %s | Result failed. got %d, want: %d",
				c.nameTest, *testListPeers.CurrentPeerIndex, c.wantNextPeerIndex)
		}
	}
}

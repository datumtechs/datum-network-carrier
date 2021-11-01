package p2p

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

var indexMock = map[int][]int{
	1:  []int{2},
	2:  []int{3},
	3:  []int{4},
	4:  []int{5},
	5:  []int{6},
	6:  []int{7},
	7:  []int{8},
	8:  []int{9},
	9:  []int{10},
	10: []int{},
}

// MockStaticNode returns to a specific network topology.
func MockStaticNode(localPeerId peer.ID, staticAddrs []string) []string {
	mockNodes := make([]string, 0)
	//
	ok, idxs := needAdd(localPeerId, staticAddrs)
	for idx, n := range staticAddrs {
		if idxs == nil {
			break
		}
		for _, i := range idxs {
			if ok && i == (idx+1) {
				mockNodes = append(mockNodes, n)
				break
			}
		}
	}
	return mockNodes
}

// mock
func needAdd(self peer.ID, staticAddrs []string) (bool, []int) {
	selfIndex := -1
	for idx, staticAddr := range staticAddrs {
		addr, _ := peersFromStringAddrs([]string{staticAddr})
		addrInfos, err := peer.AddrInfosFromP2pAddrs(addr...)
		if err != nil {
			return false, nil
		}
		if addrInfos[0].ID.String() == self.String() {
			selfIndex = idx
			break
		}
	}
	if selfIndex == -1 {
		return false, nil
	}
	selfIndex++
	return true, indexMock[selfIndex]
}


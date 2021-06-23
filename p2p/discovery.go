package p2p

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"net"
)

// Listener defines the discovery V5 network interface that is used
// to communicate with other peers.
type Listener interface {
	Self() *enode.Node
	Close()
	Lookup(enode.ID) []*enode.Node
	Resolve(*enode.Node) *enode.Node
	RandomNodes() enode.Iterator
	Ping(*enode.Node) error
	RequestENR(*enode.Node) (*enode.Node, error)
	LocalNode() *enode.LocalNode
}
//
//
//// RefreshENR uses an epoch to refresh the enr entry for our node
//// with the tracked committee ids for the epoch, allowing our node
//// to be dynamically discoverable by others given our tracked committee ids.
func (s *Service) RefreshENR() {
}
//
//// listen for new nodes watches for new nodes in the network and adds them to the peerstore.
func (s *Service) listenForNewNodes() {

}
//
func (s *Service) createListener(
	ipAddr net.IP,
	privKey *ecdsa.PrivateKey,
) (*discover.UDPv5, error) {

	return nil, nil
}

func (s *Service) startDiscoveryV5(
	addr net.IP,
	privKey *ecdsa.PrivateKey,
) (*discover.UDPv5, error) {
	listener, err := s.createListener(addr, privKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not create listener")
	}
	record := listener.Self()
	log.WithField("ENR", record.String()).Info("Started discovery v5")
	return listener, nil
}

//// This checks our set max peers in our config, and
//// determines whether our currently connected and
//// active peers are above our set max peer limit.
func (s *Service) isPeerAtLimit(inbound bool) bool {

	return false
}
//
func parseBootStrapAddrs(addrs []string) (discv5Nodes []string) {

	return nil
}
//

func parseGenericAddrs(addrs []string) (enodeString, multiAddrString []string) {
	for _, addr := range addrs {
		if addr == "" {
			continue
		}
		_, err := enode.Parse(enode.ValidSchemes, addr)
		if err == nil {
			enodeString = append(enodeString, addr)
			continue
		}
		_, err = multiAddrFromString(addr)
		if err == nil {
			multiAddrString = append(multiAddrString, addr)
			continue
		}
		log.Errorf("Invalid address of %s provided: %v", addr, err)
	}
	return enodeString, multiAddrString
}


func convertToMultiAddr(nodes []*enode.Node) []ma.Multiaddr {

	return nil
}


func convertToSingleMultiAddr(node *enode.Node) (ma.Multiaddr, error) {
	pubkey := node.Pubkey()
	assertedKey := convertToInterfacePubkey(pubkey)
	id, err := peer.IDFromPublicKey(assertedKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not get peer id")
	}
	return multiAddressBuilderWithID(node.IP().String(), "tcp", uint(node.TCP()), id)
}

//
func convertToUdpMultiAddr(node *enode.Node) ([]ma.Multiaddr, error) {

	return nil, nil
}
//
func (s *Service) filterPeer(node *enode.Node) bool { return false}

func peersFromStringAddrs(addrs []string) ([]ma.Multiaddr, error) {
	var allAddrs []ma.Multiaddr
	enodeString, multiAddrString := parseGenericAddrs(addrs)
	for _, stringAddr := range multiAddrString {
		addr, err := multiAddrFromString(stringAddr)
		if err != nil {
			return nil, errors.Wrap(err, "Could not get multiaddr from string")
		}
		allAddrs = append(allAddrs, addr)
	}
	for _, stringAddr := range enodeString {
		enodeAddr, err := enode.Parse(enode.ValidSchemes, stringAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not get enode from string")
		}
		addr, err := convertToSingleMultiAddr(enodeAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not get multiaddr")
		}
		allAddrs = append(allAddrs, addr)
	}
	return allAddrs, nil
}

func convertToAddrInfo(node *enode.Node) (*peer.AddrInfo, ma.Multiaddr, error) {
	return nil, nil, nil
}

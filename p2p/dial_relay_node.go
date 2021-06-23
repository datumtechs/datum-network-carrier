package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

func MakePeer(addr string) (*peer.AddrInfo, error) {
	maddr, err := multiAddrFromString(addr)
	if err != nil {
		return nil, err
	}
	return peer.AddrInfoFromP2pAddr(maddr)
}

func dialRelayNode(ctx context.Context, h host.Host, relayAddr string) error {
	//TODO: span...
	p, err := MakePeer(relayAddr)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, maxDialTimeout)
	defer cancel()
	return h.Connect(ctx, *p)
}

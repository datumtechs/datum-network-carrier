package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

// Send a message to a specific peer. The returned stream may be used for reading,
// but has been closed for writing.
//
// When done, the caller must Close or Reset on the stream.
func (s *Service) Send(ctx context.Context, message interface{}, baseTopic string, pid peer.ID) (network.Stream, error) {
	if err := VerifyTopicMapping(baseTopic, message); err != nil {
		return nil, err
	}
	topic := baseTopic + s.Encoding().ProtocolSuffix()
	//TODO: span...
	// Apply max dial timeout when opening a new stream.
	ctx, cancel := context.WithTimeout(ctx, maxDialTimeout)
	defer cancel()

	stream, err := s.host.NewStream(ctx, pid, protocol.ID(topic))
	if err != nil {
		return nil, err
	}
	// do not encode anything if we are sending a metadata request
	if baseTopic != RPCMetaDataTopic {
		if _, err := s.Encoding().EncodeWithMaxLength(stream, message); err != nil {
			_err := stream.Reset()
			_ = _err
			return nil, err
		}
	}
	// Close stream for writing
	if err := stream.CloseWrite(); err != nil {
		_err := stream.Reset()
		_ = _err
		return nil, err
	}
	return stream, nil
}
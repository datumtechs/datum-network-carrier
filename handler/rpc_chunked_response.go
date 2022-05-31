package handler

import (
	"errors"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/datumtechs/datum-network-carrier/p2p/encoder"
	libp2pcore "github.com/libp2p/go-libp2p-core"
)

// chunkWriter writes the given message as a chunked response to the given network stream.
// response_chunk ::= | <result> | <encoding-dependent-header>  | <encoded-payload>
func (s *Service) chunkWriter(stream libp2pcore.Stream, msg interface{}) error {
	SetStreamWriteDeadline(stream, defaultWriteDuration)
	return WriteChunk(stream, s.cfg.P2P.Encoding(), msg)
}

// WriteChunk object to stream.
// response_chunk ::= | <result> | <encoding-dependent-header> | <encoded-payload>
func WriteChunk(stream libp2pcore.Stream, encoding encoder.NetworkEncoding, msg interface{}) error {
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		return err
	}
	_, err := encoding.EncodeWithMaxLength(stream, msg)
	return err
}

// ReadChunkedBlock handles each response chunk that is sent by the
// peer and converts it into a beacon block.
func ReadChunkedBlock(stream libp2pcore.Stream, p2p p2p.P2P, isFirstChunk bool) (*carriertypespb.BlockData, error) {
	// Handle deadlines differently for first chunk
	if isFirstChunk {
		return readFirstChunkedBlock(stream, p2p)
	}
	blk := &carriertypespb.BlockData{}
	if err := readResponseChunk(stream, p2p, blk); err != nil {
		return nil, err
	}
	return blk, nil
}

// readFirstChunkedBlock reads the first chunked block and applies the appropriate deadlines to
// it.
func readFirstChunkedBlock(stream libp2pcore.Stream, p2p p2p.P2P) (*carriertypespb.BlockData, error) {
	blk := &carriertypespb.BlockData{}
	code, errMsg, err := ReadStatusCode(stream, p2p.Encoding())
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, errors.New(errMsg)
	}
	err = p2p.Encoding().DecodeWithMaxLength(stream, blk)
	return blk, err
}

// readResponseChunk reads the response from the stream and decodes it into the
// provided message type.
func readResponseChunk(stream libp2pcore.Stream, p2p p2p.P2P, to interface{}) error {
	SetStreamReadDeadline(stream, respTimeout)
	code, errMsg, err := readStatusCodeNoDeadline(stream, p2p.Encoding())
	if err != nil {
		return err
	}

	if code != 0 {
		return errors.New(errMsg)
	}
	return p2p.Encoding().DecodeWithMaxLength(stream, to)
}

func ReadChunkedGossipTestData(stream libp2pcore.Stream, p2p p2p.P2P, isFirstChunk bool) (*carrierrpcdebugpbv1.SignedGossipTestData, error) {
	if isFirstChunk {
		return readFirstChunkedGossipTestData(stream, p2p)
	}
	blk := &carrierrpcdebugpbv1.SignedGossipTestData{}
	if err := readResponseChunk(stream, p2p, blk); err != nil {
		return nil, err
	}
	return blk, nil
}

func readFirstChunkedGossipTestData(stream libp2pcore.Stream, p2p p2p.P2P) (*carrierrpcdebugpbv1.SignedGossipTestData, error) {
	blk := &carrierrpcdebugpbv1.SignedGossipTestData{}
	code, errMsg, err := ReadStatusCode(stream, p2p.Encoding())
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, errors.New(errMsg)
	}
	err = p2p.Encoding().DecodeWithMaxLength(stream, blk)
	return blk, err
}
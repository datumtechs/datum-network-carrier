package encoder_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	librpcpb "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/RosettaFlow/Carrier-Go/p2p/encoder"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"io"
	"math"
	"testing"
)

func TestSszNetworkEncoder_RoundTrip(t *testing.T) {
	e := &encoder.SszNetworkEncoder{}
	testRoundTripWithLength(t, e)
	testRoundTripWithGossip(t, e)
}

func TestSszNetworkEncoder_FailsSnappyLength(t *testing.T) {
	e := &encoder.SszNetworkEncoder{}
	goosipData := &librpcpb.GossipTestData{}
	data := make([]byte, 32)
	binary.PutUvarint(data, encoder.MaxGossipSize+32)
	err := e.DecodeGossip(data, goosipData)
	require.Contains(t, err.Error(), "snappy message exceeds max size")
}

func testRoundTripWithLength(t *testing.T, e *encoder.SszNetworkEncoder) {
	buf := new(bytes.Buffer)
	msg := &librpcpb.GossipTestData{
		Data:                 []byte("my gossip test data"),
		Count:                2,
		Step:                 3,
	}
	_, err := e.EncodeWithMaxLength(buf, msg)
	require.NoError(t, err)
	decoded := &librpcpb.GossipTestData{}
	require.NoError(t, e.DecodeWithMaxLength(buf, decoded))
	if !proto.Equal(decoded, msg) {
		t.Logf("decoded=%+v\n", decoded)
		t.Error("Decoded message is not the same as original")
	}
}

func testRoundTripWithGossip(t *testing.T, e *encoder.SszNetworkEncoder) {
	buf := new(bytes.Buffer)
	msg := &librpcpb.GossipTestData{
		Data:                 []byte("my gossip test data"),
		Count:                2,
		Step:                 3,
	}
	_, err := e.EncodeGossip(buf, msg)
	require.NoError(t, err)
	decoded := &librpcpb.GossipTestData{}
	require.NoError(t, e.DecodeGossip(buf.Bytes(), decoded))
	if !proto.Equal(decoded, msg) {
		t.Logf("decoded=%+v\n", decoded)
		t.Error("Decoded message is not the same as original")
	}
}

func TestSszNetworkEncoder_EncodeWithMaxLength(t *testing.T) {
	buf := new(bytes.Buffer)
	msg := &librpcpb.GossipTestData{
		Data:                 []byte("my gossip test data"),
		Count:                2,
		Step:                 3,
	}
	e := &encoder.SszNetworkEncoder{}
	_, err := e.EncodeWithMaxLength(buf, msg)
	assert.NilError(t, err)
}

func TestSszNetworkEncoder_DecodeWithMaxLength(t *testing.T) {
	buf := new(bytes.Buffer)
	msg := &librpcpb.GossipTestData{
		Data:                 []byte("my gossip test data"),
		Count:                2,
		Step:                 3,
	}
	e := &encoder.SszNetworkEncoder{}
	maxChunkSize := uint64(5)
	_, err := e.EncodeGossip(buf, msg)
	require.NoError(t, err)
	decoded := &librpcpb.GossipTestData{}
	err = e.DecodeWithMaxLength(buf, decoded)
	wanted := fmt.Sprintf("goes over the provided max limit of %d", maxChunkSize)
	assert.ErrorContains(t, err, wanted,)
}

func TestSszNetworkEncoder_DecodeWithMultipleFrames(t *testing.T) {

}
func TestSszNetworkEncoder_NegativeMaxLength(t *testing.T) {
	e := &encoder.SszNetworkEncoder{}
	length, err := e.MaxLength(0xfffffffffff)

	assert.Equal(t, 0, length, "Received non zero length on bad message length")
	assert.ErrorContains(t, err,"max encoded length is negative")
}

func TestSszNetworkEncoder_MaxInt64(t *testing.T) {
	e := &encoder.SszNetworkEncoder{}
	length, err := e.MaxLength(math.MaxInt64 + 1)

	assert.Equal(t, 0, length, "Received non zero length on bad message length")
	assert.ErrorContains(t, err, "invalid length provided")
}

func TestSszNetworkEncoder_DecodeWithBadSnappyStream(t *testing.T) {
	st := newBadSnappyStream()
	e := &encoder.SszNetworkEncoder{}
	decoded := new(librpcpb.GossipTestData)
	err := e.DecodeWithMaxLength(st, decoded)
	assert.ErrorContains(t, err, io.EOF.Error())
}

type badSnappyStream struct {
	varint []byte
	header []byte
	repeat []byte
	i      int
	// count how many times it was read
	counter int
	// count bytes read so far
	total int
}

func newBadSnappyStream() *badSnappyStream {
	const (
		magicBody  = "sNaPpY"
		magicChunk = "\xff\x06\x00\x00" + magicBody
	)

	header := make([]byte, len(magicChunk))
	// magicChunk == chunkTypeStreamIdentifier byte ++ 3 byte little endian len(magic body) ++ 6 byte magic body

	// header is a special chunk type, with small fixed length, to add some magic to claim it's really snappy.
	copy(header, magicChunk) // snappy library constants help us construct the common header chunk easily.

	payload := make([]byte, 4)

	// byte 0 is chunk type
	// Exploit any fancy ignored chunk type
	//   Section 4.4 Padding (chunk type 0xfe).
	//   Section 4.6. Reserved skippable chunks (chunk types 0x80-0xfd).
	payload[0] = 0xfe

	// byte 1,2,3 are chunk length (little endian)
	payload[1] = 0
	payload[2] = 0
	payload[3] = 0

	return &badSnappyStream{
		varint:  proto.EncodeVarint(1000),
		header:  header,
		repeat:  payload,
		i:       0,
		counter: 0,
		total:   0,
	}
}

func (b *badSnappyStream) Read(p []byte) (n int, err error) {
	// Stream out varint bytes first to make test happy.
	if len(b.varint) > 0 {
		copy(p, b.varint[:1])
		b.varint = b.varint[1:]
		return 1, nil
	}
	defer func() {
		b.counter += 1
		b.total += n
	}()
	if len(b.repeat) == 0 {
		panic("no bytes to repeat")
	}
	if len(b.header) > 0 {
		n = copy(p, b.header)
		b.header = b.header[n:]
		return
	}
	for n < len(p) {
		n += copy(p[n:], b.repeat[b.i:])
		b.i = (b.i + n) % len(b.repeat)
	}
	return
}

package network

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/p2p/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"testing"
)

func TestRegularSync_generateErrorResponse(t *testing.T) {
	r := &Service{
		//cfg: &Config{P2P: p2ptest.NewTestP2P(t)},
	}
	data, err := r.generateErrorResponse(responseCodeServerError, "something bad happened")
	require.NoError(t, err)

	buf := bytes.NewBuffer(data)
	b := make([]byte, 1)
	_, err = buf.Read(b)
	require.NoError(t, err)
	assert.Equal(t, responseCodeServerError, b[0], "The first byte was not the status code")
	msg := &types.ErrorMessage{}
	require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(buf, msg))
	assert.Equal(t, "something bad happened", string(*msg), "Received the wrong message")
}

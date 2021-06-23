package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"testing"
)

func TestErrorResponse_Limit(t *testing.T) {
	errorMessage := make([]byte, 0)
	// Provide a message of size 6400 bytes.
	for i := uint64(0); i < 200; i++ {
		byteArr := [32]byte{byte(i)}
		errorMessage = append(errorMessage, byteArr[:]...)
	}
	errMsg := ErrorMessage{}
	assert.ErrorContains(t, errMsg.UnmarshalSSZ(errorMessage), "expected buffer with length of upto")
}

func TestRoundTripSerialization(t *testing.T) {
	roundTripTestErrorMessage(t)
}


func roundTripTestErrorMessage(t *testing.T) {
	errMsg := []byte{'e', 'r', 'r', 'o', 'r'}
	sszErr := make(ErrorMessage, len(errMsg))
	copy(sszErr, errMsg)

	marshalledObj, err := sszErr.MarshalSSZ()
	require.NoError(t, err)
	newVal := ErrorMessage(nil)

	require.NoError(t, newVal.UnmarshalSSZ(marshalledObj))
	assert.DeepEqual(t, []byte(newVal), errMsg)
}

func TestSSZBytes_HashTreeRoot(t *testing.T) {
	tests := []struct {
		name        string
		actualValue []byte
		root        []byte
		wantErr     bool
	}{
		{
			name:        "random1",
			actualValue: hexDecodeOrDie(t, "844e1063e0b396eed17be8eddb7eecd1fe3ea46542a4b72f7466e77325e5aa6d"),
			root:        hexDecodeOrDie(t, "844e1063e0b396eed17be8eddb7eecd1fe3ea46542a4b72f7466e77325e5aa6d"),
			wantErr:     false,
		},
		{
			name:        "random1",
			actualValue: hexDecodeOrDie(t, "7b16162ecd9a28fa80a475080b0e4fff4c27efe19ce5134ce3554b72274d59fd534400ba4c7f699aa1c307cd37c2b103"),
			root:        hexDecodeOrDie(t, "128ed34ee798b9f00716f9ba5c000df5c99443dabc4d3f2e9bb86c77c732e007"),
			wantErr:     false,
		},
		{
			name:        "random2",
			actualValue: []byte{},
			root:        hexDecodeOrDie(t, "0000000000000000000000000000000000000000000000000000000000000000"),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SSZBytes(tt.actualValue)
			htr, err := s.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, tt.root, htr[:])
		})
	}
}

func hexDecodeOrDie(t *testing.T, str string) []byte {
	decoded, err := hex.DecodeString(str)
	require.NoError(t, err)
	return decoded
}

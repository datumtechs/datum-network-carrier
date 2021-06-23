package types

import (
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

const rootLength = 32

const maxErrorLength = 256

// SSZBytes is a bytes slice that satisfies the fast-ssz interface.
type SSZBytes []byte

// HashTreeRoot hashes the unit64 object following the SSZ standard.
func (b *SSZBytes) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

// HashTreeRootWith hashes the uint64 object with the given hasher.
func (b *SSZBytes) HashTreeRootWith(hh *ssz.Hasher) error {
	indx := hh.Index()
	hh.PutBytes(*b)
	hh.Merkleize(indx)
	return nil
}

// ErrorMessage describes the error message type.
type ErrorMessage []byte

// MarshalSSZTo marshals the error message with the provided byte slice.
func (m *ErrorMessage) MarshalSSZTo(dst []byte) ([]byte, error) {
	marshalledObj, err := m.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	return append(dst, marshalledObj...), nil
}

// MarshalSSZ marshals the error message into the serialized object.
func (m *ErrorMessage) MarshalSSZ() ([]byte, error) {
	if len(*m) > maxErrorLength {
		return nil, errors.Errorf("error message exceeds max size: %d > %d", len(*m), maxErrorLength)
	}
	buf := make([]byte, m.SizeSSZ())
	copy(buf, *m)
	return buf, nil
}

// SizeSSZ returns the size of the serialized representation.
func (m *ErrorMessage) SizeSSZ() int {
	return len(*m)
}

// UnmarshalSSZ unmarshals the provided bytes buffer into the
// error message object.
func (m *ErrorMessage) UnmarshalSSZ(buf []byte) error {
	bufLen := len(buf)
	maxLength := maxErrorLength
	if bufLen > maxLength {
		return errors.Errorf("expected buffer with length of upto %d but received length %d", maxLength, bufLen)
	}
	errMsg := make([]byte, bufLen)
	copy(errMsg, buf)
	*m = errMsg
	return nil
}


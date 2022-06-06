// Code generated by fastssz. DO NOT EDIT.
package v1

import (
	ssz "github.com/ferranbt/fastssz"
	github_com_prysmaticlabs_eth2_types "github.com/prysmaticlabs/eth2-types"
)

// MarshalSSZ ssz marshals the Status object
func (s *Status) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(s)
}

// MarshalSSZTo ssz marshals the Status object to a target array
func (s *Status) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'ForkDigest'
	if len(s.ForkDigest) != 4 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, s.ForkDigest...)

	// Field (1) 'FinalizedRoot'
	if len(s.FinalizedRoot) != 32 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, s.FinalizedRoot...)

	// Field (2) 'FinalizedEpoch'
	dst = ssz.MarshalUint64(dst, uint64(s.FinalizedEpoch))

	// Field (3) 'HeadRoot'
	if len(s.HeadRoot) != 32 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, s.HeadRoot...)

	// Field (4) 'HeadSlot'
	dst = ssz.MarshalUint64(dst, uint64(s.HeadSlot))

	return
}

// UnmarshalSSZ ssz unmarshals the Status object
func (s *Status) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 84 {
		return ssz.ErrSize
	}

	// Field (0) 'ForkDigest'
	if cap(s.ForkDigest) == 0 {
		s.ForkDigest = make([]byte, 0, len(buf[0:4]))
	}
	s.ForkDigest = append(s.ForkDigest, buf[0:4]...)

	// Field (1) 'FinalizedRoot'
	if cap(s.FinalizedRoot) == 0 {
		s.FinalizedRoot = make([]byte, 0, len(buf[4:36]))
	}
	s.FinalizedRoot = append(s.FinalizedRoot, buf[4:36]...)

	// Field (2) 'FinalizedEpoch'
	s.FinalizedEpoch = github_com_prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(buf[36:44]))

	// Field (3) 'HeadRoot'
	if cap(s.HeadRoot) == 0 {
		s.HeadRoot = make([]byte, 0, len(buf[44:76]))
	}
	s.HeadRoot = append(s.HeadRoot, buf[44:76]...)

	// Field (4) 'HeadSlot'
	s.HeadSlot = github_com_prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(buf[76:84]))

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the Status object
func (s *Status) SizeSSZ() (size int) {
	size = 84
	return
}

// HashTreeRoot ssz hashes the Status object
func (s *Status) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(s)
}

// HashTreeRootWith ssz hashes the Status object with a hasher
func (s *Status) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'ForkDigest'
	if len(s.ForkDigest) != 4 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(s.ForkDigest)

	// Field (1) 'FinalizedRoot'
	if len(s.FinalizedRoot) != 32 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(s.FinalizedRoot)

	// Field (2) 'FinalizedEpoch'
	hh.PutUint64(uint64(s.FinalizedEpoch))

	// Field (3) 'HeadRoot'
	if len(s.HeadRoot) != 32 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(s.HeadRoot)

	// Field (4) 'HeadSlot'
	hh.PutUint64(uint64(s.HeadSlot))

	hh.Merkleize(indx)
	return
}

// MarshalSSZ ssz marshals the ENRForkID object
func (e *ENRForkID) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(e)
}

// MarshalSSZTo ssz marshals the ENRForkID object to a target array
func (e *ENRForkID) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'CurrentForkDigest'
	if len(e.CurrentForkDigest) != 4 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, e.CurrentForkDigest...)

	// Field (1) 'NextForkVersion'
	if len(e.NextForkVersion) != 4 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, e.NextForkVersion...)

	// Field (2) 'NextForkEpoch'
	dst = ssz.MarshalUint64(dst, uint64(e.NextForkEpoch))

	return
}

// UnmarshalSSZ ssz unmarshals the ENRForkID object
func (e *ENRForkID) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	// Field (0) 'CurrentForkDigest'
	if cap(e.CurrentForkDigest) == 0 {
		e.CurrentForkDigest = make([]byte, 0, len(buf[0:4]))
	}
	e.CurrentForkDigest = append(e.CurrentForkDigest, buf[0:4]...)

	// Field (1) 'NextForkVersion'
	if cap(e.NextForkVersion) == 0 {
		e.NextForkVersion = make([]byte, 0, len(buf[4:8]))
	}
	e.NextForkVersion = append(e.NextForkVersion, buf[4:8]...)

	// Field (2) 'NextForkEpoch'
	e.NextForkEpoch = github_com_prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(buf[8:16]))

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the ENRForkID object
func (e *ENRForkID) SizeSSZ() (size int) {
	size = 16
	return
}

// HashTreeRoot ssz hashes the ENRForkID object
func (e *ENRForkID) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(e)
}

// HashTreeRootWith ssz hashes the ENRForkID object with a hasher
func (e *ENRForkID) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'CurrentForkDigest'
	if len(e.CurrentForkDigest) != 4 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(e.CurrentForkDigest)

	// Field (1) 'NextForkVersion'
	if len(e.NextForkVersion) != 4 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(e.NextForkVersion)

	// Field (2) 'NextForkEpoch'
	hh.PutUint64(uint64(e.NextForkEpoch))

	hh.Merkleize(indx)
	return
}

// MarshalSSZ ssz marshals the MetaData object
func (m *MetaData) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(m)
}

// MarshalSSZTo ssz marshals the MetaData object to a target array
func (m *MetaData) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'SeqNumber'
	dst = ssz.MarshalUint64(dst, m.SeqNumber)

	// Field (1) 'Attnets'
	if len(m.Attnets) != 8 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, m.Attnets...)

	return
}

// UnmarshalSSZ ssz unmarshals the MetaData object
func (m *MetaData) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	// Field (0) 'SeqNumber'
	m.SeqNumber = ssz.UnmarshallUint64(buf[0:8])

	// Field (1) 'Attnets'
	if cap(m.Attnets) == 0 {
		m.Attnets = make([]byte, 0, len(buf[8:16]))
	}
	m.Attnets = append(m.Attnets, buf[8:16]...)

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the MetaData object
func (m *MetaData) SizeSSZ() (size int) {
	size = 16
	return
}

// HashTreeRoot ssz hashes the MetaData object
func (m *MetaData) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(m)
}

// HashTreeRootWith ssz hashes the MetaData object with a hasher
func (m *MetaData) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'SeqNumber'
	hh.PutUint64(m.SeqNumber)

	// Field (1) 'Attnets'
	if len(m.Attnets) != 8 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(m.Attnets)

	hh.Merkleize(indx)
	return
}
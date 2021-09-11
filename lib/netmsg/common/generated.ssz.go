package common

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the MsgOption object
func (m *MsgOption) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(m)
}

// MarshalSSZTo ssz marshals the MsgOption object to a target array
func (m *MsgOption) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(32)

	// Offset (0) 'ProposalId'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(m.ProposalId)

	// Field (1) 'SenderRole'
	dst = ssz.MarshalUint64(dst, m.SenderRole)

	// Offset (2) 'SenderPartyId'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(m.SenderPartyId)

	// Field (3) 'ReceiverRole'
	dst = ssz.MarshalUint64(dst, m.ReceiverRole)

	// Offset (4) 'ReceiverPartyId'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(m.ReceiverPartyId)

	// Offset (5) 'MsgOwner'
	dst = ssz.WriteOffset(dst, offset)
	if m.MsgOwner == nil {
		m.MsgOwner = new(TaskOrganizationIdentityInfo)
	}
	offset += m.MsgOwner.SizeSSZ()

	// Field (0) 'ProposalId'
	if len(m.ProposalId) > 1024 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, m.ProposalId...)

	// Field (2) 'SenderPartyId'
	if len(m.SenderPartyId) > 32 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, m.SenderPartyId...)

	// Field (4) 'ReceiverPartyId'
	if len(m.ReceiverPartyId) > 32 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, m.ReceiverPartyId...)

	// Field (5) 'MsgOwner'
	if dst, err = m.MsgOwner.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the MsgOption object
func (m *MsgOption) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 32 {
		return ssz.ErrSize
	}

	tail := buf
	var o0, o2, o4, o5 uint64

	// Offset (0) 'ProposalId'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 32 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (1) 'SenderRole'
	m.SenderRole = ssz.UnmarshallUint64(buf[4:12])

	// Offset (2) 'SenderPartyId'
	if o2 = ssz.ReadOffset(buf[12:16]); o2 > size || o0 > o2 {
		return ssz.ErrOffset
	}

	// Field (3) 'ReceiverRole'
	m.ReceiverRole = ssz.UnmarshallUint64(buf[16:24])

	// Offset (4) 'ReceiverPartyId'
	if o4 = ssz.ReadOffset(buf[24:28]); o4 > size || o2 > o4 {
		return ssz.ErrOffset
	}

	// Offset (5) 'MsgOwner'
	if o5 = ssz.ReadOffset(buf[28:32]); o5 > size || o4 > o5 {
		return ssz.ErrOffset
	}

	// Field (0) 'ProposalId'
	{
		buf = tail[o0:o2]
		if len(buf) > 1024 {
			return ssz.ErrBytesLength
		}
		if cap(m.ProposalId) == 0 {
			m.ProposalId = make([]byte, 0, len(buf))
		}
		m.ProposalId = append(m.ProposalId, buf...)
	}

	// Field (2) 'SenderPartyId'
	{
		buf = tail[o2:o4]
		if len(buf) > 32 {
			return ssz.ErrBytesLength
		}
		if cap(m.SenderPartyId) == 0 {
			m.SenderPartyId = make([]byte, 0, len(buf))
		}
		m.SenderPartyId = append(m.SenderPartyId, buf...)
	}

	// Field (4) 'ReceiverPartyId'
	{
		buf = tail[o4:o5]
		if len(buf) > 32 {
			return ssz.ErrBytesLength
		}
		if cap(m.ReceiverPartyId) == 0 {
			m.ReceiverPartyId = make([]byte, 0, len(buf))
		}
		m.ReceiverPartyId = append(m.ReceiverPartyId, buf...)
	}

	// Field (5) 'MsgOwner'
	{
		buf = tail[o5:]
		if m.MsgOwner == nil {
			m.MsgOwner = new(TaskOrganizationIdentityInfo)
		}
		if err = m.MsgOwner.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the MsgOption object
func (m *MsgOption) SizeSSZ() (size int) {
	size = 32

	// Field (0) 'ProposalId'
	size += len(m.ProposalId)

	// Field (2) 'SenderPartyId'
	size += len(m.SenderPartyId)

	// Field (4) 'ReceiverPartyId'
	size += len(m.ReceiverPartyId)

	// Field (5) 'MsgOwner'
	if m.MsgOwner == nil {
		m.MsgOwner = new(TaskOrganizationIdentityInfo)
	}
	size += m.MsgOwner.SizeSSZ()

	return
}

// HashTreeRoot ssz hashes the MsgOption object
func (m *MsgOption) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(m)
}

// HashTreeRootWith ssz hashes the MsgOption object with a hasher
func (m *MsgOption) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'ProposalId'
	if len(m.ProposalId) > 1024 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(m.ProposalId)

	// Field (1) 'SenderRole'
	hh.PutUint64(m.SenderRole)

	// Field (2) 'SenderPartyId'
	if len(m.SenderPartyId) > 32 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(m.SenderPartyId)

	// Field (3) 'ReceiverRole'
	hh.PutUint64(m.ReceiverRole)

	// Field (4) 'ReceiverPartyId'
	if len(m.ReceiverPartyId) > 32 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(m.ReceiverPartyId)

	// Field (5) 'MsgOwner'
	if err = m.MsgOwner.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// MarshalSSZ ssz marshals the TaskOrganizationIdentityInfo object
func (t *TaskOrganizationIdentityInfo) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(t)
}

// MarshalSSZTo ssz marshals the TaskOrganizationIdentityInfo object to a target array
func (t *TaskOrganizationIdentityInfo) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(16)

	// Offset (0) 'Name'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(t.Name)

	// Offset (1) 'NodeId'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(t.NodeId)

	// Offset (2) 'IdentityId'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(t.IdentityId)

	// Offset (3) 'PartyId'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(t.PartyId)

	// Field (0) 'Name'
	if len(t.Name) > 64 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, t.Name...)

	// Field (1) 'NodeId'
	if len(t.NodeId) > 1024 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, t.NodeId...)

	// Field (2) 'IdentityId'
	if len(t.IdentityId) > 1024 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, t.IdentityId...)

	// Field (3) 'PartyId'
	if len(t.PartyId) > 32 {
		err = ssz.ErrBytesLength
		return
	}
	dst = append(dst, t.PartyId...)

	return
}

// UnmarshalSSZ ssz unmarshals the TaskOrganizationIdentityInfo object
func (t *TaskOrganizationIdentityInfo) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 16 {
		return ssz.ErrSize
	}

	tail := buf
	var o0, o1, o2, o3 uint64

	// Offset (0) 'Name'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 16 {
		return ssz.ErrInvalidVariableOffset
	}

	// Offset (1) 'NodeId'
	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
		return ssz.ErrOffset
	}

	// Offset (2) 'IdentityId'
	if o2 = ssz.ReadOffset(buf[8:12]); o2 > size || o1 > o2 {
		return ssz.ErrOffset
	}

	// Offset (3) 'PartyId'
	if o3 = ssz.ReadOffset(buf[12:16]); o3 > size || o2 > o3 {
		return ssz.ErrOffset
	}

	// Field (0) 'Name'
	{
		buf = tail[o0:o1]
		if len(buf) > 64 {
			return ssz.ErrBytesLength
		}
		if cap(t.Name) == 0 {
			t.Name = make([]byte, 0, len(buf))
		}
		t.Name = append(t.Name, buf...)
	}

	// Field (1) 'NodeId'
	{
		buf = tail[o1:o2]
		if len(buf) > 1024 {
			return ssz.ErrBytesLength
		}
		if cap(t.NodeId) == 0 {
			t.NodeId = make([]byte, 0, len(buf))
		}
		t.NodeId = append(t.NodeId, buf...)
	}

	// Field (2) 'IdentityId'
	{
		buf = tail[o2:o3]
		if len(buf) > 1024 {
			return ssz.ErrBytesLength
		}
		if cap(t.IdentityId) == 0 {
			t.IdentityId = make([]byte, 0, len(buf))
		}
		t.IdentityId = append(t.IdentityId, buf...)
	}

	// Field (3) 'PartyId'
	{
		buf = tail[o3:]
		if len(buf) > 32 {
			return ssz.ErrBytesLength
		}
		if cap(t.PartyId) == 0 {
			t.PartyId = make([]byte, 0, len(buf))
		}
		t.PartyId = append(t.PartyId, buf...)
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the TaskOrganizationIdentityInfo object
func (t *TaskOrganizationIdentityInfo) SizeSSZ() (size int) {
	size = 16

	// Field (0) 'Name'
	size += len(t.Name)

	// Field (1) 'NodeId'
	size += len(t.NodeId)

	// Field (2) 'IdentityId'
	size += len(t.IdentityId)

	// Field (3) 'PartyId'
	size += len(t.PartyId)

	return
}

// HashTreeRoot ssz hashes the TaskOrganizationIdentityInfo object
func (t *TaskOrganizationIdentityInfo) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(t)
}

// HashTreeRootWith ssz hashes the TaskOrganizationIdentityInfo object with a hasher
func (t *TaskOrganizationIdentityInfo) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'Name'
	if len(t.Name) > 64 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(t.Name)

	// Field (1) 'NodeId'
	if len(t.NodeId) > 1024 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(t.NodeId)

	// Field (2) 'IdentityId'
	if len(t.IdentityId) > 1024 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(t.IdentityId)

	// Field (3) 'PartyId'
	if len(t.PartyId) > 32 {
		err = ssz.ErrBytesLength
		return
	}
	hh.PutBytes(t.PartyId)

	hh.Merkleize(indx)
	return
}

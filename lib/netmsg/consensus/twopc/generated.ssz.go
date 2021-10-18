package twopc
//
//import (
//	"github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
//	ssz "github.com/ferranbt/fastssz"
//)
//
//// MarshalSSZ ssz marshals the PrepareMsg object
//func (p *PrepareMsg) MarshalSSZ() ([]byte, error) {
//	return ssz.MarshalSSZ(p)
//}
//
//// MarshalSSZTo ssz marshals the PrepareMsg object to a target array
//func (p *PrepareMsg) MarshalSSZTo(buf []byte) (dst []byte, err error) {
//	dst = buf
//	offset := int(20)
//
//	// Offset (0) 'MsgOption'
//	dst = ssz.WriteOffset(dst, offset)
//	if p.MsgOption == nil {
//		p.MsgOption = new(common.MsgOption)
//	}
//	offset += p.MsgOption.SizeSSZ()
//
//	// Offset (1) 'TaskInfo'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(p.TaskInfo)
//
//	// Field (2) 'CreateAt'
//	dst = ssz.MarshalUint64(dst, p.CreateAt)
//
//	// Offset (3) 'Sign'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(p.Sign)
//
//	// Field (0) 'MsgOption'
//	if dst, err = p.MsgOption.MarshalSSZTo(dst); err != nil {
//		return
//	}
//
//	// Field (1) 'TaskInfo'
//	if len(p.TaskInfo) > 16777216 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, p.TaskInfo...)
//
//	// Field (3) 'Sign'
//	if len(p.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, p.Sign...)
//
//	return
//}
//
//// UnmarshalSSZ ssz unmarshals the PrepareMsg object
//func (p *PrepareMsg) UnmarshalSSZ(buf []byte) error {
//	var err error
//	size := uint64(len(buf))
//	if size < 20 {
//		return ssz.ErrSize
//	}
//
//	tail := buf
//	var o0, o1, o3 uint64
//
//	// Offset (0) 'MsgOption'
//	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
//		return ssz.ErrOffset
//	}
//
//	if o0 < 20 {
//		return ssz.ErrInvalidVariableOffset
//	}
//
//	// Offset (1) 'TaskInfo'
//	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
//		return ssz.ErrOffset
//	}
//
//	// Field (2) 'CreateAt'
//	p.CreateAt = ssz.UnmarshallUint64(buf[8:16])
//
//	// Offset (3) 'Sign'
//	if o3 = ssz.ReadOffset(buf[16:20]); o3 > size || o1 > o3 {
//		return ssz.ErrOffset
//	}
//
//	// Field (0) 'MsgOption'
//	{
//		buf = tail[o0:o1]
//		if p.MsgOption == nil {
//			p.MsgOption = new(common.MsgOption)
//		}
//		if err = p.MsgOption.UnmarshalSSZ(buf); err != nil {
//			return err
//		}
//	}
//
//	// Field (1) 'TaskInfo'
//	{
//		buf = tail[o1:o3]
//		if len(buf) > 16777216 {
//			return ssz.ErrBytesLength
//		}
//		if cap(p.TaskInfo) == 0 {
//			p.TaskInfo = make([]byte, 0, len(buf))
//		}
//		p.TaskInfo = append(p.TaskInfo, buf...)
//	}
//
//	// Field (3) 'Sign'
//	{
//		buf = tail[o3:]
//		if len(buf) > 1024 {
//			return ssz.ErrBytesLength
//		}
//		if cap(p.Sign) == 0 {
//			p.Sign = make([]byte, 0, len(buf))
//		}
//		p.Sign = append(p.Sign, buf...)
//	}
//	return err
//}
//
//// SizeSSZ returns the ssz encoded size in bytes for the PrepareMsg object
//func (p *PrepareMsg) SizeSSZ() (size int) {
//	size = 20
//
//	// Field (0) 'MsgOption'
//	if p.MsgOption == nil {
//		p.MsgOption = new(common.MsgOption)
//	}
//	size += p.MsgOption.SizeSSZ()
//
//	// Field (1) 'TaskInfo'
//	size += len(p.TaskInfo)
//
//	// Field (3) 'Sign'
//	size += len(p.Sign)
//
//	return
//}
//
//// HashTreeRoot ssz hashes the PrepareMsg object
//func (p *PrepareMsg) HashTreeRoot() ([32]byte, error) {
//	return ssz.HashWithDefaultHasher(p)
//}
//
//// HashTreeRootWith ssz hashes the PrepareMsg object with a hasher
//func (p *PrepareMsg) HashTreeRootWith(hh *ssz.Hasher) (err error) {
//	indx := hh.Index()
//
//	// Field (0) 'MsgOption'
//	if err = p.MsgOption.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (1) 'TaskInfo'
//	if len(p.TaskInfo) > 16777216 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(p.TaskInfo)
//
//	// Field (2) 'CreateAt'
//	hh.PutUint64(p.CreateAt)
//
//	// Field (3) 'Sign'
//	if len(p.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(p.Sign)
//
//	hh.Merkleize(indx)
//	return
//}
//
//// MarshalSSZ ssz marshals the PrepareVote object
//func (p *PrepareVote) MarshalSSZ() ([]byte, error) {
//	return ssz.MarshalSSZ(p)
//}
//
//// MarshalSSZTo ssz marshals the PrepareVote object to a target array
//func (p *PrepareVote) MarshalSSZTo(buf []byte) (dst []byte, err error) {
//	dst = buf
//	offset := int(24)
//
//	// Offset (0) 'MsgOption'
//	dst = ssz.WriteOffset(dst, offset)
//	if p.MsgOption == nil {
//		p.MsgOption = new(common.MsgOption)
//	}
//	offset += p.MsgOption.SizeSSZ()
//
//	// Offset (1) 'VoteOption'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(p.VoteOption)
//
//	// Offset (2) 'PeerInfo'
//	dst = ssz.WriteOffset(dst, offset)
//	if p.PeerInfo == nil {
//		p.PeerInfo = new(TaskPeerInfo)
//	}
//	offset += p.PeerInfo.SizeSSZ()
//
//	// Field (3) 'CreateAt'
//	dst = ssz.MarshalUint64(dst, p.CreateAt)
//
//	// Offset (4) 'Sign'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(p.Sign)
//
//	// Field (0) 'MsgOption'
//	if dst, err = p.MsgOption.MarshalSSZTo(dst); err != nil {
//		return
//	}
//
//	// Field (1) 'VoteOption'
//	if len(p.VoteOption) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, p.VoteOption...)
//
//	// Field (2) 'PeerInfo'
//	if dst, err = p.PeerInfo.MarshalSSZTo(dst); err != nil {
//		return
//	}
//
//	// Field (4) 'Sign'
//	if len(p.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, p.Sign...)
//
//	return
//}
//
//// UnmarshalSSZ ssz unmarshals the PrepareVote object
//func (p *PrepareVote) UnmarshalSSZ(buf []byte) error {
//	var err error
//	size := uint64(len(buf))
//	if size < 24 {
//		return ssz.ErrSize
//	}
//
//	tail := buf
//	var o0, o1, o2, o4 uint64
//
//	// Offset (0) 'MsgOption'
//	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
//		return ssz.ErrOffset
//	}
//
//	if o0 < 24 {
//		return ssz.ErrInvalidVariableOffset
//	}
//
//	// Offset (1) 'VoteOption'
//	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
//		return ssz.ErrOffset
//	}
//
//	// Offset (2) 'PeerInfo'
//	if o2 = ssz.ReadOffset(buf[8:12]); o2 > size || o1 > o2 {
//		return ssz.ErrOffset
//	}
//
//	// Field (3) 'CreateAt'
//	p.CreateAt = ssz.UnmarshallUint64(buf[12:20])
//
//	// Offset (4) 'Sign'
//	if o4 = ssz.ReadOffset(buf[20:24]); o4 > size || o2 > o4 {
//		return ssz.ErrOffset
//	}
//
//	// Field (0) 'MsgOption'
//	{
//		buf = tail[o0:o1]
//		if p.MsgOption == nil {
//			p.MsgOption = new(common.MsgOption)
//		}
//		if err = p.MsgOption.UnmarshalSSZ(buf); err != nil {
//			return err
//		}
//	}
//
//	// Field (1) 'VoteOption'
//	{
//		buf = tail[o1:o2]
//		if len(buf) > 32 {
//			return ssz.ErrBytesLength
//		}
//		if cap(p.VoteOption) == 0 {
//			p.VoteOption = make([]byte, 0, len(buf))
//		}
//		p.VoteOption = append(p.VoteOption, buf...)
//	}
//
//	// Field (2) 'PeerInfo'
//	{
//		buf = tail[o2:o4]
//		if p.PeerInfo == nil {
//			p.PeerInfo = new(TaskPeerInfo)
//		}
//		if err = p.PeerInfo.UnmarshalSSZ(buf); err != nil {
//			return err
//		}
//	}
//
//	// Field (4) 'Sign'
//	{
//		buf = tail[o4:]
//		if len(buf) > 1024 {
//			return ssz.ErrBytesLength
//		}
//		if cap(p.Sign) == 0 {
//			p.Sign = make([]byte, 0, len(buf))
//		}
//		p.Sign = append(p.Sign, buf...)
//	}
//	return err
//}
//
//// SizeSSZ returns the ssz encoded size in bytes for the PrepareVote object
//func (p *PrepareVote) SizeSSZ() (size int) {
//	size = 24
//
//	// Field (0) 'MsgOption'
//	if p.MsgOption == nil {
//		p.MsgOption = new(common.MsgOption)
//	}
//	size += p.MsgOption.SizeSSZ()
//
//	// Field (1) 'VoteOption'
//	size += len(p.VoteOption)
//
//	// Field (2) 'PeerInfo'
//	if p.PeerInfo == nil {
//		p.PeerInfo = new(TaskPeerInfo)
//	}
//	size += p.PeerInfo.SizeSSZ()
//
//	// Field (4) 'Sign'
//	size += len(p.Sign)
//
//	return
//}
//
//// HashTreeRoot ssz hashes the PrepareVote object
//func (p *PrepareVote) HashTreeRoot() ([32]byte, error) {
//	return ssz.HashWithDefaultHasher(p)
//}
//
//// HashTreeRootWith ssz hashes the PrepareVote object with a hasher
//func (p *PrepareVote) HashTreeRootWith(hh *ssz.Hasher) (err error) {
//	indx := hh.Index()
//
//	// Field (0) 'MsgOption'
//	if err = p.MsgOption.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (1) 'VoteOption'
//	if len(p.VoteOption) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(p.VoteOption)
//
//	// Field (2) 'PeerInfo'
//	if err = p.PeerInfo.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (3) 'CreateAt'
//	hh.PutUint64(p.CreateAt)
//
//	// Field (4) 'Sign'
//	if len(p.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(p.Sign)
//
//	hh.Merkleize(indx)
//	return
//}
//
//// MarshalSSZ ssz marshals the ConfirmMsg object
//func (c *ConfirmMsg) MarshalSSZ() ([]byte, error) {
//	return ssz.MarshalSSZ(c)
//}
//
//// MarshalSSZTo ssz marshals the ConfirmMsg object to a target array
//func (c *ConfirmMsg) MarshalSSZTo(buf []byte) (dst []byte, err error) {
//	dst = buf
//	offset := int(24)
//
//	// Offset (0) 'MsgOption'
//	dst = ssz.WriteOffset(dst, offset)
//	if c.MsgOption == nil {
//		c.MsgOption = new(common.MsgOption)
//	}
//	offset += c.MsgOption.SizeSSZ()
//
//	// Offset (1) 'ConfirmOption'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(c.ConfirmOption)
//
//	// Offset (2) 'Peers'
//	dst = ssz.WriteOffset(dst, offset)
//	if c.Peers == nil {
//		c.Peers = new(ConfirmTaskPeerInfo)
//	}
//	offset += c.Peers.SizeSSZ()
//
//	// Field (3) 'CreateAt'
//	dst = ssz.MarshalUint64(dst, c.CreateAt)
//
//	// Offset (4) 'Sign'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(c.Sign)
//
//	// Field (0) 'MsgOption'
//	if dst, err = c.MsgOption.MarshalSSZTo(dst); err != nil {
//		return
//	}
//
//	// Field (1) 'ConfirmOption'
//	if len(c.ConfirmOption) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, c.ConfirmOption...)
//
//	// Field (2) 'Peers'
//	if dst, err = c.Peers.MarshalSSZTo(dst); err != nil {
//		return
//	}
//
//	// Field (4) 'Sign'
//	if len(c.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, c.Sign...)
//
//	return
//}
//
//// UnmarshalSSZ ssz unmarshals the ConfirmMsg object
//func (c *ConfirmMsg) UnmarshalSSZ(buf []byte) error {
//	var err error
//	size := uint64(len(buf))
//	if size < 24 {
//		return ssz.ErrSize
//	}
//
//	tail := buf
//	var o0, o1, o2, o4 uint64
//
//	// Offset (0) 'MsgOption'
//	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
//		return ssz.ErrOffset
//	}
//
//	if o0 < 24 {
//		return ssz.ErrInvalidVariableOffset
//	}
//
//	// Offset (1) 'ConfirmOption'
//	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
//		return ssz.ErrOffset
//	}
//
//	// Offset (2) 'Peers'
//	if o2 = ssz.ReadOffset(buf[8:12]); o2 > size || o1 > o2 {
//		return ssz.ErrOffset
//	}
//
//	// Field (3) 'CreateAt'
//	c.CreateAt = ssz.UnmarshallUint64(buf[12:20])
//
//	// Offset (4) 'Sign'
//	if o4 = ssz.ReadOffset(buf[20:24]); o4 > size || o2 > o4 {
//		return ssz.ErrOffset
//	}
//
//	// Field (0) 'MsgOption'
//	{
//		buf = tail[o0:o1]
//		if c.MsgOption == nil {
//			c.MsgOption = new(common.MsgOption)
//		}
//		if err = c.MsgOption.UnmarshalSSZ(buf); err != nil {
//			return err
//		}
//	}
//
//	// Field (1) 'ConfirmOption'
//	{
//		buf = tail[o1:o2]
//		if len(buf) > 32 {
//			return ssz.ErrBytesLength
//		}
//		if cap(c.ConfirmOption) == 0 {
//			c.ConfirmOption = make([]byte, 0, len(buf))
//		}
//		c.ConfirmOption = append(c.ConfirmOption, buf...)
//	}
//
//	// Field (2) 'Peers'
//	{
//		buf = tail[o2:o4]
//		if c.Peers == nil {
//			c.Peers = new(ConfirmTaskPeerInfo)
//		}
//		if err = c.Peers.UnmarshalSSZ(buf); err != nil {
//			return err
//		}
//	}
//
//	// Field (4) 'Sign'
//	{
//		buf = tail[o4:]
//		if len(buf) > 1024 {
//			return ssz.ErrBytesLength
//		}
//		if cap(c.Sign) == 0 {
//			c.Sign = make([]byte, 0, len(buf))
//		}
//		c.Sign = append(c.Sign, buf...)
//	}
//	return err
//}
//
//// SizeSSZ returns the ssz encoded size in bytes for the ConfirmMsg object
//func (c *ConfirmMsg) SizeSSZ() (size int) {
//	size = 24
//
//	// Field (0) 'MsgOption'
//	if c.MsgOption == nil {
//		c.MsgOption = new(common.MsgOption)
//	}
//	size += c.MsgOption.SizeSSZ()
//
//	// Field (1) 'ConfirmOption'
//	size += len(c.ConfirmOption)
//
//	// Field (2) 'Peers'
//	if c.Peers == nil {
//		c.Peers = new(ConfirmTaskPeerInfo)
//	}
//	size += c.Peers.SizeSSZ()
//
//	// Field (4) 'Sign'
//	size += len(c.Sign)
//
//	return
//}
//
//// HashTreeRoot ssz hashes the ConfirmMsg object
//func (c *ConfirmMsg) HashTreeRoot() ([32]byte, error) {
//	return ssz.HashWithDefaultHasher(c)
//}
//
//// HashTreeRootWith ssz hashes the ConfirmMsg object with a hasher
//func (c *ConfirmMsg) HashTreeRootWith(hh *ssz.Hasher) (err error) {
//	indx := hh.Index()
//
//	// Field (0) 'MsgOption'
//	if err = c.MsgOption.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (1) 'ConfirmOption'
//	if len(c.ConfirmOption) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(c.ConfirmOption)
//
//	// Field (2) 'Peers'
//	if err = c.Peers.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (3) 'CreateAt'
//	hh.PutUint64(c.CreateAt)
//
//	// Field (4) 'Sign'
//	if len(c.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(c.Sign)
//
//	hh.Merkleize(indx)
//	return
//}
//
//// MarshalSSZ ssz marshals the ConfirmTaskPeerInfo object
//func (c *ConfirmTaskPeerInfo) MarshalSSZ() ([]byte, error) {
//	return ssz.MarshalSSZ(c)
//}
//
//// MarshalSSZTo ssz marshals the ConfirmTaskPeerInfo object to a target array
//func (c *ConfirmTaskPeerInfo) MarshalSSZTo(buf []byte) (dst []byte, err error) {
//	dst = buf
//	offset := int(16)
//
//	// Offset (0) 'OwnerPeerInfo'
//	dst = ssz.WriteOffset(dst, offset)
//	if c.OwnerPeerInfo == nil {
//		c.OwnerPeerInfo = new(TaskPeerInfo)
//	}
//	offset += c.OwnerPeerInfo.SizeSSZ()
//
//	// Offset (1) 'DataSupplierPeerInfoList'
//	dst = ssz.WriteOffset(dst, offset)
//	for ii := 0; ii < len(c.DataSupplierPeerInfoList); ii++ {
//		offset += 4
//		offset += c.DataSupplierPeerInfoList[ii].SizeSSZ()
//	}
//
//	// Offset (2) 'PowerSupplierPeerInfoList'
//	dst = ssz.WriteOffset(dst, offset)
//	for ii := 0; ii < len(c.PowerSupplierPeerInfoList); ii++ {
//		offset += 4
//		offset += c.PowerSupplierPeerInfoList[ii].SizeSSZ()
//	}
//
//	// Offset (3) 'ResultReceiverPeerInfoList'
//	dst = ssz.WriteOffset(dst, offset)
//	for ii := 0; ii < len(c.ResultReceiverPeerInfoList); ii++ {
//		offset += 4
//		offset += c.ResultReceiverPeerInfoList[ii].SizeSSZ()
//	}
//
//	// Field (0) 'OwnerPeerInfo'
//	if dst, err = c.OwnerPeerInfo.MarshalSSZTo(dst); err != nil {
//		return
//	}
//
//	// Field (1) 'DataSupplierPeerInfoList'
//	if len(c.DataSupplierPeerInfoList) > 16777216 {
//		err = ssz.ErrListTooBig
//		return
//	}
//	{
//		offset = 4 * len(c.DataSupplierPeerInfoList)
//		for ii := 0; ii < len(c.DataSupplierPeerInfoList); ii++ {
//			dst = ssz.WriteOffset(dst, offset)
//			offset += c.DataSupplierPeerInfoList[ii].SizeSSZ()
//		}
//	}
//	for ii := 0; ii < len(c.DataSupplierPeerInfoList); ii++ {
//		if dst, err = c.DataSupplierPeerInfoList[ii].MarshalSSZTo(dst); err != nil {
//			return
//		}
//	}
//
//	// Field (2) 'PowerSupplierPeerInfoList'
//	if len(c.PowerSupplierPeerInfoList) > 16777216 {
//		err = ssz.ErrListTooBig
//		return
//	}
//	{
//		offset = 4 * len(c.PowerSupplierPeerInfoList)
//		for ii := 0; ii < len(c.PowerSupplierPeerInfoList); ii++ {
//			dst = ssz.WriteOffset(dst, offset)
//			offset += c.PowerSupplierPeerInfoList[ii].SizeSSZ()
//		}
//	}
//	for ii := 0; ii < len(c.PowerSupplierPeerInfoList); ii++ {
//		if dst, err = c.PowerSupplierPeerInfoList[ii].MarshalSSZTo(dst); err != nil {
//			return
//		}
//	}
//
//	// Field (3) 'ResultReceiverPeerInfoList'
//	if len(c.ResultReceiverPeerInfoList) > 16777216 {
//		err = ssz.ErrListTooBig
//		return
//	}
//	{
//		offset = 4 * len(c.ResultReceiverPeerInfoList)
//		for ii := 0; ii < len(c.ResultReceiverPeerInfoList); ii++ {
//			dst = ssz.WriteOffset(dst, offset)
//			offset += c.ResultReceiverPeerInfoList[ii].SizeSSZ()
//		}
//	}
//	for ii := 0; ii < len(c.ResultReceiverPeerInfoList); ii++ {
//		if dst, err = c.ResultReceiverPeerInfoList[ii].MarshalSSZTo(dst); err != nil {
//			return
//		}
//	}
//
//	return
//}
//
//// UnmarshalSSZ ssz unmarshals the ConfirmTaskPeerInfo object
//func (c *ConfirmTaskPeerInfo) UnmarshalSSZ(buf []byte) error {
//	var err error
//	size := uint64(len(buf))
//	if size < 16 {
//		return ssz.ErrSize
//	}
//
//	tail := buf
//	var o0, o1, o2, o3 uint64
//
//	// Offset (0) 'OwnerPeerInfo'
//	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
//		return ssz.ErrOffset
//	}
//
//	if o0 < 16 {
//		return ssz.ErrInvalidVariableOffset
//	}
//
//	// Offset (1) 'DataSupplierPeerInfoList'
//	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
//		return ssz.ErrOffset
//	}
//
//	// Offset (2) 'PowerSupplierPeerInfoList'
//	if o2 = ssz.ReadOffset(buf[8:12]); o2 > size || o1 > o2 {
//		return ssz.ErrOffset
//	}
//
//	// Offset (3) 'ResultReceiverPeerInfoList'
//	if o3 = ssz.ReadOffset(buf[12:16]); o3 > size || o2 > o3 {
//		return ssz.ErrOffset
//	}
//
//	// Field (0) 'OwnerPeerInfo'
//	{
//		buf = tail[o0:o1]
//		if c.OwnerPeerInfo == nil {
//			c.OwnerPeerInfo = new(TaskPeerInfo)
//		}
//		if err = c.OwnerPeerInfo.UnmarshalSSZ(buf); err != nil {
//			return err
//		}
//	}
//
//	// Field (1) 'DataSupplierPeerInfoList'
//	{
//		buf = tail[o1:o2]
//		num, err := ssz.DecodeDynamicLength(buf, 16777216)
//		if err != nil {
//			return err
//		}
//		c.DataSupplierPeerInfoList = make([]*TaskPeerInfo, num)
//		err = ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
//			if c.DataSupplierPeerInfoList[indx] == nil {
//				c.DataSupplierPeerInfoList[indx] = new(TaskPeerInfo)
//			}
//			if err = c.DataSupplierPeerInfoList[indx].UnmarshalSSZ(buf); err != nil {
//				return err
//			}
//			return nil
//		})
//		if err != nil {
//			return err
//		}
//	}
//
//	// Field (2) 'PowerSupplierPeerInfoList'
//	{
//		buf = tail[o2:o3]
//		num, err := ssz.DecodeDynamicLength(buf, 16777216)
//		if err != nil {
//			return err
//		}
//		c.PowerSupplierPeerInfoList = make([]*TaskPeerInfo, num)
//		err = ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
//			if c.PowerSupplierPeerInfoList[indx] == nil {
//				c.PowerSupplierPeerInfoList[indx] = new(TaskPeerInfo)
//			}
//			if err = c.PowerSupplierPeerInfoList[indx].UnmarshalSSZ(buf); err != nil {
//				return err
//			}
//			return nil
//		})
//		if err != nil {
//			return err
//		}
//	}
//
//	// Field (3) 'ResultReceiverPeerInfoList'
//	{
//		buf = tail[o3:]
//		num, err := ssz.DecodeDynamicLength(buf, 16777216)
//		if err != nil {
//			return err
//		}
//		c.ResultReceiverPeerInfoList = make([]*TaskPeerInfo, num)
//		err = ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
//			if c.ResultReceiverPeerInfoList[indx] == nil {
//				c.ResultReceiverPeerInfoList[indx] = new(TaskPeerInfo)
//			}
//			if err = c.ResultReceiverPeerInfoList[indx].UnmarshalSSZ(buf); err != nil {
//				return err
//			}
//			return nil
//		})
//		if err != nil {
//			return err
//		}
//	}
//	return err
//}
//
//// SizeSSZ returns the ssz encoded size in bytes for the ConfirmTaskPeerInfo object
//func (c *ConfirmTaskPeerInfo) SizeSSZ() (size int) {
//	size = 16
//
//	// Field (0) 'OwnerPeerInfo'
//	if c.OwnerPeerInfo == nil {
//		c.OwnerPeerInfo = new(TaskPeerInfo)
//	}
//	size += c.OwnerPeerInfo.SizeSSZ()
//
//	// Field (1) 'DataSupplierPeerInfoList'
//	for ii := 0; ii < len(c.DataSupplierPeerInfoList); ii++ {
//		size += 4
//		size += c.DataSupplierPeerInfoList[ii].SizeSSZ()
//	}
//
//	// Field (2) 'PowerSupplierPeerInfoList'
//	for ii := 0; ii < len(c.PowerSupplierPeerInfoList); ii++ {
//		size += 4
//		size += c.PowerSupplierPeerInfoList[ii].SizeSSZ()
//	}
//
//	// Field (3) 'ResultReceiverPeerInfoList'
//	for ii := 0; ii < len(c.ResultReceiverPeerInfoList); ii++ {
//		size += 4
//		size += c.ResultReceiverPeerInfoList[ii].SizeSSZ()
//	}
//
//	return
//}
//
//// HashTreeRoot ssz hashes the ConfirmTaskPeerInfo object
//func (c *ConfirmTaskPeerInfo) HashTreeRoot() ([32]byte, error) {
//	return ssz.HashWithDefaultHasher(c)
//}
//
//// HashTreeRootWith ssz hashes the ConfirmTaskPeerInfo object with a hasher
//func (c *ConfirmTaskPeerInfo) HashTreeRootWith(hh *ssz.Hasher) (err error) {
//	indx := hh.Index()
//
//	// Field (0) 'OwnerPeerInfo'
//	if err = c.OwnerPeerInfo.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (1) 'DataSupplierPeerInfoList'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(c.DataSupplierPeerInfoList))
//		if num > 16777216 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for i := uint64(0); i < num; i++ {
//			if err = c.DataSupplierPeerInfoList[i].HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		hh.MerkleizeWithMixin(subIndx, num, 16777216)
//	}
//
//	// Field (2) 'PowerSupplierPeerInfoList'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(c.PowerSupplierPeerInfoList))
//		if num > 16777216 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for i := uint64(0); i < num; i++ {
//			if err = c.PowerSupplierPeerInfoList[i].HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		hh.MerkleizeWithMixin(subIndx, num, 16777216)
//	}
//
//	// Field (3) 'ResultReceiverPeerInfoList'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(c.ResultReceiverPeerInfoList))
//		if num > 16777216 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for i := uint64(0); i < num; i++ {
//			if err = c.ResultReceiverPeerInfoList[i].HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		hh.MerkleizeWithMixin(subIndx, num, 16777216)
//	}
//
//	hh.Merkleize(indx)
//	return
//}
//
//// MarshalSSZ ssz marshals the ConfirmVote object
//func (c *ConfirmVote) MarshalSSZ() ([]byte, error) {
//	return ssz.MarshalSSZ(c)
//}
//
//// MarshalSSZTo ssz marshals the ConfirmVote object to a target array
//func (c *ConfirmVote) MarshalSSZTo(buf []byte) (dst []byte, err error) {
//	dst = buf
//	offset := int(20)
//
//	// Offset (0) 'MsgOption'
//	dst = ssz.WriteOffset(dst, offset)
//	if c.MsgOption == nil {
//		c.MsgOption = new(common.MsgOption)
//	}
//	offset += c.MsgOption.SizeSSZ()
//
//	// Offset (1) 'VoteOption'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(c.VoteOption)
//
//	// Field (2) 'CreateAt'
//	dst = ssz.MarshalUint64(dst, c.CreateAt)
//
//	// Offset (3) 'Sign'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(c.Sign)
//
//	// Field (0) 'MsgOption'
//	if dst, err = c.MsgOption.MarshalSSZTo(dst); err != nil {
//		return
//	}
//
//	// Field (1) 'VoteOption'
//	if len(c.VoteOption) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, c.VoteOption...)
//
//	// Field (3) 'Sign'
//	if len(c.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, c.Sign...)
//
//	return
//}
//
//// UnmarshalSSZ ssz unmarshals the ConfirmVote object
//func (c *ConfirmVote) UnmarshalSSZ(buf []byte) error {
//	var err error
//	size := uint64(len(buf))
//	if size < 20 {
//		return ssz.ErrSize
//	}
//
//	tail := buf
//	var o0, o1, o3 uint64
//
//	// Offset (0) 'MsgOption'
//	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
//		return ssz.ErrOffset
//	}
//
//	if o0 < 20 {
//		return ssz.ErrInvalidVariableOffset
//	}
//
//	// Offset (1) 'VoteOption'
//	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
//		return ssz.ErrOffset
//	}
//
//	// Field (2) 'CreateAt'
//	c.CreateAt = ssz.UnmarshallUint64(buf[8:16])
//
//	// Offset (3) 'Sign'
//	if o3 = ssz.ReadOffset(buf[16:20]); o3 > size || o1 > o3 {
//		return ssz.ErrOffset
//	}
//
//	// Field (0) 'MsgOption'
//	{
//		buf = tail[o0:o1]
//		if c.MsgOption == nil {
//			c.MsgOption = new(common.MsgOption)
//		}
//		if err = c.MsgOption.UnmarshalSSZ(buf); err != nil {
//			return err
//		}
//	}
//
//	// Field (1) 'VoteOption'
//	{
//		buf = tail[o1:o3]
//		if len(buf) > 32 {
//			return ssz.ErrBytesLength
//		}
//		if cap(c.VoteOption) == 0 {
//			c.VoteOption = make([]byte, 0, len(buf))
//		}
//		c.VoteOption = append(c.VoteOption, buf...)
//	}
//
//	// Field (3) 'Sign'
//	{
//		buf = tail[o3:]
//		if len(buf) > 1024 {
//			return ssz.ErrBytesLength
//		}
//		if cap(c.Sign) == 0 {
//			c.Sign = make([]byte, 0, len(buf))
//		}
//		c.Sign = append(c.Sign, buf...)
//	}
//	return err
//}
//
//// SizeSSZ returns the ssz encoded size in bytes for the ConfirmVote object
//func (c *ConfirmVote) SizeSSZ() (size int) {
//	size = 20
//
//	// Field (0) 'MsgOption'
//	if c.MsgOption == nil {
//		c.MsgOption = new(common.MsgOption)
//	}
//	size += c.MsgOption.SizeSSZ()
//
//	// Field (1) 'VoteOption'
//	size += len(c.VoteOption)
//
//	// Field (3) 'Sign'
//	size += len(c.Sign)
//
//	return
//}
//
//// HashTreeRoot ssz hashes the ConfirmVote object
//func (c *ConfirmVote) HashTreeRoot() ([32]byte, error) {
//	return ssz.HashWithDefaultHasher(c)
//}
//
//// HashTreeRootWith ssz hashes the ConfirmVote object with a hasher
//func (c *ConfirmVote) HashTreeRootWith(hh *ssz.Hasher) (err error) {
//	indx := hh.Index()
//
//	// Field (0) 'MsgOption'
//	if err = c.MsgOption.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (1) 'VoteOption'
//	if len(c.VoteOption) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(c.VoteOption)
//
//	// Field (2) 'CreateAt'
//	hh.PutUint64(c.CreateAt)
//
//	// Field (3) 'Sign'
//	if len(c.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(c.Sign)
//
//	hh.Merkleize(indx)
//	return
//}
//
//// MarshalSSZ ssz marshals the CommitMsg object
//func (c *CommitMsg) MarshalSSZ() ([]byte, error) {
//	return ssz.MarshalSSZ(c)
//}
//
//// MarshalSSZTo ssz marshals the CommitMsg object to a target array
//func (c *CommitMsg) MarshalSSZTo(buf []byte) (dst []byte, err error) {
//	dst = buf
//	offset := int(20)
//
//	// Offset (0) 'MsgOption'
//	dst = ssz.WriteOffset(dst, offset)
//	if c.MsgOption == nil {
//		c.MsgOption = new(common.MsgOption)
//	}
//	offset += c.MsgOption.SizeSSZ()
//
//	// Offset (1) 'CommitOption'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(c.CommitOption)
//
//	// Field (2) 'CreateAt'
//	dst = ssz.MarshalUint64(dst, c.CreateAt)
//
//	// Offset (3) 'Sign'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(c.Sign)
//
//	// Field (0) 'MsgOption'
//	if dst, err = c.MsgOption.MarshalSSZTo(dst); err != nil {
//		return
//	}
//
//	// Field (1) 'CommitOption'
//	if len(c.CommitOption) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, c.CommitOption...)
//
//	// Field (3) 'Sign'
//	if len(c.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, c.Sign...)
//
//	return
//}
//
//// UnmarshalSSZ ssz unmarshals the CommitMsg object
//func (c *CommitMsg) UnmarshalSSZ(buf []byte) error {
//	var err error
//	size := uint64(len(buf))
//	if size < 20 {
//		return ssz.ErrSize
//	}
//
//	tail := buf
//	var o0, o1, o3 uint64
//
//	// Offset (0) 'MsgOption'
//	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
//		return ssz.ErrOffset
//	}
//
//	if o0 < 20 {
//		return ssz.ErrInvalidVariableOffset
//	}
//
//	// Offset (1) 'CommitOption'
//	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
//		return ssz.ErrOffset
//	}
//
//	// Field (2) 'CreateAt'
//	c.CreateAt = ssz.UnmarshallUint64(buf[8:16])
//
//	// Offset (3) 'Sign'
//	if o3 = ssz.ReadOffset(buf[16:20]); o3 > size || o1 > o3 {
//		return ssz.ErrOffset
//	}
//
//	// Field (0) 'MsgOption'
//	{
//		buf = tail[o0:o1]
//		if c.MsgOption == nil {
//			c.MsgOption = new(common.MsgOption)
//		}
//		if err = c.MsgOption.UnmarshalSSZ(buf); err != nil {
//			return err
//		}
//	}
//
//	// Field (1) 'CommitOption'
//	{
//		buf = tail[o1:o3]
//		if len(buf) > 32 {
//			return ssz.ErrBytesLength
//		}
//		if cap(c.CommitOption) == 0 {
//			c.CommitOption = make([]byte, 0, len(buf))
//		}
//		c.CommitOption = append(c.CommitOption, buf...)
//	}
//
//	// Field (3) 'Sign'
//	{
//		buf = tail[o3:]
//		if len(buf) > 1024 {
//			return ssz.ErrBytesLength
//		}
//		if cap(c.Sign) == 0 {
//			c.Sign = make([]byte, 0, len(buf))
//		}
//		c.Sign = append(c.Sign, buf...)
//	}
//	return err
//}
//
//// SizeSSZ returns the ssz encoded size in bytes for the CommitMsg object
//func (c *CommitMsg) SizeSSZ() (size int) {
//	size = 20
//
//	// Field (0) 'MsgOption'
//	if c.MsgOption == nil {
//		c.MsgOption = new(common.MsgOption)
//	}
//	size += c.MsgOption.SizeSSZ()
//
//	// Field (1) 'CommitOption'
//	size += len(c.CommitOption)
//
//	// Field (3) 'Sign'
//	size += len(c.Sign)
//
//	return
//}
//
//// HashTreeRoot ssz hashes the CommitMsg object
//func (c *CommitMsg) HashTreeRoot() ([32]byte, error) {
//	return ssz.HashWithDefaultHasher(c)
//}
//
//// HashTreeRootWith ssz hashes the CommitMsg object with a hasher
//func (c *CommitMsg) HashTreeRootWith(hh *ssz.Hasher) (err error) {
//	indx := hh.Index()
//
//	// Field (0) 'MsgOption'
//	if err = c.MsgOption.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (1) 'CommitOption'
//	if len(c.CommitOption) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(c.CommitOption)
//
//	// Field (2) 'CreateAt'
//	hh.PutUint64(c.CreateAt)
//
//	// Field (3) 'Sign'
//	if len(c.Sign) > 1024 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(c.Sign)
//
//	hh.Merkleize(indx)
//	return
//}
//
//// MarshalSSZ ssz marshals the TaskPeerInfo object
//func (t *TaskPeerInfo) MarshalSSZ() ([]byte, error) {
//	return ssz.MarshalSSZ(t)
//}
//
//// MarshalSSZTo ssz marshals the TaskPeerInfo object to a target array
//func (t *TaskPeerInfo) MarshalSSZTo(buf []byte) (dst []byte, err error) {
//	dst = buf
//	offset := int(12)
//
//	// Offset (0) 'Ip'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(t.Ip)
//
//	// Offset (1) 'Port'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(t.Port)
//
//	// Offset (2) 'PartyId'
//	dst = ssz.WriteOffset(dst, offset)
//	offset += len(t.PartyId)
//
//	// Field (0) 'Ip'
//	if len(t.Ip) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, t.Ip...)
//
//	// Field (1) 'Port'
//	if len(t.Port) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, t.Port...)
//
//	// Field (2) 'PartyId'
//	if len(t.PartyId) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	dst = append(dst, t.PartyId...)
//
//	return
//}
//
//// UnmarshalSSZ ssz unmarshals the TaskPeerInfo object
//func (t *TaskPeerInfo) UnmarshalSSZ(buf []byte) error {
//	var err error
//	size := uint64(len(buf))
//	if size < 12 {
//		return ssz.ErrSize
//	}
//
//	tail := buf
//	var o0, o1, o2 uint64
//
//	// Offset (0) 'Ip'
//	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
//		return ssz.ErrOffset
//	}
//
//	if o0 < 12 {
//		return ssz.ErrInvalidVariableOffset
//	}
//
//	// Offset (1) 'Port'
//	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
//		return ssz.ErrOffset
//	}
//
//	// Offset (2) 'PartyId'
//	if o2 = ssz.ReadOffset(buf[8:12]); o2 > size || o1 > o2 {
//		return ssz.ErrOffset
//	}
//
//	// Field (0) 'Ip'
//	{
//		buf = tail[o0:o1]
//		if len(buf) > 32 {
//			return ssz.ErrBytesLength
//		}
//		if cap(t.Ip) == 0 {
//			t.Ip = make([]byte, 0, len(buf))
//		}
//		t.Ip = append(t.Ip, buf...)
//	}
//
//	// Field (1) 'Port'
//	{
//		buf = tail[o1:o2]
//		if len(buf) > 32 {
//			return ssz.ErrBytesLength
//		}
//		if cap(t.Port) == 0 {
//			t.Port = make([]byte, 0, len(buf))
//		}
//		t.Port = append(t.Port, buf...)
//	}
//
//	// Field (2) 'PartyId'
//	{
//		buf = tail[o2:]
//		if len(buf) > 32 {
//			return ssz.ErrBytesLength
//		}
//		if cap(t.PartyId) == 0 {
//			t.PartyId = make([]byte, 0, len(buf))
//		}
//		t.PartyId = append(t.PartyId, buf...)
//	}
//	return err
//}
//
//// SizeSSZ returns the ssz encoded size in bytes for the TaskPeerInfo object
//func (t *TaskPeerInfo) SizeSSZ() (size int) {
//	size = 12
//
//	// Field (0) 'Ip'
//	size += len(t.Ip)
//
//	// Field (1) 'Port'
//	size += len(t.Port)
//
//	// Field (2) 'PartyId'
//	size += len(t.PartyId)
//
//	return
//}
//
//// HashTreeRoot ssz hashes the TaskPeerInfo object
//func (t *TaskPeerInfo) HashTreeRoot() ([32]byte, error) {
//	return ssz.HashWithDefaultHasher(t)
//}
//
//// HashTreeRootWith ssz hashes the TaskPeerInfo object with a hasher
//func (t *TaskPeerInfo) HashTreeRootWith(hh *ssz.Hasher) (err error) {
//	indx := hh.Index()
//
//	// Field (0) 'Ip'
//	if len(t.Ip) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(t.Ip)
//
//	// Field (1) 'Port'
//	if len(t.Port) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(t.Port)
//
//	// Field (2) 'PartyId'
//	if len(t.PartyId) > 32 {
//		err = ssz.ErrBytesLength
//		return
//	}
//	hh.PutBytes(t.PartyId)
//
//	hh.Merkleize(indx)
//	return
//}
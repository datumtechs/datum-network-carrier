// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lib/p2p/v1/messages.proto

package v1

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_prysmaticlabs_eth2_types "github.com/prysmaticlabs/eth2-types"
	github_com_prysmaticlabs_go_bitfield "github.com/prysmaticlabs/go-bitfield"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Status struct {
	ForkDigest           []byte                                    `protobuf:"bytes,1,opt,name=fork_digest,json=forkDigest,proto3" json:"fork_digest,omitempty" ssz-size:"4"`
	FinalizedRoot        []byte                                    `protobuf:"bytes,2,opt,name=finalized_root,json=finalizedRoot,proto3" json:"finalized_root,omitempty" ssz-size:"32"`
	FinalizedEpoch       github_com_prysmaticlabs_eth2_types.Epoch `protobuf:"varint,3,opt,name=finalized_epoch,json=finalizedEpoch,proto3,casttype=github.com/prysmaticlabs/eth2-types.Epoch" json:"finalized_epoch,omitempty"`
	HeadRoot             []byte                                    `protobuf:"bytes,4,opt,name=head_root,json=headRoot,proto3" json:"head_root,omitempty" ssz-size:"32"`
	HeadSlot             github_com_prysmaticlabs_eth2_types.Slot  `protobuf:"varint,5,opt,name=head_slot,json=headSlot,proto3,casttype=github.com/prysmaticlabs/eth2-types.Slot" json:"head_slot,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                  `json:"-"`
	XXX_unrecognized     []byte                                    `json:"-"`
	XXX_sizecache        int32                                     `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c5bd9e5eeffff24, []int{0}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Status.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return m.Size()
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetForkDigest() []byte {
	if m != nil {
		return m.ForkDigest
	}
	return nil
}

func (m *Status) GetFinalizedRoot() []byte {
	if m != nil {
		return m.FinalizedRoot
	}
	return nil
}

func (m *Status) GetFinalizedEpoch() github_com_prysmaticlabs_eth2_types.Epoch {
	if m != nil {
		return m.FinalizedEpoch
	}
	return 0
}

func (m *Status) GetHeadRoot() []byte {
	if m != nil {
		return m.HeadRoot
	}
	return nil
}

func (m *Status) GetHeadSlot() github_com_prysmaticlabs_eth2_types.Slot {
	if m != nil {
		return m.HeadSlot
	}
	return 0
}

type CarrierBlocksByRangeRequest struct {
	StartSlot            github_com_prysmaticlabs_eth2_types.Slot `protobuf:"varint,1,opt,name=start_slot,json=startSlot,proto3,casttype=github.com/prysmaticlabs/eth2-types.Slot" json:"start_slot,omitempty"`
	Count                uint64                                   `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	Step                 uint64                                   `protobuf:"varint,3,opt,name=step,proto3" json:"step,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                 `json:"-"`
	XXX_unrecognized     []byte                                   `json:"-"`
	XXX_sizecache        int32                                    `json:"-"`
}

func (m *CarrierBlocksByRangeRequest) Reset()         { *m = CarrierBlocksByRangeRequest{} }
func (m *CarrierBlocksByRangeRequest) String() string { return proto.CompactTextString(m) }
func (*CarrierBlocksByRangeRequest) ProtoMessage()    {}
func (*CarrierBlocksByRangeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c5bd9e5eeffff24, []int{1}
}
func (m *CarrierBlocksByRangeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CarrierBlocksByRangeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CarrierBlocksByRangeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CarrierBlocksByRangeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CarrierBlocksByRangeRequest.Merge(m, src)
}
func (m *CarrierBlocksByRangeRequest) XXX_Size() int {
	return m.Size()
}
func (m *CarrierBlocksByRangeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CarrierBlocksByRangeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CarrierBlocksByRangeRequest proto.InternalMessageInfo

func (m *CarrierBlocksByRangeRequest) GetStartSlot() github_com_prysmaticlabs_eth2_types.Slot {
	if m != nil {
		return m.StartSlot
	}
	return 0
}

func (m *CarrierBlocksByRangeRequest) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *CarrierBlocksByRangeRequest) GetStep() uint64 {
	if m != nil {
		return m.Step
	}
	return 0
}

type ENRForkID struct {
	CurrentForkDigest    []byte                                    `protobuf:"bytes,1,opt,name=current_fork_digest,json=currentForkDigest,proto3" json:"current_fork_digest,omitempty" ssz-size:"4"`
	NextForkVersion      []byte                                    `protobuf:"bytes,2,opt,name=next_fork_version,json=nextForkVersion,proto3" json:"next_fork_version,omitempty" ssz-size:"4"`
	NextForkEpoch        github_com_prysmaticlabs_eth2_types.Epoch `protobuf:"varint,3,opt,name=next_fork_epoch,json=nextForkEpoch,proto3,casttype=github.com/prysmaticlabs/eth2-types.Epoch" json:"next_fork_epoch,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                  `json:"-"`
	XXX_unrecognized     []byte                                    `json:"-"`
	XXX_sizecache        int32                                     `json:"-"`
}

func (m *ENRForkID) Reset()         { *m = ENRForkID{} }
func (m *ENRForkID) String() string { return proto.CompactTextString(m) }
func (*ENRForkID) ProtoMessage()    {}
func (*ENRForkID) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c5bd9e5eeffff24, []int{2}
}
func (m *ENRForkID) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ENRForkID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ENRForkID.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ENRForkID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ENRForkID.Merge(m, src)
}
func (m *ENRForkID) XXX_Size() int {
	return m.Size()
}
func (m *ENRForkID) XXX_DiscardUnknown() {
	xxx_messageInfo_ENRForkID.DiscardUnknown(m)
}

var xxx_messageInfo_ENRForkID proto.InternalMessageInfo

func (m *ENRForkID) GetCurrentForkDigest() []byte {
	if m != nil {
		return m.CurrentForkDigest
	}
	return nil
}

func (m *ENRForkID) GetNextForkVersion() []byte {
	if m != nil {
		return m.NextForkVersion
	}
	return nil
}

func (m *ENRForkID) GetNextForkEpoch() github_com_prysmaticlabs_eth2_types.Epoch {
	if m != nil {
		return m.NextForkEpoch
	}
	return 0
}

//
//Spec Definition:
//MetaData
//(
//seq_number: uint64
//attnets: Bitvector[ATTESTATION_SUBNET_COUNT]
//)
type MetaData struct {
	SeqNumber            uint64                                           `protobuf:"varint,1,opt,name=seq_number,json=seqNumber,proto3" json:"seq_number,omitempty"`
	Attnets              github_com_prysmaticlabs_go_bitfield.Bitvector64 `protobuf:"bytes,2,opt,name=attnets,proto3,casttype=github.com/prysmaticlabs/go-bitfield.Bitvector64" json:"attnets,omitempty" ssz-size:"8"`
	XXX_NoUnkeyedLiteral struct{}                                         `json:"-"`
	XXX_unrecognized     []byte                                           `json:"-"`
	XXX_sizecache        int32                                            `json:"-"`
}

func (m *MetaData) Reset()         { *m = MetaData{} }
func (m *MetaData) String() string { return proto.CompactTextString(m) }
func (*MetaData) ProtoMessage()    {}
func (*MetaData) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c5bd9e5eeffff24, []int{3}
}
func (m *MetaData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MetaData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MetaData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MetaData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetaData.Merge(m, src)
}
func (m *MetaData) XXX_Size() int {
	return m.Size()
}
func (m *MetaData) XXX_DiscardUnknown() {
	xxx_messageInfo_MetaData.DiscardUnknown(m)
}

var xxx_messageInfo_MetaData proto.InternalMessageInfo

func (m *MetaData) GetSeqNumber() uint64 {
	if m != nil {
		return m.SeqNumber
	}
	return 0
}

func (m *MetaData) GetAttnets() github_com_prysmaticlabs_go_bitfield.Bitvector64 {
	if m != nil {
		return m.Attnets
	}
	return nil
}

func init() {
	proto.RegisterType((*Status)(nil), "carrier.p2p.v1.Status")
	proto.RegisterType((*CarrierBlocksByRangeRequest)(nil), "carrier.p2p.v1.CarrierBlocksByRangeRequest")
	proto.RegisterType((*ENRForkID)(nil), "carrier.p2p.v1.ENRForkID")
	proto.RegisterType((*MetaData)(nil), "carrier.p2p.v1.MetaData")
}

func init() { proto.RegisterFile("lib/p2p/v1/messages.proto", fileDescriptor_5c5bd9e5eeffff24) }

var fileDescriptor_5c5bd9e5eeffff24 = []byte{
	// 542 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xe5, 0x92, 0x96, 0x66, 0xe9, 0x1f, 0xba, 0x70, 0x08, 0x20, 0x9a, 0xca, 0x5c, 0x8a,
	0xd4, 0xd8, 0x34, 0xad, 0x50, 0x55, 0x71, 0x40, 0x26, 0x0d, 0xaa, 0x10, 0x3d, 0x6c, 0x45, 0x0f,
	0x1c, 0x88, 0xd6, 0xce, 0xc4, 0x59, 0xc5, 0xf1, 0x3a, 0xbb, 0xe3, 0x40, 0xf2, 0x06, 0xdc, 0x79,
	0xa8, 0x1e, 0x79, 0x82, 0x08, 0xe5, 0x11, 0x7a, 0xcc, 0x09, 0x79, 0xed, 0x62, 0x10, 0xaa, 0x54,
	0xc1, 0x6d, 0x26, 0xf9, 0x7e, 0x33, 0x9f, 0xe6, 0xf3, 0x92, 0x47, 0x91, 0xf0, 0xdd, 0xa4, 0x99,
	0xb8, 0xe3, 0x7d, 0x77, 0x08, 0x5a, 0xf3, 0x10, 0xb4, 0x93, 0x28, 0x89, 0x92, 0x6e, 0x04, 0x5c,
	0x29, 0x01, 0xca, 0x49, 0x9a, 0x89, 0x33, 0xde, 0x7f, 0xfc, 0x4c, 0x41, 0x22, 0xb5, 0x6b, 0xfe,
	0xf4, 0xd3, 0x9e, 0x1b, 0xca, 0x50, 0x9a, 0xc6, 0x54, 0x39, 0x64, 0x5f, 0x2e, 0x91, 0x95, 0x73,
	0xe4, 0x98, 0x6a, 0xba, 0x4f, 0xee, 0xf5, 0xa4, 0x1a, 0x74, 0xba, 0x22, 0x04, 0x8d, 0x35, 0x6b,
	0xc7, 0xda, 0x5d, 0xf3, 0xee, 0x5f, 0xcd, 0xea, 0x6b, 0x5a, 0x4f, 0x1b, 0x5a, 0x4c, 0xe1, 0xd8,
	0x3e, 0xb4, 0x19, 0xc9, 0x44, 0x2d, 0xa3, 0xa1, 0x47, 0x64, 0xa3, 0x27, 0x62, 0x1e, 0x89, 0x29,
	0x74, 0x3b, 0x4a, 0x4a, 0xac, 0x2d, 0x19, 0x6a, 0xeb, 0x6a, 0x56, 0x5f, 0x2f, 0xa9, 0x83, 0xa6,
	0xcd, 0xd6, 0x7f, 0x09, 0x99, 0x94, 0x48, 0x2f, 0xc8, 0x66, 0x49, 0x42, 0x22, 0x83, 0x7e, 0xed,
	0xce, 0x8e, 0xb5, 0x5b, 0xf1, 0x1a, 0x8b, 0x59, 0xfd, 0x79, 0x28, 0xb0, 0x9f, 0xfa, 0x4e, 0x20,
	0x87, 0x6e, 0xa2, 0x26, 0x7a, 0xc8, 0x51, 0x04, 0x11, 0xf7, 0xb5, 0x0b, 0xd8, 0x6f, 0x36, 0x70,
	0x92, 0x80, 0x76, 0x4e, 0x32, 0x88, 0x95, 0xfb, 0x4d, 0x4f, 0x1d, 0x52, 0xed, 0x03, 0x2f, 0xcc,
	0x54, 0x6e, 0x32, 0xb3, 0x9a, 0x69, 0x8c, 0x8f, 0xd3, 0x42, 0xaf, 0x23, 0x89, 0xb5, 0x65, 0xe3,
	0x60, 0x6f, 0x31, 0xab, 0xef, 0xde, 0xc6, 0xc1, 0x79, 0x24, 0x31, 0x1f, 0x95, 0x55, 0xf6, 0x37,
	0x8b, 0x3c, 0x79, 0x93, 0x47, 0xe0, 0x45, 0x32, 0x18, 0x68, 0x6f, 0xc2, 0x78, 0x1c, 0x02, 0x83,
	0x51, 0x9a, 0x1d, 0xeb, 0x1d, 0x21, 0x1a, 0xb9, 0xc2, 0x7c, 0x97, 0xf5, 0x0f, 0xbb, 0xaa, 0x86,
	0xcf, 0x4a, 0xfa, 0x90, 0x2c, 0x07, 0x32, 0x8d, 0xf3, 0x83, 0x57, 0x58, 0xde, 0x50, 0x4a, 0x2a,
	0x1a, 0x21, 0xc9, 0x4f, 0xc9, 0x4c, 0x6d, 0xcf, 0x2d, 0x52, 0x3d, 0x39, 0x63, 0x6d, 0xa9, 0x06,
	0xa7, 0x2d, 0xfa, 0x9a, 0x3c, 0x08, 0x52, 0xa5, 0x20, 0xc6, 0xce, 0x6d, 0xc2, 0xde, 0x2a, 0xc4,
	0xed, 0x32, 0xf3, 0x57, 0x64, 0x2b, 0x86, 0x2f, 0x05, 0x3e, 0x06, 0xa5, 0x85, 0x8c, 0x8b, 0xd8,
	0xff, 0xe6, 0x37, 0x33, 0x69, 0x06, 0x5f, 0xe4, 0x42, 0xfa, 0x81, 0x6c, 0x96, 0xf4, 0x7f, 0xe4,
	0xbe, 0x7e, 0x3d, 0xd8, 0xb4, 0xf6, 0x57, 0x8b, 0xac, 0xbe, 0x07, 0xe4, 0x2d, 0x8e, 0x9c, 0x3e,
	0x25, 0x44, 0xc3, 0xa8, 0x13, 0xa7, 0x43, 0x1f, 0x54, 0x7e, 0x68, 0x56, 0xd5, 0x30, 0x3a, 0x33,
	0x3f, 0xd0, 0x4f, 0xe4, 0x2e, 0x47, 0x8c, 0x01, 0x75, 0x61, 0xbb, 0xf5, 0xa7, 0xed, 0x23, 0x7b,
	0x31, 0xab, 0xbf, 0xb8, 0xd1, 0x4a, 0x28, 0x1b, 0xbe, 0xc0, 0x9e, 0x80, 0xa8, 0xeb, 0x78, 0x02,
	0xc7, 0x10, 0xa0, 0x54, 0x2f, 0x0f, 0xd9, 0xf5, 0x50, 0xef, 0xf8, 0x72, 0xbe, 0x6d, 0x7d, 0x9f,
	0x6f, 0x5b, 0x3f, 0xe6, 0xdb, 0xd6, 0xc7, 0xbd, 0xdf, 0x06, 0x31, 0xa9, 0x01, 0x91, 0xb7, 0x23,
	0xf9, 0xd9, 0x2d, 0xbe, 0x94, 0xc6, 0x5b, 0xe9, 0x96, 0x4f, 0xda, 0x5f, 0x31, 0xaf, 0xf2, 0xe0,
	0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0x37, 0x4e, 0x34, 0x63, 0xe7, 0x03, 0x00, 0x00,
}

func (m *Status) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.HeadSlot != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.HeadSlot))
		i--
		dAtA[i] = 0x28
	}
	if len(m.HeadRoot) > 0 {
		i -= len(m.HeadRoot)
		copy(dAtA[i:], m.HeadRoot)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.HeadRoot)))
		i--
		dAtA[i] = 0x22
	}
	if m.FinalizedEpoch != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.FinalizedEpoch))
		i--
		dAtA[i] = 0x18
	}
	if len(m.FinalizedRoot) > 0 {
		i -= len(m.FinalizedRoot)
		copy(dAtA[i:], m.FinalizedRoot)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.FinalizedRoot)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ForkDigest) > 0 {
		i -= len(m.ForkDigest)
		copy(dAtA[i:], m.ForkDigest)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.ForkDigest)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *CarrierBlocksByRangeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CarrierBlocksByRangeRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CarrierBlocksByRangeRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Step != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.Step))
		i--
		dAtA[i] = 0x18
	}
	if m.Count != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.Count))
		i--
		dAtA[i] = 0x10
	}
	if m.StartSlot != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.StartSlot))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ENRForkID) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ENRForkID) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ENRForkID) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.NextForkEpoch != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.NextForkEpoch))
		i--
		dAtA[i] = 0x18
	}
	if len(m.NextForkVersion) > 0 {
		i -= len(m.NextForkVersion)
		copy(dAtA[i:], m.NextForkVersion)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.NextForkVersion)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.CurrentForkDigest) > 0 {
		i -= len(m.CurrentForkDigest)
		copy(dAtA[i:], m.CurrentForkDigest)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.CurrentForkDigest)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MetaData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetaData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetaData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Attnets) > 0 {
		i -= len(m.Attnets)
		copy(dAtA[i:], m.Attnets)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.Attnets)))
		i--
		dAtA[i] = 0x12
	}
	if m.SeqNumber != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.SeqNumber))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintMessages(dAtA []byte, offset int, v uint64) int {
	offset -= sovMessages(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ForkDigest)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	l = len(m.FinalizedRoot)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.FinalizedEpoch != 0 {
		n += 1 + sovMessages(uint64(m.FinalizedEpoch))
	}
	l = len(m.HeadRoot)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.HeadSlot != 0 {
		n += 1 + sovMessages(uint64(m.HeadSlot))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *CarrierBlocksByRangeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.StartSlot != 0 {
		n += 1 + sovMessages(uint64(m.StartSlot))
	}
	if m.Count != 0 {
		n += 1 + sovMessages(uint64(m.Count))
	}
	if m.Step != 0 {
		n += 1 + sovMessages(uint64(m.Step))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ENRForkID) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CurrentForkDigest)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	l = len(m.NextForkVersion)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.NextForkEpoch != 0 {
		n += 1 + sovMessages(uint64(m.NextForkEpoch))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *MetaData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SeqNumber != 0 {
		n += 1 + sovMessages(uint64(m.SeqNumber))
	}
	l = len(m.Attnets)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMessages(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMessages(x uint64) (n int) {
	return sovMessages(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Status) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Status: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Status: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ForkDigest", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ForkDigest = append(m.ForkDigest[:0], dAtA[iNdEx:postIndex]...)
			if m.ForkDigest == nil {
				m.ForkDigest = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FinalizedRoot", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FinalizedRoot = append(m.FinalizedRoot[:0], dAtA[iNdEx:postIndex]...)
			if m.FinalizedRoot == nil {
				m.FinalizedRoot = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FinalizedEpoch", wireType)
			}
			m.FinalizedEpoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FinalizedEpoch |= github_com_prysmaticlabs_eth2_types.Epoch(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HeadRoot", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HeadRoot = append(m.HeadRoot[:0], dAtA[iNdEx:postIndex]...)
			if m.HeadRoot == nil {
				m.HeadRoot = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HeadSlot", wireType)
			}
			m.HeadSlot = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HeadSlot |= github_com_prysmaticlabs_eth2_types.Slot(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CarrierBlocksByRangeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CarrierBlocksByRangeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CarrierBlocksByRangeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartSlot", wireType)
			}
			m.StartSlot = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartSlot |= github_com_prysmaticlabs_eth2_types.Slot(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Count |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Step", wireType)
			}
			m.Step = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Step |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ENRForkID) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ENRForkID: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ENRForkID: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentForkDigest", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CurrentForkDigest = append(m.CurrentForkDigest[:0], dAtA[iNdEx:postIndex]...)
			if m.CurrentForkDigest == nil {
				m.CurrentForkDigest = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextForkVersion", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NextForkVersion = append(m.NextForkVersion[:0], dAtA[iNdEx:postIndex]...)
			if m.NextForkVersion == nil {
				m.NextForkVersion = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextForkEpoch", wireType)
			}
			m.NextForkEpoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextForkEpoch |= github_com_prysmaticlabs_eth2_types.Epoch(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MetaData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetaData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetaData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SeqNumber", wireType)
			}
			m.SeqNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SeqNumber |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attnets", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Attnets = append(m.Attnets[:0], dAtA[iNdEx:postIndex]...)
			if m.Attnets == nil {
				m.Attnets = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMessages(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMessages
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMessages
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMessages
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMessages        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessages          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMessages = fmt.Errorf("proto: unexpected end of group")
)
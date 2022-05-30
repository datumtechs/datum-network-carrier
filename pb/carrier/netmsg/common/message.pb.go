// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: carrier/netmsg/common/message.proto

package common

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// 消息选项
type MsgOption struct {
	ProposalId           []byte                        `protobuf:"bytes,1,opt,name=proposal_id,json=proposalId,proto3" json:"proposal_id,omitempty" ssz-max:"1024"`
	SenderRole           uint64                        `protobuf:"varint,2,opt,name=sender_role,json=senderRole,proto3" json:"sender_role,omitempty" ssz-size:"8"`
	SenderPartyId        []byte                        `protobuf:"bytes,3,opt,name=sender_party_id,json=senderPartyId,proto3" json:"sender_party_id,omitempty" ssz-max:"32"`
	ReceiverRole         uint64                        `protobuf:"varint,4,opt,name=receiver_role,json=receiverRole,proto3" json:"receiver_role,omitempty" ssz-size:"8"`
	ReceiverPartyId      []byte                        `protobuf:"bytes,5,opt,name=receiver_party_id,json=receiverPartyId,proto3" json:"receiver_party_id,omitempty" ssz-max:"32"`
	MsgOwner             *TaskOrganizationIdentityInfo `protobuf:"bytes,6,opt,name=msg_owner,json=msgOwner,proto3" json:"msg_owner,omitempty" ssz-max:"16777216"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *MsgOption) Reset()         { *m = MsgOption{} }
func (m *MsgOption) String() string { return proto.CompactTextString(m) }
func (*MsgOption) ProtoMessage()    {}
func (*MsgOption) Descriptor() ([]byte, []int) {
	return fileDescriptor_17b6b368a7264c14, []int{0}
}
func (m *MsgOption) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgOption) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgOption.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgOption) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgOption.Merge(m, src)
}
func (m *MsgOption) XXX_Size() int {
	return m.Size()
}
func (m *MsgOption) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgOption.DiscardUnknown(m)
}

var xxx_messageInfo_MsgOption proto.InternalMessageInfo

func (m *MsgOption) GetProposalId() []byte {
	if m != nil {
		return m.ProposalId
	}
	return nil
}

func (m *MsgOption) GetSenderRole() uint64 {
	if m != nil {
		return m.SenderRole
	}
	return 0
}

func (m *MsgOption) GetSenderPartyId() []byte {
	if m != nil {
		return m.SenderPartyId
	}
	return nil
}

func (m *MsgOption) GetReceiverRole() uint64 {
	if m != nil {
		return m.ReceiverRole
	}
	return 0
}

func (m *MsgOption) GetReceiverPartyId() []byte {
	if m != nil {
		return m.ReceiverPartyId
	}
	return nil
}

func (m *MsgOption) GetMsgOwner() *TaskOrganizationIdentityInfo {
	if m != nil {
		return m.MsgOwner
	}
	return nil
}

// 组织(节点)唯一标识抽象
type TaskOrganizationIdentityInfo struct {
	Name                 []byte   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" ssz-max:"64"`
	NodeId               []byte   `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty" ssz-max:"1024"`
	IdentityId           []byte   `protobuf:"bytes,3,opt,name=identity_id,json=identityId,proto3" json:"identity_id,omitempty" ssz-max:"1024"`
	PartyId              []byte   `protobuf:"bytes,4,opt,name=party_id,json=partyId,proto3" json:"party_id,omitempty" ssz-max:"32"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskOrganizationIdentityInfo) Reset()         { *m = TaskOrganizationIdentityInfo{} }
func (m *TaskOrganizationIdentityInfo) String() string { return proto.CompactTextString(m) }
func (*TaskOrganizationIdentityInfo) ProtoMessage()    {}
func (*TaskOrganizationIdentityInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_17b6b368a7264c14, []int{1}
}
func (m *TaskOrganizationIdentityInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TaskOrganizationIdentityInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TaskOrganizationIdentityInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TaskOrganizationIdentityInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskOrganizationIdentityInfo.Merge(m, src)
}
func (m *TaskOrganizationIdentityInfo) XXX_Size() int {
	return m.Size()
}
func (m *TaskOrganizationIdentityInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskOrganizationIdentityInfo.DiscardUnknown(m)
}

var xxx_messageInfo_TaskOrganizationIdentityInfo proto.InternalMessageInfo

func (m *TaskOrganizationIdentityInfo) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *TaskOrganizationIdentityInfo) GetNodeId() []byte {
	if m != nil {
		return m.NodeId
	}
	return nil
}

func (m *TaskOrganizationIdentityInfo) GetIdentityId() []byte {
	if m != nil {
		return m.IdentityId
	}
	return nil
}

func (m *TaskOrganizationIdentityInfo) GetPartyId() []byte {
	if m != nil {
		return m.PartyId
	}
	return nil
}

// 资源消耗概览 (系统 和 task 的通用)
type ResourceUsage struct {
	// 服务系统的总内存 (单位: byte)
	TotalMem uint64 `protobuf:"varint,1,opt,name=total_mem,json=totalMem,proto3" json:"total_mem,omitempty" ssz-size:"8"`
	// 服务系统的已用内存  (单位: byte)
	UsedMem uint64 `protobuf:"varint,2,opt,name=used_mem,json=usedMem,proto3" json:"used_mem,omitempty" ssz-size:"8"`
	// 服务的总内核数 (单位: 个)
	TotalProcessor uint64 `protobuf:"varint,3,opt,name=total_processor,json=totalProcessor,proto3" json:"total_processor,omitempty" ssz-size:"8"`
	// 服务的已用内核数 (单位: 个)
	UsedProcessor uint64 `protobuf:"varint,4,opt,name=used_processor,json=usedProcessor,proto3" json:"used_processor,omitempty" ssz-size:"8"`
	// 服务的总带宽数 (单位: bps)
	TotalBandwidth uint64 `protobuf:"varint,5,opt,name=total_bandwidth,json=totalBandwidth,proto3" json:"total_bandwidth,omitempty" ssz-size:"8"`
	// 服务的已用带宽数 (单位: bps)
	UsedBandwidth uint64 `protobuf:"varint,6,opt,name=used_bandwidth,json=usedBandwidth,proto3" json:"used_bandwidth,omitempty" ssz-size:"8"`
	// 服务的总磁盘空间 (单位: byte)
	TotalDisk uint64 `protobuf:"varint,7,opt,name=total_disk,json=totalDisk,proto3" json:"total_disk,omitempty" ssz-size:"8"`
	// 服务的已用磁盘空间 (单位: byte)
	UsedDisk             uint64   `protobuf:"varint,8,opt,name=used_disk,json=usedDisk,proto3" json:"used_disk,omitempty" ssz-size:"8"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResourceUsage) Reset()         { *m = ResourceUsage{} }
func (m *ResourceUsage) String() string { return proto.CompactTextString(m) }
func (*ResourceUsage) ProtoMessage()    {}
func (*ResourceUsage) Descriptor() ([]byte, []int) {
	return fileDescriptor_17b6b368a7264c14, []int{2}
}
func (m *ResourceUsage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ResourceUsage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ResourceUsage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ResourceUsage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourceUsage.Merge(m, src)
}
func (m *ResourceUsage) XXX_Size() int {
	return m.Size()
}
func (m *ResourceUsage) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourceUsage.DiscardUnknown(m)
}

var xxx_messageInfo_ResourceUsage proto.InternalMessageInfo

func (m *ResourceUsage) GetTotalMem() uint64 {
	if m != nil {
		return m.TotalMem
	}
	return 0
}

func (m *ResourceUsage) GetUsedMem() uint64 {
	if m != nil {
		return m.UsedMem
	}
	return 0
}

func (m *ResourceUsage) GetTotalProcessor() uint64 {
	if m != nil {
		return m.TotalProcessor
	}
	return 0
}

func (m *ResourceUsage) GetUsedProcessor() uint64 {
	if m != nil {
		return m.UsedProcessor
	}
	return 0
}

func (m *ResourceUsage) GetTotalBandwidth() uint64 {
	if m != nil {
		return m.TotalBandwidth
	}
	return 0
}

func (m *ResourceUsage) GetUsedBandwidth() uint64 {
	if m != nil {
		return m.UsedBandwidth
	}
	return 0
}

func (m *ResourceUsage) GetTotalDisk() uint64 {
	if m != nil {
		return m.TotalDisk
	}
	return 0
}

func (m *ResourceUsage) GetUsedDisk() uint64 {
	if m != nil {
		return m.UsedDisk
	}
	return 0
}

func init() {
	proto.RegisterType((*MsgOption)(nil), "carrier.netmsg.common.MsgOption")
	proto.RegisterType((*TaskOrganizationIdentityInfo)(nil), "carrier.netmsg.common.TaskOrganizationIdentityInfo")
	proto.RegisterType((*ResourceUsage)(nil), "carrier.netmsg.common.ResourceUsage")
}

func init() {
	proto.RegisterFile("carrier/netmsg/common/message.proto", fileDescriptor_17b6b368a7264c14)
}

var fileDescriptor_17b6b368a7264c14 = []byte{
	// 587 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0xc1, 0x4e, 0xdb, 0x3e,
	0x1c, 0xc7, 0x15, 0xe8, 0xbf, 0x2d, 0x86, 0xc2, 0x7f, 0x96, 0x36, 0x55, 0xd3, 0x04, 0x28, 0xec,
	0x80, 0x84, 0xda, 0x8c, 0x96, 0x51, 0x86, 0x76, 0xaa, 0xb8, 0xe4, 0x80, 0x40, 0xd1, 0x76, 0xd9,
	0xa5, 0x72, 0xe3, 0x1f, 0xa9, 0xd5, 0xda, 0x8e, 0x6c, 0x77, 0x0c, 0x9e, 0x61, 0x6f, 0xb3, 0x97,
	0xd8, 0x71, 0xd2, 0xee, 0x68, 0xea, 0x23, 0xf4, 0x09, 0x26, 0x3b, 0x4d, 0xa2, 0x49, 0x0d, 0xb7,
	0x20, 0x3e, 0x5f, 0x7f, 0xea, 0xef, 0xef, 0x27, 0xa3, 0xa3, 0x98, 0x28, 0xc5, 0x40, 0x05, 0x02,
	0x0c, 0xd7, 0x49, 0x10, 0x4b, 0xce, 0xa5, 0x08, 0x38, 0x68, 0x4d, 0x12, 0xe8, 0xa6, 0x4a, 0x1a,
	0x89, 0x5f, 0xae, 0xa0, 0x6e, 0x06, 0x75, 0x33, 0xe8, 0xf5, 0x91, 0x82, 0x54, 0xea, 0xc0, 0x31,
	0xe3, 0xf9, 0x5d, 0x90, 0xc8, 0x44, 0xba, 0x3f, 0xdc, 0x57, 0x96, 0xf5, 0xbf, 0x6f, 0xa2, 0xad,
	0x6b, 0x9d, 0xdc, 0xa4, 0x86, 0x49, 0x81, 0xfb, 0x68, 0x3b, 0x55, 0x32, 0x95, 0x9a, 0xcc, 0x46,
	0x8c, 0xb6, 0xbd, 0x43, 0xef, 0x78, 0x67, 0x88, 0x97, 0x4f, 0x07, 0xbb, 0x5a, 0x3f, 0x76, 0x38,
	0xf9, 0x76, 0xe9, 0x9f, 0xbe, 0xeb, 0x9d, 0xf9, 0x11, 0xca, 0xb1, 0x90, 0xe2, 0x53, 0xb4, 0xad,
	0x41, 0x50, 0x50, 0x23, 0x25, 0x67, 0xd0, 0xde, 0x38, 0xf4, 0x8e, 0x6b, 0xc3, 0xff, 0x97, 0x4f,
	0x07, 0x3b, 0x36, 0xa4, 0xd9, 0x23, 0x5c, 0xfa, 0x17, 0x7e, 0x84, 0x32, 0x28, 0x92, 0x33, 0xc0,
	0x17, 0x68, 0x6f, 0x15, 0x49, 0x89, 0x32, 0x0f, 0xd6, 0xb5, 0xe9, 0x5c, 0x45, 0xcc, 0xb9, 0xfa,
	0x3d, 0x3f, 0x6a, 0x65, 0xe0, 0xad, 0xe5, 0x42, 0x8a, 0xdf, 0xa3, 0x96, 0x82, 0x18, 0xd8, 0xd7,
	0x5c, 0x57, 0xab, 0xd0, 0xed, 0xe4, 0x98, 0x13, 0x7e, 0x44, 0x2f, 0x8a, 0x58, 0xa1, 0xfc, 0xaf,
	0x42, 0xb9, 0x97, 0xa3, 0xb9, 0x74, 0x82, 0xb6, 0xb8, 0x4e, 0x46, 0xf2, 0x5e, 0x80, 0x6a, 0xd7,
	0x0f, 0xbd, 0xe3, 0xed, 0x5e, 0xbf, 0xbb, 0xb6, 0xf4, 0xee, 0x27, 0xa2, 0xa7, 0x37, 0x2a, 0x21,
	0x82, 0x3d, 0x12, 0x5b, 0x69, 0x48, 0x41, 0x18, 0x66, 0x1e, 0x42, 0x71, 0x27, 0x87, 0xaf, 0x96,
	0x4f, 0x07, 0xb8, 0x6c, 0xf2, 0x7c, 0x30, 0x18, 0xf4, 0x4e, 0xcf, 0xfd, 0xa8, 0xc9, 0x75, 0x72,
	0x63, 0x0f, 0xf7, 0x7f, 0x7b, 0xe8, 0xcd, 0x73, 0x47, 0xe0, 0xb7, 0xa8, 0x26, 0x08, 0x87, 0xd5,
	0x68, 0xfe, 0xfd, 0xed, 0xe7, 0x67, 0x7e, 0xe4, 0xfe, 0x8b, 0x4f, 0x50, 0x43, 0x48, 0x0a, 0xf6,
	0x92, 0x1b, 0x95, 0x33, 0xac, 0x5b, 0x24, 0xa4, 0x76, 0xe8, 0x6c, 0xa5, 0x28, 0x07, 0xb1, 0x76,
	0xe8, 0x39, 0x16, 0x52, 0x7c, 0x82, 0x9a, 0x45, 0x8f, 0xb5, 0x8a, 0x1e, 0x1b, 0x69, 0xd6, 0x9f,
	0xff, 0x63, 0x13, 0xb5, 0x22, 0xd0, 0x72, 0xae, 0x62, 0xf8, 0x6c, 0x17, 0x17, 0x77, 0xd0, 0x96,
	0x91, 0x86, 0xcc, 0x46, 0x1c, 0xb8, 0xbb, 0xcb, 0xba, 0x11, 0x36, 0x1d, 0x72, 0x0d, 0xdc, 0xda,
	0xe6, 0x1a, 0xa8, 0xa3, 0xab, 0xf6, 0xab, 0x61, 0x09, 0x0b, 0x7f, 0x40, 0x7b, 0xd9, 0xd9, 0xa9,
	0x92, 0x31, 0x68, 0x2d, 0x95, 0xbb, 0xd3, 0xba, 0xcc, 0xae, 0x03, 0x6f, 0x73, 0x0e, 0x0f, 0xd0,
	0xae, 0xf3, 0x94, 0xc9, 0xaa, 0xf5, 0x6a, 0x59, 0xae, 0x0c, 0x16, 0xce, 0x31, 0x11, 0xf4, 0x9e,
	0x51, 0x33, 0x71, 0xdb, 0x55, 0xed, 0x1c, 0xe6, 0x5c, 0xe1, 0x2c, 0x93, 0xf5, 0xe7, 0x9c, 0x65,
	0x30, 0x40, 0x28, 0x73, 0x52, 0xa6, 0xa7, 0xed, 0x46, 0x45, 0x28, 0xeb, 0xf9, 0x8a, 0xe9, 0xa9,
	0x2d, 0xdd, 0x99, 0x1c, 0xdf, 0xac, 0x2a, 0xdd, 0x22, 0x16, 0x1f, 0x46, 0x3f, 0x17, 0xfb, 0xde,
	0xaf, 0xc5, 0xbe, 0xf7, 0x67, 0xb1, 0xef, 0x7d, 0xb9, 0x4a, 0x98, 0x99, 0xcc, 0xc7, 0x76, 0xcf,
	0x03, 0x4a, 0xcc, 0x9c, 0x1b, 0x88, 0x27, 0x3a, 0xfb, 0xec, 0x08, 0x30, 0xf7, 0x52, 0x4d, 0x3b,
	0xf9, 0x6b, 0x95, 0x8e, 0x83, 0xb5, 0x0f, 0xd7, 0xb8, 0xee, 0x5e, 0x9d, 0xfe, 0xdf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x70, 0x57, 0x8f, 0x2e, 0xd8, 0x04, 0x00, 0x00,
}

func (m *MsgOption) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgOption) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgOption) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.MsgOwner != nil {
		{
			size, err := m.MsgOwner.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMessage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if len(m.ReceiverPartyId) > 0 {
		i -= len(m.ReceiverPartyId)
		copy(dAtA[i:], m.ReceiverPartyId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.ReceiverPartyId)))
		i--
		dAtA[i] = 0x2a
	}
	if m.ReceiverRole != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.ReceiverRole))
		i--
		dAtA[i] = 0x20
	}
	if len(m.SenderPartyId) > 0 {
		i -= len(m.SenderPartyId)
		copy(dAtA[i:], m.SenderPartyId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.SenderPartyId)))
		i--
		dAtA[i] = 0x1a
	}
	if m.SenderRole != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.SenderRole))
		i--
		dAtA[i] = 0x10
	}
	if len(m.ProposalId) > 0 {
		i -= len(m.ProposalId)
		copy(dAtA[i:], m.ProposalId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.ProposalId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TaskOrganizationIdentityInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TaskOrganizationIdentityInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TaskOrganizationIdentityInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.PartyId) > 0 {
		i -= len(m.PartyId)
		copy(dAtA[i:], m.PartyId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.PartyId)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.IdentityId) > 0 {
		i -= len(m.IdentityId)
		copy(dAtA[i:], m.IdentityId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.IdentityId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.NodeId) > 0 {
		i -= len(m.NodeId)
		copy(dAtA[i:], m.NodeId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.NodeId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ResourceUsage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ResourceUsage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ResourceUsage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.UsedDisk != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.UsedDisk))
		i--
		dAtA[i] = 0x40
	}
	if m.TotalDisk != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.TotalDisk))
		i--
		dAtA[i] = 0x38
	}
	if m.UsedBandwidth != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.UsedBandwidth))
		i--
		dAtA[i] = 0x30
	}
	if m.TotalBandwidth != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.TotalBandwidth))
		i--
		dAtA[i] = 0x28
	}
	if m.UsedProcessor != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.UsedProcessor))
		i--
		dAtA[i] = 0x20
	}
	if m.TotalProcessor != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.TotalProcessor))
		i--
		dAtA[i] = 0x18
	}
	if m.UsedMem != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.UsedMem))
		i--
		dAtA[i] = 0x10
	}
	if m.TotalMem != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.TotalMem))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintMessage(dAtA []byte, offset int, v uint64) int {
	offset -= sovMessage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgOption) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ProposalId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.SenderRole != 0 {
		n += 1 + sovMessage(uint64(m.SenderRole))
	}
	l = len(m.SenderPartyId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.ReceiverRole != 0 {
		n += 1 + sovMessage(uint64(m.ReceiverRole))
	}
	l = len(m.ReceiverPartyId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.MsgOwner != nil {
		l = m.MsgOwner.Size()
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TaskOrganizationIdentityInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.NodeId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.IdentityId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.PartyId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ResourceUsage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TotalMem != 0 {
		n += 1 + sovMessage(uint64(m.TotalMem))
	}
	if m.UsedMem != 0 {
		n += 1 + sovMessage(uint64(m.UsedMem))
	}
	if m.TotalProcessor != 0 {
		n += 1 + sovMessage(uint64(m.TotalProcessor))
	}
	if m.UsedProcessor != 0 {
		n += 1 + sovMessage(uint64(m.UsedProcessor))
	}
	if m.TotalBandwidth != 0 {
		n += 1 + sovMessage(uint64(m.TotalBandwidth))
	}
	if m.UsedBandwidth != 0 {
		n += 1 + sovMessage(uint64(m.UsedBandwidth))
	}
	if m.TotalDisk != 0 {
		n += 1 + sovMessage(uint64(m.TotalDisk))
	}
	if m.UsedDisk != 0 {
		n += 1 + sovMessage(uint64(m.UsedDisk))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMessage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMessage(x uint64) (n int) {
	return sovMessage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgOption) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
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
			return fmt.Errorf("proto: MsgOption: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgOption: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProposalId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
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
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProposalId = append(m.ProposalId[:0], dAtA[iNdEx:postIndex]...)
			if m.ProposalId == nil {
				m.ProposalId = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SenderRole", wireType)
			}
			m.SenderRole = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SenderRole |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SenderPartyId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
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
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SenderPartyId = append(m.SenderPartyId[:0], dAtA[iNdEx:postIndex]...)
			if m.SenderPartyId == nil {
				m.SenderPartyId = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceiverRole", wireType)
			}
			m.ReceiverRole = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReceiverRole |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceiverPartyId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
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
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceiverPartyId = append(m.ReceiverPartyId[:0], dAtA[iNdEx:postIndex]...)
			if m.ReceiverPartyId == nil {
				m.ReceiverPartyId = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgOwner", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MsgOwner == nil {
				m.MsgOwner = &TaskOrganizationIdentityInfo{}
			}
			if err := m.MsgOwner.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
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
func (m *TaskOrganizationIdentityInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
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
			return fmt.Errorf("proto: TaskOrganizationIdentityInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TaskOrganizationIdentityInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
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
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = append(m.Name[:0], dAtA[iNdEx:postIndex]...)
			if m.Name == nil {
				m.Name = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
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
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeId = append(m.NodeId[:0], dAtA[iNdEx:postIndex]...)
			if m.NodeId == nil {
				m.NodeId = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IdentityId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
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
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IdentityId = append(m.IdentityId[:0], dAtA[iNdEx:postIndex]...)
			if m.IdentityId == nil {
				m.IdentityId = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PartyId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
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
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PartyId = append(m.PartyId[:0], dAtA[iNdEx:postIndex]...)
			if m.PartyId == nil {
				m.PartyId = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
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
func (m *ResourceUsage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
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
			return fmt.Errorf("proto: ResourceUsage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ResourceUsage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalMem", wireType)
			}
			m.TotalMem = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalMem |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UsedMem", wireType)
			}
			m.UsedMem = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UsedMem |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalProcessor", wireType)
			}
			m.TotalProcessor = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalProcessor |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UsedProcessor", wireType)
			}
			m.UsedProcessor = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UsedProcessor |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalBandwidth", wireType)
			}
			m.TotalBandwidth = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalBandwidth |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UsedBandwidth", wireType)
			}
			m.UsedBandwidth = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UsedBandwidth |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalDisk", wireType)
			}
			m.TotalDisk = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalDisk |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UsedDisk", wireType)
			}
			m.UsedDisk = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UsedDisk |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMessage
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
func skipMessage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessage
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
					return 0, ErrIntOverflowMessage
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
					return 0, ErrIntOverflowMessage
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
				return 0, ErrInvalidLengthMessage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMessage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMessage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMessage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMessage = fmt.Errorf("proto: unexpected end of group")
)

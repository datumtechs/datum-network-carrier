// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lib/netmsg/taskmng/message.proto

package taskmng

import (
	fmt "fmt"
	common "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
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

// 参与方反馈 各自对某个task的执行结果 (发给发起方)
type TaskResultMsg struct {
	MsgOption            *common.MsgOption `protobuf:"bytes,1,opt,name=msg_option,json=msgOption,proto3" json:"msg_option,omitempty" ssz-max:"16777216"`
	TaskEventList        []*TaskEvent      `protobuf:"bytes,2,rep,name=task_event_list,json=taskEventList,proto3" json:"task_event_list,omitempty" ssz-max:"16777216"`
	CreateAt             uint64            `protobuf:"varint,3,opt,name=create_at,json=createAt,proto3" json:"create_at,omitempty" ssz-size:"8"`
	Sign                 []byte            `protobuf:"bytes,4,opt,name=sign,proto3" json:"sign,omitempty" ssz-max:"1024"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *TaskResultMsg) Reset()         { *m = TaskResultMsg{} }
func (m *TaskResultMsg) String() string { return proto.CompactTextString(m) }
func (*TaskResultMsg) ProtoMessage()    {}
func (*TaskResultMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f475e7ffd8a9ae2, []int{0}
}
func (m *TaskResultMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TaskResultMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TaskResultMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TaskResultMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskResultMsg.Merge(m, src)
}
func (m *TaskResultMsg) XXX_Size() int {
	return m.Size()
}
func (m *TaskResultMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskResultMsg.DiscardUnknown(m)
}

var xxx_messageInfo_TaskResultMsg proto.InternalMessageInfo

func (m *TaskResultMsg) GetMsgOption() *common.MsgOption {
	if m != nil {
		return m.MsgOption
	}
	return nil
}

func (m *TaskResultMsg) GetTaskEventList() []*TaskEvent {
	if m != nil {
		return m.TaskEventList
	}
	return nil
}

func (m *TaskResultMsg) GetCreateAt() uint64 {
	if m != nil {
		return m.CreateAt
	}
	return 0
}

func (m *TaskResultMsg) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

type TaskEvent struct {
	Type                 []byte   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty" ssz-max:"32"`
	TaskId               []byte   `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty" ssz-max:"128"`
	IdentityId           []byte   `protobuf:"bytes,3,opt,name=identity_id,json=identityId,proto3" json:"identity_id,omitempty" ssz-max:"1024"`
	Content              []byte   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty" ssz-max:"2048"`
	CreateAt             uint64   `protobuf:"varint,5,opt,name=create_at,json=createAt,proto3" json:"create_at,omitempty" ssz-size:"8"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskEvent) Reset()         { *m = TaskEvent{} }
func (m *TaskEvent) String() string { return proto.CompactTextString(m) }
func (*TaskEvent) ProtoMessage()    {}
func (*TaskEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f475e7ffd8a9ae2, []int{1}
}
func (m *TaskEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TaskEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TaskEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TaskEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskEvent.Merge(m, src)
}
func (m *TaskEvent) XXX_Size() int {
	return m.Size()
}
func (m *TaskEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskEvent.DiscardUnknown(m)
}

var xxx_messageInfo_TaskEvent proto.InternalMessageInfo

func (m *TaskEvent) GetType() []byte {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *TaskEvent) GetTaskId() []byte {
	if m != nil {
		return m.TaskId
	}
	return nil
}

func (m *TaskEvent) GetIdentityId() []byte {
	if m != nil {
		return m.IdentityId
	}
	return nil
}

func (m *TaskEvent) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *TaskEvent) GetCreateAt() uint64 {
	if m != nil {
		return m.CreateAt
	}
	return 0
}

func init() {
	proto.RegisterType((*TaskResultMsg)(nil), "taskmng.TaskResultMsg")
	proto.RegisterType((*TaskEvent)(nil), "taskmng.TaskEvent")
}

func init() { proto.RegisterFile("lib/netmsg/taskmng/message.proto", fileDescriptor_2f475e7ffd8a9ae2) }

var fileDescriptor_2f475e7ffd8a9ae2 = []byte{
	// 438 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xdf, 0x6e, 0xd3, 0x30,
	0x14, 0xc6, 0x95, 0xb6, 0x6c, 0xd4, 0x6d, 0x81, 0xf9, 0x02, 0x45, 0xbb, 0x68, 0x22, 0x83, 0x50,
	0x85, 0x68, 0xb2, 0xa5, 0xd5, 0x5a, 0xed, 0x8e, 0x20, 0xfe, 0x4c, 0x62, 0x42, 0xb2, 0xb8, 0xe2,
	0x26, 0x4a, 0x5b, 0x63, 0xac, 0xd5, 0x76, 0x95, 0x73, 0x0a, 0x6c, 0xd7, 0x3c, 0x1c, 0x97, 0x3c,
	0x41, 0x85, 0xfa, 0x02, 0x48, 0x7d, 0x02, 0x14, 0x77, 0xa5, 0x4c, 0x50, 0xed, 0x2e, 0xf1, 0xf7,
	0x3b, 0x9f, 0x3f, 0x9f, 0x73, 0x48, 0x38, 0x55, 0xa3, 0xd8, 0x08, 0xd4, 0x20, 0x63, 0xcc, 0xe1,
	0x42, 0x1b, 0x19, 0x6b, 0x01, 0x90, 0x4b, 0x11, 0xcd, 0x0a, 0x8b, 0x96, 0xee, 0x5f, 0x1f, 0x1f,
	0x3e, 0x2a, 0xc4, 0xcc, 0x42, 0xec, 0x4e, 0x47, 0xf3, 0x8f, 0xb1, 0xb4, 0xd2, 0xba, 0x1f, 0xf7,
	0xb5, 0xa6, 0x0f, 0x83, 0xbf, 0xfc, 0xc6, 0x56, 0x6b, 0x6b, 0x6e, 0xda, 0xb1, 0x6f, 0x15, 0xd2,
	0x7a, 0x9f, 0xc3, 0x05, 0x17, 0x30, 0x9f, 0xe2, 0x39, 0x48, 0xfa, 0x86, 0x10, 0x0d, 0x32, 0xb3,
	0x33, 0x54, 0xd6, 0xf8, 0x5e, 0xe8, 0x75, 0x1a, 0xc9, 0x41, 0xb4, 0x2e, 0x8e, 0xce, 0x41, 0xbe,
	0x73, 0x42, 0xfa, 0x70, 0xb5, 0x08, 0x28, 0xc0, 0x55, 0x57, 0xe7, 0x5f, 0x4f, 0xd9, 0xf1, 0xc9,
	0x60, 0x30, 0x48, 0x8e, 0x4f, 0x18, 0xaf, 0xeb, 0x0d, 0x42, 0x39, 0xb9, 0x5f, 0x86, 0xcd, 0xc4,
	0x67, 0x61, 0x30, 0x9b, 0x2a, 0x40, 0xbf, 0x12, 0x56, 0x3b, 0x8d, 0x84, 0x46, 0xd7, 0x8f, 0x88,
	0xca, 0xab, 0x5f, 0x96, 0xf2, 0x4e, 0xbf, 0x16, 0x6e, 0x90, 0xb7, 0x0a, 0x90, 0x76, 0x49, 0x7d,
	0x5c, 0x88, 0x1c, 0x45, 0x96, 0xa3, 0x5f, 0x0d, 0xbd, 0x4e, 0x2d, 0x7d, 0xb0, 0x5a, 0x04, 0xcd,
	0xb2, 0x12, 0xd4, 0x95, 0x38, 0x65, 0x43, 0xc6, 0xef, 0xae, 0x91, 0xe7, 0x48, 0x9f, 0x90, 0x1a,
	0x28, 0x69, 0xfc, 0x5a, 0xe8, 0x75, 0x9a, 0x29, 0x5d, 0x2d, 0x82, 0x7b, 0xdb, 0x3b, 0x8e, 0x92,
	0x3e, 0xe3, 0x4e, 0x67, 0xbf, 0x3c, 0x52, 0xff, 0x93, 0x85, 0x3e, 0x26, 0x35, 0xbc, 0x9c, 0x09,
	0xf7, 0xf8, 0xe6, 0xd6, 0xdf, 0x55, 0xf5, 0x12, 0xc6, 0x9d, 0x4a, 0x9f, 0x12, 0x37, 0x8b, 0x4c,
	0x4d, 0xfc, 0x8a, 0x03, 0x0f, 0x56, 0x8b, 0xa0, 0xb5, 0xb5, 0x4f, 0x86, 0x8c, 0xef, 0x95, 0xc4,
	0xd9, 0x84, 0xf6, 0x48, 0x43, 0x4d, 0x84, 0x41, 0x85, 0x97, 0x25, 0x5f, 0xdd, 0x19, 0x87, 0x6c,
	0xb0, 0xb3, 0x09, 0x7d, 0x46, 0xf6, 0xc7, 0xd6, 0xa0, 0x30, 0xf8, 0xdf, 0xfc, 0xc9, 0x51, 0x7f,
	0xc8, 0xf8, 0x06, 0xb9, 0xd9, 0x99, 0x3b, 0xb7, 0x75, 0x26, 0x4d, 0xbf, 0x2f, 0xdb, 0xde, 0x8f,
	0x65, 0xdb, 0xfb, 0xb9, 0x6c, 0x7b, 0x1f, 0xfa, 0x52, 0xe1, 0xa7, 0xf9, 0xa8, 0x1c, 0x73, 0xcc,
	0x2d, 0x08, 0xc4, 0xfc, 0xd5, 0xd4, 0x7e, 0x89, 0x5f, 0xe4, 0x45, 0xa1, 0x44, 0xd1, 0x7d, 0x6d,
	0xe3, 0x7f, 0x37, 0x73, 0xb4, 0xe7, 0x76, 0xa8, 0xf7, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x69, 0x07,
	0x3a, 0x7c, 0xb6, 0x02, 0x00, 0x00,
}

func (m *TaskResultMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TaskResultMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TaskResultMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Sign) > 0 {
		i -= len(m.Sign)
		copy(dAtA[i:], m.Sign)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Sign)))
		i--
		dAtA[i] = 0x22
	}
	if m.CreateAt != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.CreateAt))
		i--
		dAtA[i] = 0x18
	}
	if len(m.TaskEventList) > 0 {
		for iNdEx := len(m.TaskEventList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TaskEventList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMessage(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.MsgOption != nil {
		{
			size, err := m.MsgOption.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMessage(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TaskEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TaskEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TaskEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.CreateAt != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.CreateAt))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Content) > 0 {
		i -= len(m.Content)
		copy(dAtA[i:], m.Content)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Content)))
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
	if len(m.TaskId) > 0 {
		i -= len(m.TaskId)
		copy(dAtA[i:], m.TaskId)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.TaskId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0xa
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
func (m *TaskResultMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MsgOption != nil {
		l = m.MsgOption.Size()
		n += 1 + l + sovMessage(uint64(l))
	}
	if len(m.TaskEventList) > 0 {
		for _, e := range m.TaskEventList {
			l = e.Size()
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	if m.CreateAt != 0 {
		n += 1 + sovMessage(uint64(m.CreateAt))
	}
	l = len(m.Sign)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TaskEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.TaskId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.IdentityId)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.Content)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.CreateAt != 0 {
		n += 1 + sovMessage(uint64(m.CreateAt))
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
func (m *TaskResultMsg) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TaskResultMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TaskResultMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgOption", wireType)
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
			if m.MsgOption == nil {
				m.MsgOption = &common.MsgOption{}
			}
			if err := m.MsgOption.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TaskEventList", wireType)
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
			m.TaskEventList = append(m.TaskEventList, &TaskEvent{})
			if err := m.TaskEventList[len(m.TaskEventList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreateAt", wireType)
			}
			m.CreateAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreateAt |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sign", wireType)
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
			m.Sign = append(m.Sign[:0], dAtA[iNdEx:postIndex]...)
			if m.Sign == nil {
				m.Sign = []byte{}
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
func (m *TaskEvent) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TaskEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TaskEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
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
			m.Type = append(m.Type[:0], dAtA[iNdEx:postIndex]...)
			if m.Type == nil {
				m.Type = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TaskId", wireType)
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
			m.TaskId = append(m.TaskId[:0], dAtA[iNdEx:postIndex]...)
			if m.TaskId == nil {
				m.TaskId = []byte{}
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
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
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
			m.Content = append(m.Content[:0], dAtA[iNdEx:postIndex]...)
			if m.Content == nil {
				m.Content = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreateAt", wireType)
			}
			m.CreateAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreateAt |= uint64(b&0x7F) << shift
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
// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lib/center/api/identity.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	_ "google.golang.org/protobuf/types/known/emptypb"
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

type SaveIdentityRequest struct {
	Member *Organization `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
	// 节点的身份凭证（DID中的凭证信息）
	Credential           string   `protobuf:"bytes,2,opt,name=credential,proto3" json:"credential,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SaveIdentityRequest) Reset()         { *m = SaveIdentityRequest{} }
func (m *SaveIdentityRequest) String() string { return proto.CompactTextString(m) }
func (*SaveIdentityRequest) ProtoMessage()    {}
func (*SaveIdentityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77fa157ad346cc0a, []int{0}
}
func (m *SaveIdentityRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SaveIdentityRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SaveIdentityRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SaveIdentityRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveIdentityRequest.Merge(m, src)
}
func (m *SaveIdentityRequest) XXX_Size() int {
	return m.Size()
}
func (m *SaveIdentityRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveIdentityRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SaveIdentityRequest proto.InternalMessageInfo

func (m *SaveIdentityRequest) GetMember() *Organization {
	if m != nil {
		return m.Member
	}
	return nil
}

func (m *SaveIdentityRequest) GetCredential() string {
	if m != nil {
		return m.Credential
	}
	return ""
}

type RevokeIdentityJoinRequest struct {
	Member               *Organization `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RevokeIdentityJoinRequest) Reset()         { *m = RevokeIdentityJoinRequest{} }
func (m *RevokeIdentityJoinRequest) String() string { return proto.CompactTextString(m) }
func (*RevokeIdentityJoinRequest) ProtoMessage()    {}
func (*RevokeIdentityJoinRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77fa157ad346cc0a, []int{1}
}
func (m *RevokeIdentityJoinRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RevokeIdentityJoinRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RevokeIdentityJoinRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RevokeIdentityJoinRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RevokeIdentityJoinRequest.Merge(m, src)
}
func (m *RevokeIdentityJoinRequest) XXX_Size() int {
	return m.Size()
}
func (m *RevokeIdentityJoinRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RevokeIdentityJoinRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RevokeIdentityJoinRequest proto.InternalMessageInfo

func (m *RevokeIdentityJoinRequest) GetMember() *Organization {
	if m != nil {
		return m.Member
	}
	return nil
}

type IdentityListRequest struct {
	// 同步时间点，用于进行数据增量拉去
	LastUpdateTime       uint64   `protobuf:"varint,1,opt,name=last_update_time,json=lastUpdateTime,proto3" json:"last_update_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdentityListRequest) Reset()         { *m = IdentityListRequest{} }
func (m *IdentityListRequest) String() string { return proto.CompactTextString(m) }
func (*IdentityListRequest) ProtoMessage()    {}
func (*IdentityListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77fa157ad346cc0a, []int{2}
}
func (m *IdentityListRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IdentityListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IdentityListRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IdentityListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdentityListRequest.Merge(m, src)
}
func (m *IdentityListRequest) XXX_Size() int {
	return m.Size()
}
func (m *IdentityListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IdentityListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IdentityListRequest proto.InternalMessageInfo

func (m *IdentityListRequest) GetLastUpdateTime() uint64 {
	if m != nil {
		return m.LastUpdateTime
	}
	return 0
}

type IdentityListResponse struct {
	IdentityList []*Organization `protobuf:"bytes,1,rep,name=identity_list,json=identityList,proto3" json:"identity_list,omitempty"`
	// 数据的最后更新点（秒级时间戳）
	LastUpdateTime       uint64   `protobuf:"varint,2,opt,name=last_update_time,json=lastUpdateTime,proto3" json:"last_update_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdentityListResponse) Reset()         { *m = IdentityListResponse{} }
func (m *IdentityListResponse) String() string { return proto.CompactTextString(m) }
func (*IdentityListResponse) ProtoMessage()    {}
func (*IdentityListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77fa157ad346cc0a, []int{3}
}
func (m *IdentityListResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IdentityListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IdentityListResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IdentityListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdentityListResponse.Merge(m, src)
}
func (m *IdentityListResponse) XXX_Size() int {
	return m.Size()
}
func (m *IdentityListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IdentityListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IdentityListResponse proto.InternalMessageInfo

func (m *IdentityListResponse) GetIdentityList() []*Organization {
	if m != nil {
		return m.IdentityList
	}
	return nil
}

func (m *IdentityListResponse) GetLastUpdateTime() uint64 {
	if m != nil {
		return m.LastUpdateTime
	}
	return 0
}

func init() {
	proto.RegisterType((*SaveIdentityRequest)(nil), "api.SaveIdentityRequest")
	proto.RegisterType((*RevokeIdentityJoinRequest)(nil), "api.RevokeIdentityJoinRequest")
	proto.RegisterType((*IdentityListRequest)(nil), "api.IdentityListRequest")
	proto.RegisterType((*IdentityListResponse)(nil), "api.IdentityListResponse")
}

func init() { proto.RegisterFile("lib/center/api/identity.proto", fileDescriptor_77fa157ad346cc0a) }

var fileDescriptor_77fa157ad346cc0a = []byte{
	// 382 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xb5, 0x09, 0x8a, 0xc4, 0x12, 0x08, 0x6c, 0x38, 0x38, 0x46, 0x58, 0x91, 0x4f, 0xe1,
	0x80, 0x17, 0x05, 0x89, 0x1b, 0x20, 0x81, 0x94, 0x28, 0x08, 0x09, 0xc9, 0x81, 0x0b, 0x97, 0xb0,
	0x76, 0x06, 0x33, 0xc2, 0xf6, 0x6e, 0x77, 0xd7, 0xe9, 0x9f, 0x27, 0xec, 0xb1, 0x8f, 0x50, 0xa5,
	0x2f, 0x52, 0xd5, 0xb1, 0x23, 0xa7, 0x75, 0x0f, 0x3d, 0x7a, 0x66, 0xbe, 0xdf, 0x37, 0xde, 0x6f,
	0xe8, 0xeb, 0x14, 0x23, 0x1e, 0x43, 0x6e, 0x41, 0x73, 0xa1, 0x90, 0xe3, 0x1a, 0x72, 0x8b, 0xf6,
	0x34, 0x50, 0x5a, 0x5a, 0xc9, 0xba, 0x42, 0xa1, 0x3b, 0xba, 0x35, 0x13, 0x09, 0x03, 0xbb, 0xbe,
	0xfb, 0x2a, 0x91, 0x32, 0x49, 0x81, 0x97, 0x5f, 0x51, 0xf1, 0x97, 0x43, 0xa6, 0x6a, 0xb1, 0xff,
	0x87, 0x0e, 0x97, 0x62, 0x03, 0x8b, 0x0a, 0x19, 0xc2, 0x51, 0x01, 0xc6, 0xb2, 0x37, 0xb4, 0x97,
	0x41, 0x16, 0x81, 0x76, 0xc8, 0x98, 0x4c, 0x9e, 0x4c, 0x5f, 0x04, 0x42, 0x61, 0xf0, 0x43, 0x27,
	0x22, 0xc7, 0x33, 0x61, 0x51, 0xe6, 0x61, 0x35, 0xc0, 0x3c, 0x4a, 0x63, 0x0d, 0xa5, 0x5e, 0xa4,
	0x4e, 0x67, 0x4c, 0x26, 0x8f, 0xc3, 0x46, 0xc5, 0x9f, 0xd1, 0x51, 0x08, 0x1b, 0xf9, 0x7f, 0xef,
	0xf1, 0x4d, 0x62, 0xfe, 0x70, 0x1f, 0xff, 0x33, 0x1d, 0xd6, 0x84, 0xef, 0x68, 0x6c, 0x4d, 0x98,
	0xd0, 0xe7, 0xa9, 0x30, 0x76, 0x55, 0xa8, 0xb5, 0xb0, 0xb0, 0xb2, 0x98, 0x41, 0xc9, 0x7a, 0x14,
	0x3e, 0xbb, 0xa9, 0xff, 0x2a, 0xcb, 0x3f, 0x31, 0x03, 0xff, 0x84, 0xbe, 0x3c, 0x04, 0x18, 0x25,
	0x73, 0x03, 0xec, 0x03, 0x7d, 0x5a, 0xbf, 0xe8, 0x2a, 0x45, 0x63, 0x1d, 0x32, 0xee, 0xb6, 0xaf,
	0xd2, 0xc7, 0x86, 0xbe, 0xd5, 0xb9, 0xd3, 0xe6, 0x3c, 0xbd, 0x22, 0x74, 0x50, 0x5b, 0x2f, 0x41,
	0x6f, 0x30, 0x06, 0x36, 0xa3, 0x83, 0x39, 0xd8, 0xe6, 0x42, 0xcc, 0x29, 0x1d, 0x5b, 0x7e, 0xd2,
	0x1d, 0xb5, 0x74, 0xaa, 0xed, 0x3f, 0xd2, 0x7e, 0x33, 0xc0, 0x0a, 0xd2, 0x92, 0xa9, 0x3b, 0xdc,
	0x75, 0x30, 0x53, 0x29, 0xec, 0xe5, 0x0b, 0xca, 0xee, 0xa6, 0xc3, 0xbc, 0x72, 0xf4, 0xde, 0xd8,
	0x5a, 0x51, 0x5f, 0x3e, 0x9d, 0x6f, 0x3d, 0x72, 0xb1, 0xf5, 0xc8, 0xe5, 0xd6, 0x23, 0xbf, 0xdf,
	0x25, 0x68, 0xff, 0x15, 0x51, 0x10, 0xcb, 0x8c, 0x87, 0xd2, 0x80, 0xb5, 0x62, 0x96, 0xca, 0x63,
	0xfe, 0x55, 0x68, 0x8d, 0xa0, 0xdf, 0xce, 0x25, 0x3f, 0x3c, 0xd9, 0xa8, 0x57, 0x5e, 0xe4, 0xfb,
	0xeb, 0x00, 0x00, 0x00, 0xff, 0xff, 0x44, 0x1c, 0x2b, 0x39, 0xef, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IdentityServiceClient is the client API for IdentityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IdentityServiceClient interface {
	// 拉去所有的身份数据
	GetIdentityList(ctx context.Context, in *IdentityListRequest, opts ...grpc.CallOption) (*IdentityListResponse, error)
	// 存储身份信息（节点用于申请接入网络的基本信息，详细的存于本地）
	SaveIdentity(ctx context.Context, in *SaveIdentityRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	// 注销准入网络
	RevokeIdentityJoin(ctx context.Context, in *RevokeIdentityJoinRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
}

type identityServiceClient struct {
	cc *grpc.ClientConn
}

func NewIdentityServiceClient(cc *grpc.ClientConn) IdentityServiceClient {
	return &identityServiceClient{cc}
}

func (c *identityServiceClient) GetIdentityList(ctx context.Context, in *IdentityListRequest, opts ...grpc.CallOption) (*IdentityListResponse, error) {
	out := new(IdentityListResponse)
	err := c.cc.Invoke(ctx, "/api.IdentityService/GetIdentityList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityServiceClient) SaveIdentity(ctx context.Context, in *SaveIdentityRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/api.IdentityService/SaveIdentity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityServiceClient) RevokeIdentityJoin(ctx context.Context, in *RevokeIdentityJoinRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/api.IdentityService/RevokeIdentityJoin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IdentityServiceServer is the server API for IdentityService service.
type IdentityServiceServer interface {
	// 拉去所有的身份数据
	GetIdentityList(context.Context, *IdentityListRequest) (*IdentityListResponse, error)
	// 存储身份信息（节点用于申请接入网络的基本信息，详细的存于本地）
	SaveIdentity(context.Context, *SaveIdentityRequest) (*SimpleResponse, error)
	// 注销准入网络
	RevokeIdentityJoin(context.Context, *RevokeIdentityJoinRequest) (*SimpleResponse, error)
}

// UnimplementedIdentityServiceServer can be embedded to have forward compatible implementations.
type UnimplementedIdentityServiceServer struct {
}

func (*UnimplementedIdentityServiceServer) GetIdentityList(ctx context.Context, req *IdentityListRequest) (*IdentityListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIdentityList not implemented")
}
func (*UnimplementedIdentityServiceServer) SaveIdentity(ctx context.Context, req *SaveIdentityRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveIdentity not implemented")
}
func (*UnimplementedIdentityServiceServer) RevokeIdentityJoin(ctx context.Context, req *RevokeIdentityJoinRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeIdentityJoin not implemented")
}

func RegisterIdentityServiceServer(s *grpc.Server, srv IdentityServiceServer) {
	s.RegisterService(&_IdentityService_serviceDesc, srv)
}

func _IdentityService_GetIdentityList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdentityListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).GetIdentityList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.IdentityService/GetIdentityList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).GetIdentityList(ctx, req.(*IdentityListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityService_SaveIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveIdentityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).SaveIdentity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.IdentityService/SaveIdentity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).SaveIdentity(ctx, req.(*SaveIdentityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityService_RevokeIdentityJoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeIdentityJoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityServiceServer).RevokeIdentityJoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.IdentityService/RevokeIdentityJoin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityServiceServer).RevokeIdentityJoin(ctx, req.(*RevokeIdentityJoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IdentityService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.IdentityService",
	HandlerType: (*IdentityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetIdentityList",
			Handler:    _IdentityService_GetIdentityList_Handler,
		},
		{
			MethodName: "SaveIdentity",
			Handler:    _IdentityService_SaveIdentity_Handler,
		},
		{
			MethodName: "RevokeIdentityJoin",
			Handler:    _IdentityService_RevokeIdentityJoin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lib/center/api/identity.proto",
}

func (m *SaveIdentityRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SaveIdentityRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SaveIdentityRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Credential) > 0 {
		i -= len(m.Credential)
		copy(dAtA[i:], m.Credential)
		i = encodeVarintIdentity(dAtA, i, uint64(len(m.Credential)))
		i--
		dAtA[i] = 0x12
	}
	if m.Member != nil {
		{
			size, err := m.Member.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIdentity(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RevokeIdentityJoinRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RevokeIdentityJoinRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RevokeIdentityJoinRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Member != nil {
		{
			size, err := m.Member.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIdentity(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *IdentityListRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IdentityListRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IdentityListRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.LastUpdateTime != 0 {
		i = encodeVarintIdentity(dAtA, i, uint64(m.LastUpdateTime))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *IdentityListResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IdentityListResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IdentityListResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.LastUpdateTime != 0 {
		i = encodeVarintIdentity(dAtA, i, uint64(m.LastUpdateTime))
		i--
		dAtA[i] = 0x10
	}
	if len(m.IdentityList) > 0 {
		for iNdEx := len(m.IdentityList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.IdentityList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintIdentity(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintIdentity(dAtA []byte, offset int, v uint64) int {
	offset -= sovIdentity(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SaveIdentityRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Member != nil {
		l = m.Member.Size()
		n += 1 + l + sovIdentity(uint64(l))
	}
	l = len(m.Credential)
	if l > 0 {
		n += 1 + l + sovIdentity(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *RevokeIdentityJoinRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Member != nil {
		l = m.Member.Size()
		n += 1 + l + sovIdentity(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *IdentityListRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.LastUpdateTime != 0 {
		n += 1 + sovIdentity(uint64(m.LastUpdateTime))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *IdentityListResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.IdentityList) > 0 {
		for _, e := range m.IdentityList {
			l = e.Size()
			n += 1 + l + sovIdentity(uint64(l))
		}
	}
	if m.LastUpdateTime != 0 {
		n += 1 + sovIdentity(uint64(m.LastUpdateTime))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovIdentity(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozIdentity(x uint64) (n int) {
	return sovIdentity(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SaveIdentityRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentity
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
			return fmt.Errorf("proto: SaveIdentityRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SaveIdentityRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Member", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentity
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
				return ErrInvalidLengthIdentity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Member == nil {
				m.Member = &Organization{}
			}
			if err := m.Member.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Credential", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIdentity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Credential = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIdentity
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
func (m *RevokeIdentityJoinRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentity
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
			return fmt.Errorf("proto: RevokeIdentityJoinRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RevokeIdentityJoinRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Member", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentity
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
				return ErrInvalidLengthIdentity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Member == nil {
				m.Member = &Organization{}
			}
			if err := m.Member.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIdentity
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
func (m *IdentityListRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentity
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
			return fmt.Errorf("proto: IdentityListRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IdentityListRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastUpdateTime", wireType)
			}
			m.LastUpdateTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastUpdateTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipIdentity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIdentity
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
func (m *IdentityListResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentity
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
			return fmt.Errorf("proto: IdentityListResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IdentityListResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IdentityList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentity
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
				return ErrInvalidLengthIdentity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IdentityList = append(m.IdentityList, &Organization{})
			if err := m.IdentityList[len(m.IdentityList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastUpdateTime", wireType)
			}
			m.LastUpdateTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastUpdateTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipIdentity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIdentity
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
func skipIdentity(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIdentity
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
					return 0, ErrIntOverflowIdentity
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
					return 0, ErrIntOverflowIdentity
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
				return 0, ErrInvalidLengthIdentity
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupIdentity
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthIdentity
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthIdentity        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIdentity          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupIdentity = fmt.Errorf("proto: unexpected end of group")
)

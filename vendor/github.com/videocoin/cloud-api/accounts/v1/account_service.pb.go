// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: accounts/v1/account_service.proto

package v1

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	rpc "github.com/videocoin/cloud-api/rpc"
	grpc "google.golang.org/grpc"
	io "io"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type AccountRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OwnerId              string   `protobuf:"bytes,2,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountRequest) Reset()         { *m = AccountRequest{} }
func (m *AccountRequest) String() string { return proto.CompactTextString(m) }
func (*AccountRequest) ProtoMessage()    {}
func (*AccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a57b966a6f05cc7, []int{0}
}
func (m *AccountRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountRequest.Merge(m, src)
}
func (m *AccountRequest) XXX_Size() int {
	return m.Size()
}
func (m *AccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AccountRequest proto.InternalMessageInfo

func (m *AccountRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AccountRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (*AccountRequest) XXX_MessageName() string {
	return "cloud.api.account.v1.AccountRequest"
}

type Address struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a57b966a6f05cc7, []int{1}
}
func (m *Address) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Address) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Address.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Address) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Address.Merge(m, src)
}
func (m *Address) XXX_Size() int {
	return m.Size()
}
func (m *Address) XXX_DiscardUnknown() {
	xxx_messageInfo_Address.DiscardUnknown(m)
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (*Address) XXX_MessageName() string {
	return "cloud.api.account.v1.Address"
}

type ListResponse struct {
	Items                []*AccountProfile `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ListResponse) Reset()         { *m = ListResponse{} }
func (m *ListResponse) String() string { return proto.CompactTextString(m) }
func (*ListResponse) ProtoMessage()    {}
func (*ListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a57b966a6f05cc7, []int{2}
}
func (m *ListResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ListResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListResponse.Merge(m, src)
}
func (m *ListResponse) XXX_Size() int {
	return m.Size()
}
func (m *ListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListResponse proto.InternalMessageInfo

func (m *ListResponse) GetItems() []*AccountProfile {
	if m != nil {
		return m.Items
	}
	return nil
}

func (*ListResponse) XXX_MessageName() string {
	return "cloud.api.account.v1.ListResponse"
}
func init() {
	proto.RegisterType((*AccountRequest)(nil), "cloud.api.account.v1.AccountRequest")
	golang_proto.RegisterType((*AccountRequest)(nil), "cloud.api.account.v1.AccountRequest")
	proto.RegisterType((*Address)(nil), "cloud.api.account.v1.Address")
	golang_proto.RegisterType((*Address)(nil), "cloud.api.account.v1.Address")
	proto.RegisterType((*ListResponse)(nil), "cloud.api.account.v1.ListResponse")
	golang_proto.RegisterType((*ListResponse)(nil), "cloud.api.account.v1.ListResponse")
}

func init() { proto.RegisterFile("accounts/v1/account_service.proto", fileDescriptor_0a57b966a6f05cc7) }
func init() {
	golang_proto.RegisterFile("accounts/v1/account_service.proto", fileDescriptor_0a57b966a6f05cc7)
}

var fileDescriptor_0a57b966a6f05cc7 = []byte{
	// 481 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x52, 0xbd, 0x6e, 0x13, 0x41,
	0x10, 0xf6, 0x9d, 0x83, 0x1d, 0x86, 0x28, 0x42, 0x2b, 0x7e, 0x1c, 0x87, 0x9c, 0xc2, 0x41, 0x91,
	0x26, 0x7b, 0x72, 0xe8, 0xa0, 0x4a, 0x00, 0x05, 0x08, 0x92, 0x91, 0x23, 0x28, 0xd2, 0x44, 0xeb,
	0xdb, 0xc9, 0x79, 0x25, 0xfb, 0xf6, 0xb8, 0xdd, 0x3b, 0x64, 0x10, 0x0d, 0xaf, 0xc0, 0x0b, 0x51,
	0xa6, 0x44, 0xe2, 0x05, 0x90, 0xc3, 0x2b, 0xd0, 0xa3, 0xdb, 0x5d, 0x13, 0x17, 0x4e, 0xe4, 0xc2,
	0xdd, 0xcc, 0xce, 0xf7, 0xcd, 0x7c, 0xfb, 0xcd, 0xc0, 0x43, 0x16, 0xc7, 0xb2, 0x48, 0xb5, 0x8a,
	0xca, 0x4e, 0xe4, 0xe2, 0x53, 0x85, 0x79, 0x29, 0x62, 0xa4, 0x59, 0x2e, 0xb5, 0x24, 0x77, 0xe2,
	0xa1, 0x2c, 0x38, 0x65, 0x99, 0xa0, 0x0e, 0x40, 0xcb, 0x4e, 0x7b, 0x63, 0x0e, 0xd1, 0x12, 0xda,
	0x51, 0x22, 0xf4, 0xa0, 0xe8, 0xd3, 0x58, 0x8e, 0xa2, 0x52, 0x70, 0x94, 0xb1, 0x14, 0x69, 0x64,
	0xba, 0xec, 0xb2, 0x4c, 0x44, 0x79, 0x16, 0x47, 0x03, 0x64, 0x43, 0x3d, 0x70, 0x84, 0xdd, 0x19,
	0x42, 0x22, 0x13, 0x19, 0x99, 0xe7, 0x7e, 0x71, 0x66, 0x32, 0x93, 0x98, 0xc8, 0xc1, 0x37, 0x13,
	0x29, 0x93, 0x21, 0x5e, 0xa2, 0x70, 0x94, 0xe9, 0xb1, 0x2b, 0x3e, 0x70, 0xc5, 0x6a, 0x10, 0x4b,
	0x53, 0xa9, 0x99, 0x16, 0x32, 0x55, 0xb6, 0x1a, 0x3e, 0x83, 0xf5, 0x7d, 0xab, 0xb5, 0x87, 0x1f,
	0x0b, 0x54, 0x9a, 0xac, 0x83, 0x2f, 0x78, 0xcb, 0xdb, 0xf6, 0x76, 0x6e, 0xf6, 0x7c, 0xc1, 0xc9,
	0x06, 0xac, 0xca, 0x4f, 0x29, 0xe6, 0xa7, 0x82, 0xb7, 0x7c, 0xf3, 0xda, 0x34, 0xf9, 0x6b, 0x1e,
	0x3e, 0x82, 0xe6, 0x3e, 0xe7, 0x39, 0x2a, 0x45, 0x5a, 0xd0, 0x64, 0x36, 0x74, 0xd4, 0x69, 0x1a,
	0xbe, 0x81, 0xb5, 0xb7, 0x42, 0xe9, 0x1e, 0xaa, 0x4c, 0xa6, 0x0a, 0xc9, 0x53, 0xb8, 0x21, 0x34,
	0x8e, 0x2a, 0x5c, 0x7d, 0xe7, 0xd6, 0xde, 0x63, 0x3a, 0xcf, 0x4d, 0xea, 0x44, 0xbd, 0xcb, 0xe5,
	0x99, 0x18, 0x62, 0xcf, 0x52, 0xf6, 0xfe, 0xae, 0xfc, 0x97, 0x7b, 0x6c, 0x57, 0x42, 0xba, 0xd0,
	0x78, 0x65, 0xac, 0x23, 0xf7, 0xa8, 0xfd, 0x29, 0x9d, 0xda, 0x40, 0x5f, 0x56, 0x36, 0xb4, 0x37,
	0x67, 0x26, 0xe4, 0x59, 0x4c, 0x2d, 0xfc, 0x58, 0x33, 0x5d, 0xa8, 0xf0, 0xf6, 0xb7, 0x5f, 0x7f,
	0xbe, 0xfb, 0x40, 0x56, 0xdd, 0x02, 0x3e, 0x93, 0x0c, 0xea, 0x87, 0xa8, 0xc9, 0xf5, 0xba, 0x9c,
	0x59, 0xed, 0x85, 0xd4, 0x87, 0x5b, 0x66, 0xc8, 0x7d, 0x72, 0xd7, 0x2c, 0xe1, 0xf2, 0x3c, 0x54,
	0xf4, 0x45, 0xf0, 0xaf, 0xe4, 0x3d, 0xac, 0x1d, 0xa2, 0x3e, 0x18, 0x4f, 0xbd, 0xdc, 0xba, 0xa2,
	0xa9, 0x2d, 0x2f, 0x38, 0xb3, 0x46, 0x4e, 0x00, 0x4c, 0xdb, 0x6e, 0xb5, 0xad, 0xa5, 0xfe, 0xa7,
	0x46, 0x3e, 0x40, 0xe3, 0x79, 0x8e, 0x4c, 0xe3, 0x92, 0xfb, 0xbe, 0x80, 0x95, 0xea, 0x58, 0xae,
	0xdc, 0x65, 0x38, 0xbf, 0xcf, 0xec, 0x81, 0x85, 0x35, 0xd2, 0x85, 0xfa, 0x11, 0x8e, 0x17, 0x94,
	0xb6, 0x7d, 0x2d, 0xea, 0x08, 0xc7, 0x61, 0xed, 0xa0, 0x75, 0x3e, 0x09, 0xbc, 0x9f, 0x93, 0xc0,
	0xfb, 0x3d, 0x09, 0xbc, 0x1f, 0x17, 0x81, 0x77, 0x7e, 0x11, 0x78, 0x27, 0x7e, 0xd9, 0xe9, 0x37,
	0x8c, 0xc0, 0x27, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0xe2, 0x41, 0xe3, 0x10, 0x37, 0x04, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AccountServiceClient interface {
	Health(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*rpc.HealthStatus, error)
	Get(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountProfile, error)
	GetByAddress(ctx context.Context, in *Address, opts ...grpc.CallOption) (*AccountProfile, error)
	GetByOwner(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountProfile, error)
	Create(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountProfile, error)
	List(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*ListResponse, error)
	Key(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountKey, error)
}

type accountServiceClient struct {
	cc *grpc.ClientConn
}

func NewAccountServiceClient(cc *grpc.ClientConn) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) Health(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*rpc.HealthStatus, error) {
	out := new(rpc.HealthStatus)
	err := c.cc.Invoke(ctx, "/cloud.api.account.v1.AccountService/Health", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Get(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountProfile, error) {
	out := new(AccountProfile)
	err := c.cc.Invoke(ctx, "/cloud.api.account.v1.AccountService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetByAddress(ctx context.Context, in *Address, opts ...grpc.CallOption) (*AccountProfile, error) {
	out := new(AccountProfile)
	err := c.cc.Invoke(ctx, "/cloud.api.account.v1.AccountService/GetByAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetByOwner(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountProfile, error) {
	out := new(AccountProfile)
	err := c.cc.Invoke(ctx, "/cloud.api.account.v1.AccountService/GetByOwner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Create(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountProfile, error) {
	out := new(AccountProfile)
	err := c.cc.Invoke(ctx, "/cloud.api.account.v1.AccountService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) List(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/cloud.api.account.v1.AccountService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Key(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountKey, error) {
	out := new(AccountKey)
	err := c.cc.Invoke(ctx, "/cloud.api.account.v1.AccountService/Key", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
type AccountServiceServer interface {
	Health(context.Context, *types.Empty) (*rpc.HealthStatus, error)
	Get(context.Context, *AccountRequest) (*AccountProfile, error)
	GetByAddress(context.Context, *Address) (*AccountProfile, error)
	GetByOwner(context.Context, *AccountRequest) (*AccountProfile, error)
	Create(context.Context, *AccountRequest) (*AccountProfile, error)
	List(context.Context, *types.Empty) (*ListResponse, error)
	Key(context.Context, *AccountRequest) (*AccountKey, error)
}

func RegisterAccountServiceServer(s *grpc.Server, srv AccountServiceServer) {
	s.RegisterService(&_AccountService_serviceDesc, srv)
}

func _AccountService_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.api.account.v1.AccountService/Health",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Health(ctx, req.(*types.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.api.account.v1.AccountService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Get(ctx, req.(*AccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetByAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetByAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.api.account.v1.AccountService/GetByAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetByAddress(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetByOwner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetByOwner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.api.account.v1.AccountService/GetByOwner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetByOwner(ctx, req.(*AccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.api.account.v1.AccountService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Create(ctx, req.(*AccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.api.account.v1.AccountService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).List(ctx, req.(*types.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Key_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Key(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.api.account.v1.AccountService/Key",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Key(ctx, req.(*AccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AccountService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cloud.api.account.v1.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Health",
			Handler:    _AccountService_Health_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _AccountService_Get_Handler,
		},
		{
			MethodName: "GetByAddress",
			Handler:    _AccountService_GetByAddress_Handler,
		},
		{
			MethodName: "GetByOwner",
			Handler:    _AccountService_GetByOwner_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _AccountService_Create_Handler,
		},
		{
			MethodName: "List",
			Handler:    _AccountService_List_Handler,
		},
		{
			MethodName: "Key",
			Handler:    _AccountService_Key_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "accounts/v1/account_service.proto",
}

func (m *AccountRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAccountService(dAtA, i, uint64(len(m.Id)))
		i += copy(dAtA[i:], m.Id)
	}
	if len(m.OwnerId) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintAccountService(dAtA, i, uint64(len(m.OwnerId)))
		i += copy(dAtA[i:], m.OwnerId)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Address) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Address) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAccountService(dAtA, i, uint64(len(m.Address)))
		i += copy(dAtA[i:], m.Address)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *ListResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ListResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Items) > 0 {
		for _, msg := range m.Items {
			dAtA[i] = 0xa
			i++
			i = encodeVarintAccountService(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintAccountService(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *AccountRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovAccountService(uint64(l))
	}
	l = len(m.OwnerId)
	if l > 0 {
		n += 1 + l + sovAccountService(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Address) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovAccountService(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ListResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovAccountService(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovAccountService(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozAccountService(x uint64) (n int) {
	return sovAccountService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AccountRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccountService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AccountRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAccountService
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OwnerId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAccountService
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OwnerId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAccountService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAccountService
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
func (m *Address) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccountService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Address: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Address: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAccountService
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAccountService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAccountService
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
func (m *ListResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccountService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ListResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ListResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAccountService
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, &AccountProfile{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAccountService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAccountService
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
func skipAccountService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAccountService
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
					return 0, ErrIntOverflowAccountService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAccountService
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthAccountService
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowAccountService
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipAccountService(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthAccountService = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAccountService   = fmt.Errorf("proto: integer overflow")
)

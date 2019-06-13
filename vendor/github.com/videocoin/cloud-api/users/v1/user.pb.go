// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: users/v1/user.proto

package v1

import (
	fmt "fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	v1 "github.com/videocoin/cloud-api/accounts/v1"
	io "io"
	math "math"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type UserRole int32

const (
	UserRoleRegular UserRole = 0
	UserRoleQa      UserRole = 3
	UserRoleManager UserRole = 6
	UserRoleSuper   UserRole = 9
)

var UserRole_name = map[int32]string{
	0: "REGULAR",
	3: "QA",
	6: "MANAGER",
	9: "SUPER",
}

var UserRole_value = map[string]int32{
	"REGULAR": 0,
	"QA":      3,
	"MANAGER": 6,
	"SUPER":   9,
}

func (x UserRole) String() string {
	return proto.EnumName(UserRole_name, int32(x))
}

func (UserRole) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_622714e2df60ae10, []int{0}
}

type User struct {
	Id                   string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" gorm:"type:varchar(36);primary_key"`
	Email                string     `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty" gorm:"unique_index"`
	Password             string     `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty" gorm:"type:varchar(100)"`
	Name                 string     `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty" gorm:"type:varchar(100)"`
	Role                 UserRole   `protobuf:"varint,5,opt,name=role,proto3,enum=cloud.api.users.v1.UserRole" json:"role,omitempty"`
	IsActive             bool       `protobuf:"varint,6,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	ActivatedAt          *time.Time `protobuf:"bytes,7,opt,name=activated_at,json=activatedAt,proto3,stdtime" json:"activated_at,omitempty"`
	CreatedAt            *time.Time `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3,stdtime" json:"created_at,omitempty"`
	Token                string     `protobuf:"bytes,13,opt,name=token,proto3" json:"token,omitempty" gorm:"type:varchar(255);index"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_622714e2df60ae10, []int{0}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_User.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return m.Size()
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetRole() UserRole {
	if m != nil {
		return m.Role
	}
	return UserRoleRegular
}

func (m *User) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

func (m *User) GetActivatedAt() *time.Time {
	if m != nil {
		return m.ActivatedAt
	}
	return nil
}

func (m *User) GetCreatedAt() *time.Time {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *User) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (*User) XXX_MessageName() string {
	return "cloud.api.users.v1.User"
}

type UserProfile struct {
	Id                   string             `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Email                string             `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Name                 string             `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	IsActive             bool               `protobuf:"varint,4,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	Account              *v1.AccountProfile `protobuf:"bytes,5,opt,name=account,proto3" json:"account,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *UserProfile) Reset()         { *m = UserProfile{} }
func (m *UserProfile) String() string { return proto.CompactTextString(m) }
func (*UserProfile) ProtoMessage()    {}
func (*UserProfile) Descriptor() ([]byte, []int) {
	return fileDescriptor_622714e2df60ae10, []int{1}
}
func (m *UserProfile) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UserProfile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UserProfile.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UserProfile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserProfile.Merge(m, src)
}
func (m *UserProfile) XXX_Size() int {
	return m.Size()
}
func (m *UserProfile) XXX_DiscardUnknown() {
	xxx_messageInfo_UserProfile.DiscardUnknown(m)
}

var xxx_messageInfo_UserProfile proto.InternalMessageInfo

func (m *UserProfile) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UserProfile) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserProfile) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UserProfile) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

func (m *UserProfile) GetAccount() *v1.AccountProfile {
	if m != nil {
		return m.Account
	}
	return nil
}

func (*UserProfile) XXX_MessageName() string {
	return "cloud.api.users.v1.UserProfile"
}
func init() {
	proto.RegisterEnum("cloud.api.users.v1.UserRole", UserRole_name, UserRole_value)
	golang_proto.RegisterEnum("cloud.api.users.v1.UserRole", UserRole_name, UserRole_value)
	proto.RegisterType((*User)(nil), "cloud.api.users.v1.User")
	golang_proto.RegisterType((*User)(nil), "cloud.api.users.v1.User")
	proto.RegisterType((*UserProfile)(nil), "cloud.api.users.v1.UserProfile")
	golang_proto.RegisterType((*UserProfile)(nil), "cloud.api.users.v1.UserProfile")
}

func init() { proto.RegisterFile("users/v1/user.proto", fileDescriptor_622714e2df60ae10) }
func init() { golang_proto.RegisterFile("users/v1/user.proto", fileDescriptor_622714e2df60ae10) }

var fileDescriptor_622714e2df60ae10 = []byte{
	// 648 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xc1, 0x52, 0xd3, 0x40,
	0x18, 0xc7, 0x49, 0x69, 0xa1, 0xdd, 0x0a, 0xd6, 0xc5, 0xd1, 0x4c, 0xec, 0xa4, 0x99, 0xe8, 0x8c,
	0xc5, 0xb1, 0x09, 0x2d, 0x83, 0x32, 0x30, 0xa3, 0x13, 0x1c, 0x86, 0x8b, 0x38, 0x10, 0xe4, 0xe2,
	0xa5, 0xb3, 0x24, 0x4b, 0xd8, 0x21, 0xc9, 0xc6, 0xcd, 0x26, 0xc8, 0x1b, 0x38, 0x9c, 0x3c, 0x79,
	0xe3, 0x24, 0x4f, 0xe1, 0xc9, 0x23, 0x47, 0x9f, 0xa0, 0x3a, 0x70, 0xf2, 0xda, 0x27, 0x70, 0xba,
	0x4d, 0x0a, 0x68, 0x0f, 0x9c, 0xfa, 0x7d, 0xf3, 0xfd, 0xff, 0xdf, 0xee, 0xfe, 0x7f, 0x0d, 0x98,
	0x4b, 0x62, 0xcc, 0x62, 0x33, 0x6d, 0x9b, 0x83, 0xc2, 0x88, 0x18, 0xe5, 0x14, 0x42, 0xc7, 0xa7,
	0x89, 0x6b, 0xa0, 0x88, 0x18, 0x62, 0x6c, 0xa4, 0x6d, 0x65, 0xd9, 0x23, 0xfc, 0x20, 0xd9, 0x33,
	0x1c, 0x1a, 0x98, 0x29, 0x71, 0x31, 0x75, 0x28, 0x09, 0x4d, 0x21, 0x6c, 0xa1, 0x88, 0x98, 0xc8,
	0x71, 0x68, 0x12, 0x72, 0xb1, 0x2a, 0xab, 0x87, 0xdb, 0x94, 0x86, 0x47, 0xa9, 0xe7, 0x63, 0x53,
	0x74, 0x7b, 0xc9, 0xbe, 0xc9, 0x49, 0x80, 0x63, 0x8e, 0x82, 0x28, 0x13, 0xd4, 0x33, 0x81, 0x58,
	0x13, 0x86, 0x94, 0x23, 0x4e, 0x68, 0x18, 0x67, 0xd3, 0xe7, 0xe2, 0xc7, 0x69, 0x79, 0x38, 0x6c,
	0xc5, 0x47, 0xc8, 0xf3, 0x30, 0x33, 0x69, 0x24, 0x14, 0x63, 0xd4, 0xad, 0x6b, 0xd7, 0xf4, 0xa8,
	0x47, 0xaf, 0x4e, 0x1d, 0x74, 0xa2, 0x11, 0xd5, 0x50, 0xae, 0xff, 0x99, 0x04, 0xc5, 0xdd, 0x18,
	0x33, 0xf8, 0x12, 0x14, 0x88, 0x2b, 0x4b, 0x9a, 0xd4, 0xac, 0xac, 0x3d, 0xed, 0xf7, 0x1a, 0x8f,
	0x3d, 0xca, 0x82, 0x15, 0x9d, 0x1f, 0x47, 0x78, 0x25, 0x45, 0xcc, 0x39, 0x40, 0xac, 0xb9, 0xf8,
	0x62, 0x7e, 0x35, 0x62, 0x24, 0x40, 0xec, 0xb8, 0x7b, 0x88, 0x8f, 0x75, 0xbb, 0x40, 0x5c, 0xd8,
	0x02, 0x25, 0x1c, 0x20, 0xe2, 0xcb, 0x05, 0xe1, 0x7d, 0xd8, 0xef, 0x35, 0xe6, 0x86, 0xde, 0x24,
	0x24, 0x1f, 0x13, 0xdc, 0x25, 0xa1, 0x8b, 0x3f, 0xe9, 0xf6, 0x50, 0x05, 0x97, 0x41, 0x39, 0x42,
	0x71, 0x7c, 0x44, 0x99, 0x2b, 0x4f, 0x0a, 0x47, 0xbd, 0xdf, 0x6b, 0xc8, 0x63, 0x4e, 0x6b, 0x2f,
	0x2c, 0xcc, 0xeb, 0xf6, 0x48, 0x0d, 0x17, 0x40, 0x31, 0x44, 0x01, 0x96, 0x8b, 0xb7, 0x70, 0x09,
	0xe5, 0xc0, 0xc1, 0xa8, 0x8f, 0xe5, 0x92, 0x26, 0x35, 0x67, 0x3b, 0x75, 0xe3, 0x7f, 0xaa, 0xc6,
	0xe0, 0xed, 0x36, 0xf5, 0xb1, 0x2d, 0x94, 0xf0, 0x11, 0xa8, 0x90, 0xb8, 0x8b, 0x1c, 0x4e, 0x52,
	0x2c, 0x4f, 0x69, 0x52, 0xb3, 0x6c, 0x97, 0x49, 0x6c, 0x89, 0x1e, 0xbe, 0x01, 0x77, 0xc4, 0x04,
	0x71, 0xec, 0x76, 0x11, 0x97, 0xa7, 0x35, 0xa9, 0x59, 0xed, 0x28, 0xc6, 0x90, 0x9e, 0x91, 0x07,
	0x6d, 0xbc, 0xcf, 0xf1, 0xae, 0x15, 0xbf, 0xfc, 0x6a, 0x48, 0x76, 0x75, 0xe4, 0xb2, 0x38, 0x7c,
	0x0d, 0x80, 0xc3, 0x70, 0xbe, 0xa2, 0x7c, 0xcb, 0x15, 0x95, 0xcc, 0x63, 0x71, 0xb8, 0x0c, 0x4a,
	0x9c, 0x1e, 0xe2, 0x50, 0x9e, 0x11, 0x39, 0xe8, 0xfd, 0x5e, 0x43, 0x1d, 0x93, 0x43, 0x67, 0x69,
	0x69, 0x7e, 0x35, 0x8f, 0x5e, 0x18, 0xf4, 0x33, 0x09, 0x54, 0x07, 0xef, 0xdd, 0x62, 0x74, 0x9f,
	0xf8, 0x18, 0xce, 0x5e, 0x21, 0x17, 0x24, 0xef, 0xdf, 0x20, 0x99, 0x03, 0x83, 0x59, 0xec, 0x02,
	0x56, 0x16, 0xec, 0x8d, 0x98, 0x8a, 0xff, 0xc4, 0xf4, 0x0a, 0x4c, 0x67, 0xff, 0x7f, 0x11, 0x7c,
	0xb5, 0xf3, 0xe4, 0x5a, 0xf0, 0xf9, 0x97, 0x91, 0xb6, 0x0d, 0x6b, 0x58, 0x66, 0xb7, 0xb1, 0x73,
	0xd3, 0xb3, 0xaf, 0x12, 0x28, 0xe7, 0x58, 0xa0, 0x06, 0xa6, 0xed, 0xf5, 0x8d, 0xdd, 0xb7, 0x96,
	0x5d, 0x9b, 0x50, 0xe6, 0x4e, 0x4e, 0xb5, 0xbb, 0x23, 0x62, 0xd8, 0x4b, 0x7c, 0xc4, 0xe0, 0x03,
	0x50, 0xd8, 0xb6, 0x6a, 0x93, 0xca, 0xec, 0xc9, 0xa9, 0x06, 0xf2, 0xe1, 0x36, 0x1a, 0x38, 0x37,
	0xad, 0x77, 0xd6, 0xc6, 0xba, 0x5d, 0x9b, 0xba, 0xe9, 0xdc, 0x44, 0x21, 0xf2, 0x30, 0x83, 0x75,
	0x50, 0xda, 0xd9, 0xdd, 0x5a, 0xb7, 0x6b, 0x15, 0xe5, 0xde, 0xc9, 0xa9, 0x36, 0x93, 0xcf, 0x77,
	0x92, 0x08, 0x33, 0xa5, 0xf6, 0xf9, 0x9b, 0x3a, 0xf1, 0xfd, 0x4c, 0x1d, 0xdd, 0x65, 0x4d, 0x3e,
	0xbf, 0x50, 0xa5, 0x9f, 0x17, 0xaa, 0xf4, 0xfb, 0x42, 0x95, 0x7e, 0x5c, 0xaa, 0xd2, 0xf9, 0xa5,
	0x2a, 0x7d, 0x28, 0xa4, 0xed, 0xbd, 0x29, 0x01, 0x6e, 0xf1, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x18, 0x34, 0x7c, 0xaf, 0x4d, 0x04, 0x00, 0x00,
}

func (m *User) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *User) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Id)))
		i += copy(dAtA[i:], m.Id)
	}
	if len(m.Email) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Email)))
		i += copy(dAtA[i:], m.Email)
	}
	if len(m.Password) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Password)))
		i += copy(dAtA[i:], m.Password)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.Role != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintUser(dAtA, i, uint64(m.Role))
	}
	if m.IsActive {
		dAtA[i] = 0x30
		i++
		if m.IsActive {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.ActivatedAt != nil {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.ActivatedAt)))
		n1, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.ActivatedAt, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.CreatedAt != nil {
		dAtA[i] = 0x42
		i++
		i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.CreatedAt)))
		n2, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.CreatedAt, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if len(m.Token) > 0 {
		dAtA[i] = 0x6a
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Token)))
		i += copy(dAtA[i:], m.Token)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *UserProfile) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UserProfile) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Id)))
		i += copy(dAtA[i:], m.Id)
	}
	if len(m.Email) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Email)))
		i += copy(dAtA[i:], m.Email)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.IsActive {
		dAtA[i] = 0x20
		i++
		if m.IsActive {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.Account != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintUser(dAtA, i, uint64(m.Account.Size()))
		n3, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintUser(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *User) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Password)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	if m.Role != 0 {
		n += 1 + sovUser(uint64(m.Role))
	}
	if m.IsActive {
		n += 2
	}
	if m.ActivatedAt != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.ActivatedAt)
		n += 1 + l + sovUser(uint64(l))
	}
	if m.CreatedAt != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.CreatedAt)
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *UserProfile) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	if m.IsActive {
		n += 2
	}
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovUser(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovUser(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozUser(x uint64) (n int) {
	return sovUser(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *User) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUser
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
			return fmt.Errorf("proto: User: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: User: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Email", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Email = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Password", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Password = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Role", wireType)
			}
			m.Role = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Role |= (UserRole(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsActive", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsActive = bool(v != 0)
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActivatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ActivatedAt == nil {
				m.ActivatedAt = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.ActivatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CreatedAt == nil {
				m.CreatedAt = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.CreatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUser(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthUser
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
func (m *UserProfile) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUser
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
			return fmt.Errorf("proto: UserProfile: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UserProfile: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Email", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Email = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsActive", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsActive = bool(v != 0)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Account == nil {
				m.Account = &v1.AccountProfile{}
			}
			if err := m.Account.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUser(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthUser
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
func skipUser(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUser
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
					return 0, ErrIntOverflowUser
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
					return 0, ErrIntOverflowUser
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
				return 0, ErrInvalidLengthUser
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowUser
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
				next, err := skipUser(dAtA[start:])
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
	ErrInvalidLengthUser = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUser   = fmt.Errorf("proto: integer overflow")
)

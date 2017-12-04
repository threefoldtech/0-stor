// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ztor.proto

/*
Package schema is a generated protocol buffer package.

It is generated from these files:
	ztor.proto

It has these top-level messages:
	Empty
	Namespace
	Object
	GetNamespaceRequest
	GetNamespaceReply
	ListObjectsRequest
	CreateObjectRequest
	CreateObjectReply
	ExistsObjectRequest
	ExistsObjectReply
	GetObjectRequest
	GetObjectReply
	DeleteObjectRequest
	DeleteObjectReply
	UpdateReferenceListRequest
	UpdateReferenceListReply
	CheckRequest
	CheckResponse
*/
package schema

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CheckResponse_Status int32

const (
	CheckResponse_ok        CheckResponse_Status = 0
	CheckResponse_corrupted CheckResponse_Status = 1
	CheckResponse_missing   CheckResponse_Status = 2
)

var CheckResponse_Status_name = map[int32]string{
	0: "ok",
	1: "corrupted",
	2: "missing",
}
var CheckResponse_Status_value = map[string]int32{
	"ok":        0,
	"corrupted": 1,
	"missing":   2,
}

func (x CheckResponse_Status) String() string {
	return proto.EnumName(CheckResponse_Status_name, int32(x))
}
func (CheckResponse_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{17, 0} }

// Types
type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Namespace struct {
	Label               string `protobuf:"bytes,1,opt,name=label" json:"label,omitempty"`
	SpaceAvailable      int64  `protobuf:"varint,2,opt,name=spaceAvailable" json:"spaceAvailable,omitempty"`
	SpaceUsed           int64  `protobuf:"varint,3,opt,name=spaceUsed" json:"spaceUsed,omitempty"`
	ReadRequestPerHour  int64  `protobuf:"varint,4,opt,name=readRequestPerHour" json:"readRequestPerHour,omitempty"`
	WriteRequestPerHour int64  `protobuf:"varint,5,opt,name=writeRequestPerHour" json:"writeRequestPerHour,omitempty"`
	NrObjects           int64  `protobuf:"varint,6,opt,name=nrObjects" json:"nrObjects,omitempty"`
}

func (m *Namespace) Reset()                    { *m = Namespace{} }
func (m *Namespace) String() string            { return proto.CompactTextString(m) }
func (*Namespace) ProtoMessage()               {}
func (*Namespace) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Namespace) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *Namespace) GetSpaceAvailable() int64 {
	if m != nil {
		return m.SpaceAvailable
	}
	return 0
}

func (m *Namespace) GetSpaceUsed() int64 {
	if m != nil {
		return m.SpaceUsed
	}
	return 0
}

func (m *Namespace) GetReadRequestPerHour() int64 {
	if m != nil {
		return m.ReadRequestPerHour
	}
	return 0
}

func (m *Namespace) GetWriteRequestPerHour() int64 {
	if m != nil {
		return m.WriteRequestPerHour
	}
	return 0
}

func (m *Namespace) GetNrObjects() int64 {
	if m != nil {
		return m.NrObjects
	}
	return 0
}

type Object struct {
	Key           []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	ReferenceList []string `protobuf:"bytes,3,rep,name=referenceList" json:"referenceList,omitempty"`
}

func (m *Object) Reset()                    { *m = Object{} }
func (m *Object) String() string            { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()               {}
func (*Object) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Object) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Object) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Object) GetReferenceList() []string {
	if m != nil {
		return m.ReferenceList
	}
	return nil
}

// Namespace management
type GetNamespaceRequest struct {
	Label string `protobuf:"bytes,1,opt,name=label" json:"label,omitempty"`
}

func (m *GetNamespaceRequest) Reset()                    { *m = GetNamespaceRequest{} }
func (m *GetNamespaceRequest) String() string            { return proto.CompactTextString(m) }
func (*GetNamespaceRequest) ProtoMessage()               {}
func (*GetNamespaceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetNamespaceRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

type GetNamespaceReply struct {
	Namespace *Namespace `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *GetNamespaceReply) Reset()                    { *m = GetNamespaceReply{} }
func (m *GetNamespaceReply) String() string            { return proto.CompactTextString(m) }
func (*GetNamespaceReply) ProtoMessage()               {}
func (*GetNamespaceReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetNamespaceReply) GetNamespace() *Namespace {
	if m != nil {
		return m.Namespace
	}
	return nil
}

// Object management
type ListObjectsRequest struct {
	Label string `protobuf:"bytes,1,opt,name=label" json:"label,omitempty"`
}

func (m *ListObjectsRequest) Reset()                    { *m = ListObjectsRequest{} }
func (m *ListObjectsRequest) String() string            { return proto.CompactTextString(m) }
func (*ListObjectsRequest) ProtoMessage()               {}
func (*ListObjectsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ListObjectsRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

type CreateObjectRequest struct {
	Label  string  `protobuf:"bytes,1,opt,name=Label" json:"Label,omitempty"`
	Object *Object `protobuf:"bytes,2,opt,name=object" json:"object,omitempty"`
}

func (m *CreateObjectRequest) Reset()                    { *m = CreateObjectRequest{} }
func (m *CreateObjectRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateObjectRequest) ProtoMessage()               {}
func (*CreateObjectRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CreateObjectRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *CreateObjectRequest) GetObject() *Object {
	if m != nil {
		return m.Object
	}
	return nil
}

type CreateObjectReply struct {
}

func (m *CreateObjectReply) Reset()                    { *m = CreateObjectReply{} }
func (m *CreateObjectReply) String() string            { return proto.CompactTextString(m) }
func (*CreateObjectReply) ProtoMessage()               {}
func (*CreateObjectReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type ExistsObjectRequest struct {
	Label string `protobuf:"bytes,1,opt,name=Label" json:"Label,omitempty"`
	Key   []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *ExistsObjectRequest) Reset()                    { *m = ExistsObjectRequest{} }
func (m *ExistsObjectRequest) String() string            { return proto.CompactTextString(m) }
func (*ExistsObjectRequest) ProtoMessage()               {}
func (*ExistsObjectRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ExistsObjectRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *ExistsObjectRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

type ExistsObjectReply struct {
	Exists bool `protobuf:"varint,1,opt,name=exists" json:"exists,omitempty"`
}

func (m *ExistsObjectReply) Reset()                    { *m = ExistsObjectReply{} }
func (m *ExistsObjectReply) String() string            { return proto.CompactTextString(m) }
func (*ExistsObjectReply) ProtoMessage()               {}
func (*ExistsObjectReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ExistsObjectReply) GetExists() bool {
	if m != nil {
		return m.Exists
	}
	return false
}

type GetObjectRequest struct {
	Label string `protobuf:"bytes,1,opt,name=Label" json:"Label,omitempty"`
	Key   []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *GetObjectRequest) Reset()                    { *m = GetObjectRequest{} }
func (m *GetObjectRequest) String() string            { return proto.CompactTextString(m) }
func (*GetObjectRequest) ProtoMessage()               {}
func (*GetObjectRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *GetObjectRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *GetObjectRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

type GetObjectReply struct {
	Object *Object `protobuf:"bytes,2,opt,name=object" json:"object,omitempty"`
}

func (m *GetObjectReply) Reset()                    { *m = GetObjectReply{} }
func (m *GetObjectReply) String() string            { return proto.CompactTextString(m) }
func (*GetObjectReply) ProtoMessage()               {}
func (*GetObjectReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *GetObjectReply) GetObject() *Object {
	if m != nil {
		return m.Object
	}
	return nil
}

type DeleteObjectRequest struct {
	Label string `protobuf:"bytes,1,opt,name=Label" json:"Label,omitempty"`
	Key   []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *DeleteObjectRequest) Reset()                    { *m = DeleteObjectRequest{} }
func (m *DeleteObjectRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteObjectRequest) ProtoMessage()               {}
func (*DeleteObjectRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *DeleteObjectRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *DeleteObjectRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

type DeleteObjectReply struct {
}

func (m *DeleteObjectReply) Reset()                    { *m = DeleteObjectReply{} }
func (m *DeleteObjectReply) String() string            { return proto.CompactTextString(m) }
func (*DeleteObjectReply) ProtoMessage()               {}
func (*DeleteObjectReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type UpdateReferenceListRequest struct {
	Label         string   `protobuf:"bytes,1,opt,name=label" json:"label,omitempty"`
	Key           []byte   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	ReferenceList []string `protobuf:"bytes,3,rep,name=referenceList" json:"referenceList,omitempty"`
}

func (m *UpdateReferenceListRequest) Reset()                    { *m = UpdateReferenceListRequest{} }
func (m *UpdateReferenceListRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateReferenceListRequest) ProtoMessage()               {}
func (*UpdateReferenceListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *UpdateReferenceListRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *UpdateReferenceListRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *UpdateReferenceListRequest) GetReferenceList() []string {
	if m != nil {
		return m.ReferenceList
	}
	return nil
}

type UpdateReferenceListReply struct {
}

func (m *UpdateReferenceListReply) Reset()                    { *m = UpdateReferenceListReply{} }
func (m *UpdateReferenceListReply) String() string            { return proto.CompactTextString(m) }
func (*UpdateReferenceListReply) ProtoMessage()               {}
func (*UpdateReferenceListReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type CheckRequest struct {
	Label string   `protobuf:"bytes,1,opt,name=label" json:"label,omitempty"`
	Ids   []string `protobuf:"bytes,2,rep,name=ids" json:"ids,omitempty"`
}

func (m *CheckRequest) Reset()                    { *m = CheckRequest{} }
func (m *CheckRequest) String() string            { return proto.CompactTextString(m) }
func (*CheckRequest) ProtoMessage()               {}
func (*CheckRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *CheckRequest) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *CheckRequest) GetIds() []string {
	if m != nil {
		return m.Ids
	}
	return nil
}

type CheckResponse struct {
	Id     string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Status CheckResponse_Status `protobuf:"varint,2,opt,name=status,enum=CheckResponse_Status" json:"status,omitempty"`
}

func (m *CheckResponse) Reset()                    { *m = CheckResponse{} }
func (m *CheckResponse) String() string            { return proto.CompactTextString(m) }
func (*CheckResponse) ProtoMessage()               {}
func (*CheckResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *CheckResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CheckResponse) GetStatus() CheckResponse_Status {
	if m != nil {
		return m.Status
	}
	return CheckResponse_ok
}

func init() {
	proto.RegisterType((*Empty)(nil), "Empty")
	proto.RegisterType((*Namespace)(nil), "Namespace")
	proto.RegisterType((*Object)(nil), "Object")
	proto.RegisterType((*GetNamespaceRequest)(nil), "GetNamespaceRequest")
	proto.RegisterType((*GetNamespaceReply)(nil), "GetNamespaceReply")
	proto.RegisterType((*ListObjectsRequest)(nil), "ListObjectsRequest")
	proto.RegisterType((*CreateObjectRequest)(nil), "CreateObjectRequest")
	proto.RegisterType((*CreateObjectReply)(nil), "CreateObjectReply")
	proto.RegisterType((*ExistsObjectRequest)(nil), "ExistsObjectRequest")
	proto.RegisterType((*ExistsObjectReply)(nil), "ExistsObjectReply")
	proto.RegisterType((*GetObjectRequest)(nil), "GetObjectRequest")
	proto.RegisterType((*GetObjectReply)(nil), "GetObjectReply")
	proto.RegisterType((*DeleteObjectRequest)(nil), "DeleteObjectRequest")
	proto.RegisterType((*DeleteObjectReply)(nil), "DeleteObjectReply")
	proto.RegisterType((*UpdateReferenceListRequest)(nil), "UpdateReferenceListRequest")
	proto.RegisterType((*UpdateReferenceListReply)(nil), "UpdateReferenceListReply")
	proto.RegisterType((*CheckRequest)(nil), "CheckRequest")
	proto.RegisterType((*CheckResponse)(nil), "CheckResponse")
	proto.RegisterEnum("CheckResponse_Status", CheckResponse_Status_name, CheckResponse_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for NamespaceManager service

type NamespaceManagerClient interface {
	Get(ctx context.Context, in *GetNamespaceRequest, opts ...grpc.CallOption) (*GetNamespaceReply, error)
}

type namespaceManagerClient struct {
	cc *grpc.ClientConn
}

func NewNamespaceManagerClient(cc *grpc.ClientConn) NamespaceManagerClient {
	return &namespaceManagerClient{cc}
}

func (c *namespaceManagerClient) Get(ctx context.Context, in *GetNamespaceRequest, opts ...grpc.CallOption) (*GetNamespaceReply, error) {
	out := new(GetNamespaceReply)
	err := grpc.Invoke(ctx, "/NamespaceManager/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for NamespaceManager service

type NamespaceManagerServer interface {
	Get(context.Context, *GetNamespaceRequest) (*GetNamespaceReply, error)
}

func RegisterNamespaceManagerServer(s *grpc.Server, srv NamespaceManagerServer) {
	s.RegisterService(&_NamespaceManager_serviceDesc, srv)
}

func _NamespaceManager_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NamespaceManagerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NamespaceManager/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NamespaceManagerServer).Get(ctx, req.(*GetNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NamespaceManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "NamespaceManager",
	HandlerType: (*NamespaceManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _NamespaceManager_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ztor.proto",
}

// Client API for ObjectManager service

type ObjectManagerClient interface {
	Create(ctx context.Context, in *CreateObjectRequest, opts ...grpc.CallOption) (*CreateObjectReply, error)
	List(ctx context.Context, in *ListObjectsRequest, opts ...grpc.CallOption) (ObjectManager_ListClient, error)
	Get(ctx context.Context, in *GetObjectRequest, opts ...grpc.CallOption) (*GetObjectReply, error)
	Exists(ctx context.Context, in *ExistsObjectRequest, opts ...grpc.CallOption) (*ExistsObjectReply, error)
	Delete(ctx context.Context, in *DeleteObjectRequest, opts ...grpc.CallOption) (*DeleteObjectReply, error)
	SetReferenceList(ctx context.Context, in *UpdateReferenceListRequest, opts ...grpc.CallOption) (*UpdateReferenceListReply, error)
	AppendReferenceList(ctx context.Context, in *UpdateReferenceListRequest, opts ...grpc.CallOption) (*UpdateReferenceListReply, error)
	RemoveReferenceList(ctx context.Context, in *UpdateReferenceListRequest, opts ...grpc.CallOption) (*UpdateReferenceListReply, error)
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (ObjectManager_CheckClient, error)
}

type objectManagerClient struct {
	cc *grpc.ClientConn
}

func NewObjectManagerClient(cc *grpc.ClientConn) ObjectManagerClient {
	return &objectManagerClient{cc}
}

func (c *objectManagerClient) Create(ctx context.Context, in *CreateObjectRequest, opts ...grpc.CallOption) (*CreateObjectReply, error) {
	out := new(CreateObjectReply)
	err := grpc.Invoke(ctx, "/ObjectManager/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectManagerClient) List(ctx context.Context, in *ListObjectsRequest, opts ...grpc.CallOption) (ObjectManager_ListClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ObjectManager_serviceDesc.Streams[0], c.cc, "/ObjectManager/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &objectManagerListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ObjectManager_ListClient interface {
	Recv() (*Object, error)
	grpc.ClientStream
}

type objectManagerListClient struct {
	grpc.ClientStream
}

func (x *objectManagerListClient) Recv() (*Object, error) {
	m := new(Object)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *objectManagerClient) Get(ctx context.Context, in *GetObjectRequest, opts ...grpc.CallOption) (*GetObjectReply, error) {
	out := new(GetObjectReply)
	err := grpc.Invoke(ctx, "/ObjectManager/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectManagerClient) Exists(ctx context.Context, in *ExistsObjectRequest, opts ...grpc.CallOption) (*ExistsObjectReply, error) {
	out := new(ExistsObjectReply)
	err := grpc.Invoke(ctx, "/ObjectManager/Exists", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectManagerClient) Delete(ctx context.Context, in *DeleteObjectRequest, opts ...grpc.CallOption) (*DeleteObjectReply, error) {
	out := new(DeleteObjectReply)
	err := grpc.Invoke(ctx, "/ObjectManager/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectManagerClient) SetReferenceList(ctx context.Context, in *UpdateReferenceListRequest, opts ...grpc.CallOption) (*UpdateReferenceListReply, error) {
	out := new(UpdateReferenceListReply)
	err := grpc.Invoke(ctx, "/ObjectManager/SetReferenceList", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectManagerClient) AppendReferenceList(ctx context.Context, in *UpdateReferenceListRequest, opts ...grpc.CallOption) (*UpdateReferenceListReply, error) {
	out := new(UpdateReferenceListReply)
	err := grpc.Invoke(ctx, "/ObjectManager/AppendReferenceList", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectManagerClient) RemoveReferenceList(ctx context.Context, in *UpdateReferenceListRequest, opts ...grpc.CallOption) (*UpdateReferenceListReply, error) {
	out := new(UpdateReferenceListReply)
	err := grpc.Invoke(ctx, "/ObjectManager/RemoveReferenceList", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectManagerClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (ObjectManager_CheckClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ObjectManager_serviceDesc.Streams[1], c.cc, "/ObjectManager/Check", opts...)
	if err != nil {
		return nil, err
	}
	x := &objectManagerCheckClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ObjectManager_CheckClient interface {
	Recv() (*CheckResponse, error)
	grpc.ClientStream
}

type objectManagerCheckClient struct {
	grpc.ClientStream
}

func (x *objectManagerCheckClient) Recv() (*CheckResponse, error) {
	m := new(CheckResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for ObjectManager service

type ObjectManagerServer interface {
	Create(context.Context, *CreateObjectRequest) (*CreateObjectReply, error)
	List(*ListObjectsRequest, ObjectManager_ListServer) error
	Get(context.Context, *GetObjectRequest) (*GetObjectReply, error)
	Exists(context.Context, *ExistsObjectRequest) (*ExistsObjectReply, error)
	Delete(context.Context, *DeleteObjectRequest) (*DeleteObjectReply, error)
	SetReferenceList(context.Context, *UpdateReferenceListRequest) (*UpdateReferenceListReply, error)
	AppendReferenceList(context.Context, *UpdateReferenceListRequest) (*UpdateReferenceListReply, error)
	RemoveReferenceList(context.Context, *UpdateReferenceListRequest) (*UpdateReferenceListReply, error)
	Check(*CheckRequest, ObjectManager_CheckServer) error
}

func RegisterObjectManagerServer(s *grpc.Server, srv ObjectManagerServer) {
	s.RegisterService(&_ObjectManager_serviceDesc, srv)
}

func _ObjectManager_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectManagerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ObjectManager/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectManagerServer).Create(ctx, req.(*CreateObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectManager_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListObjectsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ObjectManagerServer).List(m, &objectManagerListServer{stream})
}

type ObjectManager_ListServer interface {
	Send(*Object) error
	grpc.ServerStream
}

type objectManagerListServer struct {
	grpc.ServerStream
}

func (x *objectManagerListServer) Send(m *Object) error {
	return x.ServerStream.SendMsg(m)
}

func _ObjectManager_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectManagerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ObjectManager/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectManagerServer).Get(ctx, req.(*GetObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectManager_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistsObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectManagerServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ObjectManager/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectManagerServer).Exists(ctx, req.(*ExistsObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectManager_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectManagerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ObjectManager/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectManagerServer).Delete(ctx, req.(*DeleteObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectManager_SetReferenceList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateReferenceListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectManagerServer).SetReferenceList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ObjectManager/SetReferenceList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectManagerServer).SetReferenceList(ctx, req.(*UpdateReferenceListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectManager_AppendReferenceList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateReferenceListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectManagerServer).AppendReferenceList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ObjectManager/AppendReferenceList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectManagerServer).AppendReferenceList(ctx, req.(*UpdateReferenceListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectManager_RemoveReferenceList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateReferenceListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectManagerServer).RemoveReferenceList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ObjectManager/RemoveReferenceList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectManagerServer).RemoveReferenceList(ctx, req.(*UpdateReferenceListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectManager_Check_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CheckRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ObjectManagerServer).Check(m, &objectManagerCheckServer{stream})
}

type ObjectManager_CheckServer interface {
	Send(*CheckResponse) error
	grpc.ServerStream
}

type objectManagerCheckServer struct {
	grpc.ServerStream
}

func (x *objectManagerCheckServer) Send(m *CheckResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _ObjectManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ObjectManager",
	HandlerType: (*ObjectManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ObjectManager_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ObjectManager_Get_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _ObjectManager_Exists_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ObjectManager_Delete_Handler,
		},
		{
			MethodName: "SetReferenceList",
			Handler:    _ObjectManager_SetReferenceList_Handler,
		},
		{
			MethodName: "AppendReferenceList",
			Handler:    _ObjectManager_AppendReferenceList_Handler,
		},
		{
			MethodName: "RemoveReferenceList",
			Handler:    _ObjectManager_RemoveReferenceList_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _ObjectManager_List_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Check",
			Handler:       _ObjectManager_Check_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ztor.proto",
}

func init() { proto.RegisterFile("ztor.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 670 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x5d, 0x6f, 0xd3, 0x4a,
	0x10, 0x75, 0x9c, 0xc6, 0xbd, 0x99, 0x36, 0xb9, 0xce, 0x38, 0xf7, 0xca, 0xd7, 0x17, 0x89, 0x6a,
	0x85, 0x50, 0xd4, 0xc2, 0xaa, 0x0d, 0x88, 0x07, 0xa4, 0x3e, 0x54, 0xa5, 0x2a, 0x0f, 0x81, 0x22,
	0x57, 0xe5, 0xdd, 0x89, 0x87, 0xe2, 0xd6, 0xb1, 0x8d, 0x77, 0x53, 0x08, 0xe2, 0xe7, 0xf0, 0xdf,
	0xf8, 0x1b, 0xc8, 0x6b, 0x37, 0x8d, 0x6b, 0xf7, 0x43, 0xa8, 0x6f, 0xde, 0xd9, 0x73, 0xe6, 0xcc,
	0xcc, 0x9e, 0x91, 0x01, 0xbe, 0xcb, 0x38, 0xe5, 0x49, 0x1a, 0xcb, 0x98, 0xad, 0x42, 0xeb, 0x60,
	0x9a, 0xc8, 0x39, 0xfb, 0xd5, 0x80, 0xf6, 0x7b, 0x6f, 0x4a, 0x22, 0xf1, 0x26, 0x84, 0x7d, 0x68,
	0x85, 0xde, 0x98, 0x42, 0xbb, 0xb1, 0xd1, 0x18, 0xb4, 0xdd, 0xfc, 0x80, 0x4f, 0xa1, 0xab, 0xae,
	0xf7, 0x2e, 0xbc, 0x20, 0xf4, 0xc6, 0x21, 0xd9, 0xfa, 0x46, 0x63, 0xd0, 0x74, 0xaf, 0x45, 0xf1,
	0x11, 0xb4, 0x55, 0xe4, 0x44, 0x90, 0x6f, 0x37, 0x15, 0xe4, 0x2a, 0x80, 0x1c, 0x30, 0x25, 0xcf,
	0x77, 0xe9, 0xcb, 0x8c, 0x84, 0xfc, 0x40, 0xe9, 0xdb, 0x78, 0x96, 0xda, 0x2b, 0x0a, 0x56, 0x73,
	0x83, 0xdb, 0x60, 0x7d, 0x4d, 0x03, 0x49, 0xd7, 0x08, 0x2d, 0x45, 0xa8, 0xbb, 0xca, 0xf4, 0xa3,
	0xf4, 0x68, 0x7c, 0x46, 0x13, 0x29, 0x6c, 0x23, 0xd7, 0x5f, 0x04, 0xd8, 0x47, 0x30, 0xf2, 0x4f,
	0x34, 0xa1, 0x79, 0x4e, 0x73, 0xd5, 0xe3, 0xba, 0x9b, 0x7d, 0x66, 0x7d, 0x5f, 0x78, 0xe1, 0x2c,
	0x6f, 0x6c, 0xdd, 0xcd, 0x0f, 0xf8, 0x04, 0x3a, 0x29, 0x7d, 0xa2, 0x94, 0xa2, 0x09, 0x8d, 0x02,
	0x21, 0xed, 0xe6, 0x46, 0x73, 0xd0, 0x76, 0xcb, 0x41, 0xb6, 0x05, 0xd6, 0x21, 0xc9, 0xc5, 0x0c,
	0x8b, 0x9a, 0xea, 0x47, 0xc9, 0x76, 0xa1, 0x57, 0x06, 0x27, 0xe1, 0x1c, 0x07, 0xd0, 0x8e, 0x2e,
	0x23, 0x0a, 0xbe, 0x36, 0x04, 0x7e, 0x85, 0xb9, 0xba, 0x64, 0x9b, 0x80, 0x99, 0x66, 0xd1, 0xd2,
	0xed, 0x52, 0x23, 0xb0, 0xf6, 0x53, 0xf2, 0x24, 0xe5, 0xe8, 0x25, 0xf0, 0x68, 0x19, 0xac, 0x0e,
	0xf8, 0x18, 0x8c, 0x58, 0xc1, 0xd4, 0x04, 0xd6, 0x86, 0xab, 0xbc, 0x60, 0x15, 0x61, 0x66, 0x41,
	0xaf, 0x9c, 0x2d, 0x09, 0xe7, 0x6c, 0x17, 0xac, 0x83, 0x6f, 0x81, 0x90, 0xe2, 0x3e, 0x12, 0xc5,
	0xd4, 0xf5, 0xc5, 0xd4, 0xd9, 0x16, 0xf4, 0xca, 0xf4, 0x6c, 0x18, 0xff, 0x82, 0x41, 0x2a, 0xa8,
	0xd8, 0x7f, 0xb9, 0xc5, 0x89, 0xbd, 0x06, 0xf3, 0x90, 0xe4, 0x9f, 0x09, 0xed, 0x40, 0x77, 0x89,
	0x9b, 0xa9, 0xdc, 0xd9, 0xef, 0x2e, 0x58, 0x6f, 0x28, 0xa4, 0xfb, 0x4d, 0xaf, 0xaa, 0x68, 0x41,
	0xaf, 0x4c, 0xcf, 0xc6, 0x75, 0x06, 0xce, 0x49, 0xe2, 0x7b, 0x99, 0x6f, 0x97, 0x0c, 0x74, 0xeb,
	0x2b, 0x56, 0x53, 0xdf, 0xd3, 0x95, 0x0e, 0xd8, 0xb5, 0x5a, 0x59, 0x1d, 0xaf, 0x60, 0x7d, 0xff,
	0x33, 0x4d, 0xce, 0xef, 0x54, 0x0e, 0x7c, 0x61, 0xeb, 0x2a, 0x7b, 0xf6, 0xc9, 0x7e, 0x40, 0xa7,
	0xe0, 0x89, 0x24, 0x8e, 0x04, 0x61, 0x17, 0xf4, 0xc0, 0x2f, 0x58, 0x7a, 0xe0, 0xe3, 0x73, 0x30,
	0x84, 0xf4, 0xe4, 0x4c, 0xa8, 0x7a, 0xbb, 0xc3, 0x7f, 0x78, 0x09, 0xcf, 0x8f, 0xd5, 0xa5, 0x5b,
	0x80, 0xd8, 0x33, 0x30, 0xf2, 0x08, 0x1a, 0xa0, 0xc7, 0xe7, 0xa6, 0x86, 0x1d, 0x68, 0x4f, 0xe2,
	0x34, 0x9d, 0x25, 0x92, 0x7c, 0xb3, 0x81, 0x6b, 0xb0, 0x3a, 0x0d, 0x84, 0x08, 0xa2, 0x53, 0x53,
	0x1f, 0x1e, 0x80, 0xb9, 0xd8, 0x89, 0x77, 0x5e, 0xe4, 0x9d, 0x52, 0x8a, 0x3b, 0xd0, 0x3c, 0x24,
	0x89, 0x7d, 0x5e, 0xb3, 0x81, 0x0e, 0xf2, 0xca, 0xaa, 0x31, 0x6d, 0xf8, 0x73, 0x05, 0x3a, 0xf9,
	0xa3, 0x5c, 0x26, 0x79, 0x09, 0x46, 0x6e, 0x6d, 0xec, 0xf3, 0x9a, 0x8d, 0x71, 0x90, 0x57, 0x9d,
	0xaf, 0xe1, 0x00, 0x56, 0xb2, 0x89, 0xa2, 0xc5, 0xab, 0x1b, 0xe9, 0x5c, 0xda, 0x89, 0x69, 0xdb,
	0x0d, 0xdc, 0xca, 0x8b, 0xec, 0xf1, 0xeb, 0xfe, 0x75, 0xfe, 0xe6, 0x65, 0x5b, 0x32, 0x2d, 0x2b,
	0x26, 0xdf, 0x09, 0xec, 0xf3, 0x9a, 0xdd, 0x72, 0x90, 0x57, 0x56, 0x26, 0x67, 0xe5, 0x76, 0xc3,
	0x3e, 0xaf, 0xb1, 0xad, 0x83, 0xbc, 0xea, 0x46, 0x0d, 0x47, 0x60, 0x1e, 0x93, 0x2c, 0x19, 0x04,
	0xff, 0xe7, 0x37, 0x5b, 0xd4, 0xf9, 0x8f, 0xdf, 0xe8, 0x29, 0x0d, 0x8f, 0xc0, 0xda, 0x4b, 0x12,
	0x8a, 0xfc, 0x07, 0x4c, 0xe8, 0xd2, 0x34, 0xbe, 0xa0, 0x87, 0x4a, 0xb8, 0x09, 0x2d, 0xe5, 0x47,
	0xec, 0xf0, 0x65, 0xff, 0x3b, 0xdd, 0xb2, 0x4d, 0xb3, 0x47, 0x1b, 0x1b, 0xea, 0x3f, 0xf9, 0xe2,
	0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x61, 0x42, 0xe8, 0x35, 0x07, 0x00, 0x00,
}

// Copyright 2018 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was autogenerated. Do not edit directly.
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sdk.proto

package sdk

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_sdk_39d49c489951ac60, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

// A GameServer Custom Resource Definition object
// We will only export those resources that make the most
// sense. Can always expand to more as needed.
type GameServer struct {
	ObjectMeta           *GameServer_ObjectMeta `protobuf:"bytes,1,opt,name=object_meta,json=objectMeta" json:"object_meta,omitempty"`
	Spec                 *GameServer_Spec       `protobuf:"bytes,2,opt,name=spec" json:"spec,omitempty"`
	Status               *GameServer_Status     `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *GameServer) Reset()         { *m = GameServer{} }
func (m *GameServer) String() string { return proto.CompactTextString(m) }
func (*GameServer) ProtoMessage()    {}
func (*GameServer) Descriptor() ([]byte, []int) {
	return fileDescriptor_sdk_39d49c489951ac60, []int{1}
}
func (m *GameServer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameServer.Unmarshal(m, b)
}
func (m *GameServer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameServer.Marshal(b, m, deterministic)
}
func (dst *GameServer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameServer.Merge(dst, src)
}
func (m *GameServer) XXX_Size() int {
	return xxx_messageInfo_GameServer.Size(m)
}
func (m *GameServer) XXX_DiscardUnknown() {
	xxx_messageInfo_GameServer.DiscardUnknown(m)
}

var xxx_messageInfo_GameServer proto.InternalMessageInfo

func (m *GameServer) GetObjectMeta() *GameServer_ObjectMeta {
	if m != nil {
		return m.ObjectMeta
	}
	return nil
}

func (m *GameServer) GetSpec() *GameServer_Spec {
	if m != nil {
		return m.Spec
	}
	return nil
}

func (m *GameServer) GetStatus() *GameServer_Status {
	if m != nil {
		return m.Status
	}
	return nil
}

// representation of the K8s ObjectMeta resource
type GameServer_ObjectMeta struct {
	Name            string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Namespace       string `protobuf:"bytes,2,opt,name=namespace" json:"namespace,omitempty"`
	Uid             string `protobuf:"bytes,3,opt,name=uid" json:"uid,omitempty"`
	ResourceVersion string `protobuf:"bytes,4,opt,name=resource_version,json=resourceVersion" json:"resource_version,omitempty"`
	Generation      int64  `protobuf:"varint,5,opt,name=generation" json:"generation,omitempty"`
	// timestamp is in Epoch format, unit: seconds
	CreationTimestamp int64 `protobuf:"varint,6,opt,name=creation_timestamp,json=creationTimestamp" json:"creation_timestamp,omitempty"`
	// optional deletion timestamp in Epoch format, unit: seconds
	DeletionTimestamp    int64             `protobuf:"varint,7,opt,name=deletion_timestamp,json=deletionTimestamp" json:"deletion_timestamp,omitempty"`
	Annotations          map[string]string `protobuf:"bytes,8,rep,name=annotations" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Labels               map[string]string `protobuf:"bytes,9,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *GameServer_ObjectMeta) Reset()         { *m = GameServer_ObjectMeta{} }
func (m *GameServer_ObjectMeta) String() string { return proto.CompactTextString(m) }
func (*GameServer_ObjectMeta) ProtoMessage()    {}
func (*GameServer_ObjectMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_sdk_39d49c489951ac60, []int{1, 0}
}
func (m *GameServer_ObjectMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameServer_ObjectMeta.Unmarshal(m, b)
}
func (m *GameServer_ObjectMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameServer_ObjectMeta.Marshal(b, m, deterministic)
}
func (dst *GameServer_ObjectMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameServer_ObjectMeta.Merge(dst, src)
}
func (m *GameServer_ObjectMeta) XXX_Size() int {
	return xxx_messageInfo_GameServer_ObjectMeta.Size(m)
}
func (m *GameServer_ObjectMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_GameServer_ObjectMeta.DiscardUnknown(m)
}

var xxx_messageInfo_GameServer_ObjectMeta proto.InternalMessageInfo

func (m *GameServer_ObjectMeta) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GameServer_ObjectMeta) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *GameServer_ObjectMeta) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *GameServer_ObjectMeta) GetResourceVersion() string {
	if m != nil {
		return m.ResourceVersion
	}
	return ""
}

func (m *GameServer_ObjectMeta) GetGeneration() int64 {
	if m != nil {
		return m.Generation
	}
	return 0
}

func (m *GameServer_ObjectMeta) GetCreationTimestamp() int64 {
	if m != nil {
		return m.CreationTimestamp
	}
	return 0
}

func (m *GameServer_ObjectMeta) GetDeletionTimestamp() int64 {
	if m != nil {
		return m.DeletionTimestamp
	}
	return 0
}

func (m *GameServer_ObjectMeta) GetAnnotations() map[string]string {
	if m != nil {
		return m.Annotations
	}
	return nil
}

func (m *GameServer_ObjectMeta) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type GameServer_Spec struct {
	Health               *GameServer_Spec_Health `protobuf:"bytes,1,opt,name=health" json:"health,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *GameServer_Spec) Reset()         { *m = GameServer_Spec{} }
func (m *GameServer_Spec) String() string { return proto.CompactTextString(m) }
func (*GameServer_Spec) ProtoMessage()    {}
func (*GameServer_Spec) Descriptor() ([]byte, []int) {
	return fileDescriptor_sdk_39d49c489951ac60, []int{1, 1}
}
func (m *GameServer_Spec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameServer_Spec.Unmarshal(m, b)
}
func (m *GameServer_Spec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameServer_Spec.Marshal(b, m, deterministic)
}
func (dst *GameServer_Spec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameServer_Spec.Merge(dst, src)
}
func (m *GameServer_Spec) XXX_Size() int {
	return xxx_messageInfo_GameServer_Spec.Size(m)
}
func (m *GameServer_Spec) XXX_DiscardUnknown() {
	xxx_messageInfo_GameServer_Spec.DiscardUnknown(m)
}

var xxx_messageInfo_GameServer_Spec proto.InternalMessageInfo

func (m *GameServer_Spec) GetHealth() *GameServer_Spec_Health {
	if m != nil {
		return m.Health
	}
	return nil
}

type GameServer_Spec_Health struct {
	Disabled             bool     `protobuf:"varint,1,opt,name=Disabled" json:"Disabled,omitempty"`
	PeriodSeconds        int32    `protobuf:"varint,2,opt,name=PeriodSeconds" json:"PeriodSeconds,omitempty"`
	FailureThreshold     int32    `protobuf:"varint,3,opt,name=FailureThreshold" json:"FailureThreshold,omitempty"`
	InitialDelaySeconds  int32    `protobuf:"varint,4,opt,name=InitialDelaySeconds" json:"InitialDelaySeconds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GameServer_Spec_Health) Reset()         { *m = GameServer_Spec_Health{} }
func (m *GameServer_Spec_Health) String() string { return proto.CompactTextString(m) }
func (*GameServer_Spec_Health) ProtoMessage()    {}
func (*GameServer_Spec_Health) Descriptor() ([]byte, []int) {
	return fileDescriptor_sdk_39d49c489951ac60, []int{1, 1, 0}
}
func (m *GameServer_Spec_Health) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameServer_Spec_Health.Unmarshal(m, b)
}
func (m *GameServer_Spec_Health) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameServer_Spec_Health.Marshal(b, m, deterministic)
}
func (dst *GameServer_Spec_Health) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameServer_Spec_Health.Merge(dst, src)
}
func (m *GameServer_Spec_Health) XXX_Size() int {
	return xxx_messageInfo_GameServer_Spec_Health.Size(m)
}
func (m *GameServer_Spec_Health) XXX_DiscardUnknown() {
	xxx_messageInfo_GameServer_Spec_Health.DiscardUnknown(m)
}

var xxx_messageInfo_GameServer_Spec_Health proto.InternalMessageInfo

func (m *GameServer_Spec_Health) GetDisabled() bool {
	if m != nil {
		return m.Disabled
	}
	return false
}

func (m *GameServer_Spec_Health) GetPeriodSeconds() int32 {
	if m != nil {
		return m.PeriodSeconds
	}
	return 0
}

func (m *GameServer_Spec_Health) GetFailureThreshold() int32 {
	if m != nil {
		return m.FailureThreshold
	}
	return 0
}

func (m *GameServer_Spec_Health) GetInitialDelaySeconds() int32 {
	if m != nil {
		return m.InitialDelaySeconds
	}
	return 0
}

type GameServer_Status struct {
	State                string                    `protobuf:"bytes,1,opt,name=state" json:"state,omitempty"`
	Address              string                    `protobuf:"bytes,2,opt,name=address" json:"address,omitempty"`
	Ports                []*GameServer_Status_Port `protobuf:"bytes,3,rep,name=ports" json:"ports,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *GameServer_Status) Reset()         { *m = GameServer_Status{} }
func (m *GameServer_Status) String() string { return proto.CompactTextString(m) }
func (*GameServer_Status) ProtoMessage()    {}
func (*GameServer_Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_sdk_39d49c489951ac60, []int{1, 2}
}
func (m *GameServer_Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameServer_Status.Unmarshal(m, b)
}
func (m *GameServer_Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameServer_Status.Marshal(b, m, deterministic)
}
func (dst *GameServer_Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameServer_Status.Merge(dst, src)
}
func (m *GameServer_Status) XXX_Size() int {
	return xxx_messageInfo_GameServer_Status.Size(m)
}
func (m *GameServer_Status) XXX_DiscardUnknown() {
	xxx_messageInfo_GameServer_Status.DiscardUnknown(m)
}

var xxx_messageInfo_GameServer_Status proto.InternalMessageInfo

func (m *GameServer_Status) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *GameServer_Status) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *GameServer_Status) GetPorts() []*GameServer_Status_Port {
	if m != nil {
		return m.Ports
	}
	return nil
}

type GameServer_Status_Port struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Port                 int32    `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GameServer_Status_Port) Reset()         { *m = GameServer_Status_Port{} }
func (m *GameServer_Status_Port) String() string { return proto.CompactTextString(m) }
func (*GameServer_Status_Port) ProtoMessage()    {}
func (*GameServer_Status_Port) Descriptor() ([]byte, []int) {
	return fileDescriptor_sdk_39d49c489951ac60, []int{1, 2, 0}
}
func (m *GameServer_Status_Port) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameServer_Status_Port.Unmarshal(m, b)
}
func (m *GameServer_Status_Port) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameServer_Status_Port.Marshal(b, m, deterministic)
}
func (dst *GameServer_Status_Port) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameServer_Status_Port.Merge(dst, src)
}
func (m *GameServer_Status_Port) XXX_Size() int {
	return xxx_messageInfo_GameServer_Status_Port.Size(m)
}
func (m *GameServer_Status_Port) XXX_DiscardUnknown() {
	xxx_messageInfo_GameServer_Status_Port.DiscardUnknown(m)
}

var xxx_messageInfo_GameServer_Status_Port proto.InternalMessageInfo

func (m *GameServer_Status_Port) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GameServer_Status_Port) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func init() {
	proto.RegisterType((*Empty)(nil), "stable.agones.dev.sdk.Empty")
	proto.RegisterType((*GameServer)(nil), "stable.agones.dev.sdk.GameServer")
	proto.RegisterType((*GameServer_ObjectMeta)(nil), "stable.agones.dev.sdk.GameServer.ObjectMeta")
	proto.RegisterMapType((map[string]string)(nil), "stable.agones.dev.sdk.GameServer.ObjectMeta.AnnotationsEntry")
	proto.RegisterMapType((map[string]string)(nil), "stable.agones.dev.sdk.GameServer.ObjectMeta.LabelsEntry")
	proto.RegisterType((*GameServer_Spec)(nil), "stable.agones.dev.sdk.GameServer.Spec")
	proto.RegisterType((*GameServer_Spec_Health)(nil), "stable.agones.dev.sdk.GameServer.Spec.Health")
	proto.RegisterType((*GameServer_Status)(nil), "stable.agones.dev.sdk.GameServer.Status")
	proto.RegisterType((*GameServer_Status_Port)(nil), "stable.agones.dev.sdk.GameServer.Status.Port")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SDK service

type SDKClient interface {
	// Call when the GameServer is ready
	Ready(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	// Call when the GameServer is shutting down
	Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	// Send a Empty every d Duration to declare that this GameSever is healthy
	Health(ctx context.Context, opts ...grpc.CallOption) (SDK_HealthClient, error)
	// Retrieve the current GameServer data
	GetGameServer(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GameServer, error)
}

type sDKClient struct {
	cc *grpc.ClientConn
}

func NewSDKClient(cc *grpc.ClientConn) SDKClient {
	return &sDKClient{cc}
}

func (c *sDKClient) Ready(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/stable.agones.dev.sdk.SDK/Ready", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sDKClient) Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/stable.agones.dev.sdk.SDK/Shutdown", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sDKClient) Health(ctx context.Context, opts ...grpc.CallOption) (SDK_HealthClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SDK_serviceDesc.Streams[0], c.cc, "/stable.agones.dev.sdk.SDK/Health", opts...)
	if err != nil {
		return nil, err
	}
	x := &sDKHealthClient{stream}
	return x, nil
}

type SDK_HealthClient interface {
	Send(*Empty) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type sDKHealthClient struct {
	grpc.ClientStream
}

func (x *sDKHealthClient) Send(m *Empty) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sDKHealthClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sDKClient) GetGameServer(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GameServer, error) {
	out := new(GameServer)
	err := grpc.Invoke(ctx, "/stable.agones.dev.sdk.SDK/GetGameServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SDK service

type SDKServer interface {
	// Call when the GameServer is ready
	Ready(context.Context, *Empty) (*Empty, error)
	// Call when the GameServer is shutting down
	Shutdown(context.Context, *Empty) (*Empty, error)
	// Send a Empty every d Duration to declare that this GameSever is healthy
	Health(SDK_HealthServer) error
	// Retrieve the current GameServer data
	GetGameServer(context.Context, *Empty) (*GameServer, error)
}

func RegisterSDKServer(s *grpc.Server, srv SDKServer) {
	s.RegisterService(&_SDK_serviceDesc, srv)
}

func _SDK_Ready_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SDKServer).Ready(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stable.agones.dev.sdk.SDK/Ready",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SDKServer).Ready(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SDK_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SDKServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stable.agones.dev.sdk.SDK/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SDKServer).Shutdown(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SDK_Health_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SDKServer).Health(&sDKHealthServer{stream})
}

type SDK_HealthServer interface {
	SendAndClose(*Empty) error
	Recv() (*Empty, error)
	grpc.ServerStream
}

type sDKHealthServer struct {
	grpc.ServerStream
}

func (x *sDKHealthServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sDKHealthServer) Recv() (*Empty, error) {
	m := new(Empty)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SDK_GetGameServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SDKServer).GetGameServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stable.agones.dev.sdk.SDK/GetGameServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SDKServer).GetGameServer(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _SDK_serviceDesc = grpc.ServiceDesc{
	ServiceName: "stable.agones.dev.sdk.SDK",
	HandlerType: (*SDKServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ready",
			Handler:    _SDK_Ready_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _SDK_Shutdown_Handler,
		},
		{
			MethodName: "GetGameServer",
			Handler:    _SDK_GetGameServer_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Health",
			Handler:       _SDK_Health_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "sdk.proto",
}

func init() { proto.RegisterFile("sdk.proto", fileDescriptor_sdk_39d49c489951ac60) }

var fileDescriptor_sdk_39d49c489951ac60 = []byte{
	// 699 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0xdd, 0x6e, 0xd3, 0x4a,
	0x10, 0xc7, 0xe5, 0xc4, 0x76, 0xe2, 0xc9, 0xe9, 0x39, 0xe9, 0xb6, 0x47, 0xb2, 0xac, 0x0a, 0x95,
	0x08, 0xa1, 0x50, 0x51, 0x07, 0x95, 0x9b, 0x12, 0x09, 0xc4, 0x47, 0x4b, 0x41, 0x50, 0x51, 0x39,
	0x55, 0x2f, 0x2a, 0xa4, 0x68, 0x13, 0x8f, 0x12, 0x13, 0xc7, 0x6b, 0xed, 0x6e, 0x82, 0x72, 0xcb,
	0x2b, 0xf0, 0x10, 0x70, 0x03, 0x2f, 0xc3, 0x2b, 0xf0, 0x10, 0x5c, 0x21, 0xb4, 0x6b, 0xbb, 0x09,
	0xa5, 0xf4, 0x43, 0xbd, 0xca, 0xec, 0xcc, 0xfc, 0x7f, 0xb3, 0xd9, 0x99, 0x5d, 0x83, 0x23, 0xc2,
	0x91, 0x9f, 0x72, 0x26, 0x19, 0xf9, 0x5f, 0x48, 0xda, 0x8b, 0xd1, 0xa7, 0x03, 0x96, 0xa0, 0xf0,
	0x43, 0x9c, 0xfa, 0x22, 0x1c, 0x79, 0x6b, 0x03, 0xc6, 0x06, 0x31, 0xb6, 0x68, 0x1a, 0xb5, 0x68,
	0x92, 0x30, 0x49, 0x65, 0xc4, 0x12, 0x91, 0x89, 0x1a, 0x15, 0xb0, 0x76, 0xc7, 0xa9, 0x9c, 0x35,
	0xbe, 0x3a, 0x00, 0x7b, 0x74, 0x8c, 0x1d, 0xe4, 0x53, 0xe4, 0x64, 0x1f, 0x6a, 0xac, 0xf7, 0x0e,
	0xfb, 0xb2, 0x3b, 0x46, 0x49, 0x5d, 0x63, 0xdd, 0x68, 0xd6, 0xb6, 0xee, 0xfa, 0x67, 0x96, 0xf0,
	0xe7, 0x3a, 0xff, 0x8d, 0x16, 0xed, 0xa3, 0xa4, 0x01, 0xb0, 0x13, 0x9b, 0xb4, 0xc1, 0x14, 0x29,
	0xf6, 0xdd, 0x92, 0xe6, 0xdc, 0xbe, 0x98, 0xd3, 0x49, 0xb1, 0x1f, 0x68, 0x0d, 0x79, 0x0c, 0xb6,
	0x90, 0x54, 0x4e, 0x84, 0x5b, 0xd6, 0xea, 0xe6, 0x25, 0xd4, 0x3a, 0x3f, 0xc8, 0x75, 0xde, 0x27,
	0x13, 0x60, 0xbe, 0x31, 0x42, 0xc0, 0x4c, 0xe8, 0x18, 0xf5, 0x9f, 0x72, 0x02, 0x6d, 0x93, 0x35,
	0x70, 0xd4, 0xaf, 0x48, 0x69, 0x1f, 0xf5, 0x2e, 0x9d, 0x60, 0xee, 0x20, 0x75, 0x28, 0x4f, 0xa2,
	0x50, 0xd7, 0x77, 0x02, 0x65, 0x92, 0x3b, 0x50, 0xe7, 0x28, 0xd8, 0x84, 0xf7, 0xb1, 0x3b, 0x45,
	0x2e, 0x22, 0x96, 0xb8, 0xa6, 0x0e, 0xff, 0x57, 0xf8, 0x8f, 0x32, 0x37, 0xb9, 0x01, 0x30, 0xc0,
	0x04, 0xb9, 0x3e, 0x77, 0xd7, 0x5a, 0x37, 0x9a, 0xe5, 0x60, 0xc1, 0x43, 0x36, 0x81, 0xf4, 0x39,
	0x6a, 0xbb, 0x2b, 0xa3, 0x31, 0x0a, 0x49, 0xc7, 0xa9, 0x6b, 0xeb, 0xbc, 0xe5, 0x22, 0x72, 0x58,
	0x04, 0x54, 0x7a, 0x88, 0x31, 0x9e, 0x4a, 0xaf, 0x64, 0xe9, 0x45, 0x64, 0x9e, 0xde, 0x85, 0xda,
	0x42, 0xd7, 0xdd, 0xea, 0x7a, 0xb9, 0x59, 0xdb, 0x7a, 0x78, 0x95, 0x46, 0xfa, 0x4f, 0xe6, 0xfa,
	0xdd, 0x44, 0xf2, 0x59, 0xb0, 0x48, 0x24, 0x07, 0x60, 0xc7, 0xb4, 0x87, 0xb1, 0x70, 0x1d, 0xcd,
	0xde, 0xbe, 0x12, 0xfb, 0xb5, 0x96, 0x66, 0xd8, 0x9c, 0xe3, 0x3d, 0x82, 0xfa, 0xe9, 0x92, 0xaa,
	0x03, 0x23, 0x9c, 0xe5, 0x2d, 0x53, 0x26, 0x59, 0x05, 0x6b, 0x4a, 0xe3, 0x49, 0xd1, 0xad, 0x6c,
	0xd1, 0x2e, 0x6d, 0x1b, 0xde, 0x03, 0xa8, 0x2d, 0x60, 0xaf, 0x24, 0xfd, 0x61, 0x80, 0xa9, 0x46,
	0x8f, 0xec, 0x82, 0x3d, 0x44, 0x1a, 0xcb, 0x61, 0x3e, 0xfa, 0x9b, 0x97, 0x1b, 0x59, 0xff, 0x85,
	0x16, 0x05, 0xb9, 0xd8, 0xfb, 0x6c, 0x80, 0x9d, 0xb9, 0x88, 0x07, 0xd5, 0x9d, 0x48, 0x28, 0x46,
	0xa8, 0x99, 0xd5, 0xe0, 0x64, 0x4d, 0x6e, 0xc1, 0xd2, 0x01, 0xf2, 0x88, 0x85, 0x1d, 0xec, 0xb3,
	0x24, 0x14, 0x7a, 0x63, 0x56, 0xf0, 0xbb, 0x93, 0x6c, 0x40, 0xfd, 0x39, 0x8d, 0xe2, 0x09, 0xc7,
	0xc3, 0x21, 0x47, 0x31, 0x64, 0x71, 0x36, 0x92, 0x56, 0xf0, 0x87, 0x9f, 0xdc, 0x83, 0x95, 0x97,
	0x49, 0x24, 0x23, 0x1a, 0xef, 0x60, 0x4c, 0x67, 0x05, 0xd7, 0xd4, 0xe9, 0x67, 0x85, 0xbc, 0x2f,
	0x06, 0xd8, 0xd9, 0xbd, 0x51, 0xe7, 0xa3, 0x6e, 0x4e, 0x71, 0x43, 0xb2, 0x05, 0x71, 0xa1, 0x42,
	0xc3, 0x90, 0xa3, 0x10, 0xf9, 0xb9, 0x15, 0x4b, 0xf2, 0x0c, 0xac, 0x94, 0x71, 0xa9, 0x2e, 0x68,
	0xf9, 0x92, 0x67, 0xa5, 0x0b, 0xf9, 0x07, 0x8c, 0xcb, 0x20, 0xd3, 0x7a, 0x3e, 0x98, 0x6a, 0x79,
	0xe6, 0xed, 0x24, 0x60, 0xaa, 0xa4, 0xfc, 0x58, 0xb4, 0xbd, 0xf5, 0xb3, 0x04, 0xe5, 0xce, 0xce,
	0x2b, 0x72, 0x04, 0x56, 0x80, 0x34, 0x9c, 0x91, 0xb5, 0xbf, 0x94, 0xd5, 0xef, 0x9b, 0x77, 0x6e,
	0xb4, 0xb1, 0xfc, 0xe1, 0xdb, 0xf7, 0x8f, 0xa5, 0x5a, 0xc3, 0x6e, 0x71, 0xc5, 0x6a, 0x1b, 0x1b,
	0xe4, 0x2d, 0x54, 0x3b, 0xc3, 0x89, 0x0c, 0xd9, 0xfb, 0xe4, 0x5a, 0xe8, 0x55, 0x8d, 0xfe, 0xb7,
	0xe1, 0xb4, 0x44, 0x8e, 0x53, 0xf4, 0xe3, 0x93, 0xb9, 0xb8, 0x0e, 0x9b, 0x68, 0xf6, 0x3f, 0x8d,
	0x4a, 0x2b, 0x9b, 0xb7, 0xb6, 0xb1, 0xd1, 0x34, 0x08, 0xc2, 0xd2, 0x1e, 0xca, 0x85, 0xc7, 0xfc,
	0xfc, 0x12, 0x37, 0x2f, 0x6c, 0x57, 0x63, 0x45, 0xd7, 0x59, 0x22, 0xb5, 0xd6, 0x40, 0xbd, 0x89,
	0xda, 0xf9, 0xd4, 0x3a, 0x2e, 0x8b, 0x70, 0xd4, 0xb3, 0xf5, 0x87, 0xe4, 0xfe, 0xaf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xc5, 0x0e, 0x54, 0x28, 0x8a, 0x06, 0x00, 0x00,
}

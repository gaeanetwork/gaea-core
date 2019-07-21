// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shareddata.proto

package tee

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SharedData struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Ciphertext           string   `protobuf:"bytes,2,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
	Hash                 string   `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Owner                string   `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	CreateSeconds        int64    `protobuf:"varint,6,opt,name=create_seconds,json=createSeconds,proto3" json:"create_seconds,omitempty"`
	UpdateSeconds        int64    `protobuf:"varint,7,opt,name=update_seconds,json=updateSeconds,proto3" json:"update_seconds,omitempty"`
	Signatures           []string `protobuf:"bytes,8,rep,name=signatures,proto3" json:"signatures,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SharedData) Reset()         { *m = SharedData{} }
func (m *SharedData) String() string { return proto.CompactTextString(m) }
func (*SharedData) ProtoMessage()    {}
func (*SharedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_shareddata_a35a89393d6472ce, []int{0}
}
func (m *SharedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SharedData.Unmarshal(m, b)
}
func (m *SharedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SharedData.Marshal(b, m, deterministic)
}
func (dst *SharedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SharedData.Merge(dst, src)
}
func (m *SharedData) XXX_Size() int {
	return xxx_messageInfo_SharedData.Size(m)
}
func (m *SharedData) XXX_DiscardUnknown() {
	xxx_messageInfo_SharedData.DiscardUnknown(m)
}

var xxx_messageInfo_SharedData proto.InternalMessageInfo

func (m *SharedData) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SharedData) GetCiphertext() string {
	if m != nil {
		return m.Ciphertext
	}
	return ""
}

func (m *SharedData) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *SharedData) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *SharedData) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *SharedData) GetCreateSeconds() int64 {
	if m != nil {
		return m.CreateSeconds
	}
	return 0
}

func (m *SharedData) GetUpdateSeconds() int64 {
	if m != nil {
		return m.UpdateSeconds
	}
	return 0
}

func (m *SharedData) GetSignatures() []string {
	if m != nil {
		return m.Signatures
	}
	return nil
}

func init() {
	proto.RegisterType((*SharedData)(nil), "tee.SharedData")
}

func init() { proto.RegisterFile("shareddata.proto", fileDescriptor_shareddata_a35a89393d6472ce) }

var fileDescriptor_shareddata_a35a89393d6472ce = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xc1, 0x4a, 0x04, 0x31,
	0x0c, 0x40, 0xe9, 0x74, 0x77, 0xd5, 0x88, 0x8b, 0x04, 0x0f, 0x3d, 0x2d, 0x45, 0x10, 0xf6, 0xe4,
	0xc5, 0x5f, 0xf0, 0x0b, 0x76, 0x3e, 0x40, 0x62, 0x1b, 0x6c, 0x2f, 0x6d, 0x69, 0x33, 0xe8, 0x97,
	0x7b, 0x96, 0x69, 0x91, 0x9d, 0x5b, 0xf2, 0xde, 0x3b, 0x84, 0xc0, 0x63, 0x0b, 0x54, 0xd9, 0x7b,
	0x12, 0x7a, 0x2d, 0x35, 0x4b, 0x46, 0x2d, 0xcc, 0xcf, 0xbf, 0x0a, 0x60, 0xee, 0xe6, 0x9d, 0x84,
	0xf0, 0x08, 0x53, 0xf4, 0x46, 0x59, 0x75, 0xbe, 0xbb, 0x4c, 0xd1, 0xe3, 0x09, 0xc0, 0xc5, 0x12,
	0xb8, 0x0a, 0xff, 0x88, 0x99, 0x3a, 0xdf, 0x10, 0x44, 0xd8, 0x05, 0x6a, 0xc1, 0xe8, 0x6e, 0xfa,
	0x8c, 0x16, 0xee, 0x3d, 0x37, 0x57, 0x63, 0x91, 0x98, 0x93, 0xd9, 0x75, 0xb5, 0x45, 0xf8, 0x04,
	0xfb, 0xfc, 0x9d, 0xb8, 0x9a, 0x7d, 0x77, 0x63, 0xc1, 0x17, 0x38, 0xba, 0xca, 0x24, 0xfc, 0xd1,
	0xd8, 0xe5, 0xe4, 0x9b, 0x39, 0x58, 0x75, 0xd6, 0x97, 0x87, 0x41, 0xe7, 0x01, 0xd7, 0x6c, 0x29,
	0x7e, 0x9b, 0xdd, 0x8c, 0x6c, 0xd0, 0xff, 0xec, 0x04, 0xd0, 0xe2, 0x57, 0x22, 0x59, 0x2a, 0x37,
	0x73, 0x6b, 0xf5, 0x7a, 0xf9, 0x95, 0x7c, 0x1e, 0xfa, 0x13, 0xde, 0xfe, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x79, 0x2b, 0xd1, 0xf7, 0x18, 0x01, 0x00, 0x00,
}
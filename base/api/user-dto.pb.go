// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: user-dto.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID             int64  `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Username       string `protobuf:"bytes,2,opt,name=Username,proto3" json:"Username,omitempty"`
	Password       string `protobuf:"bytes,3,opt,name=Password,proto3" json:"Password,omitempty"`
	Telnumber      string `protobuf:"bytes,4,opt,name=Telnumber,proto3" json:"Telnumber,omitempty"`
	Email          string `protobuf:"bytes,5,opt,name=Email,proto3" json:"Email,omitempty"`
	Salt           string `protobuf:"bytes,6,opt,name=Salt,proto3" json:"Salt,omitempty"`
	Gender         int32  `protobuf:"varint,7,opt,name=Gender,proto3" json:"Gender,omitempty"`
	Status         int32  `protobuf:"varint,8,opt,name=Status,proto3" json:"Status,omitempty"`
	CreatedAt      int64  `protobuf:"varint,9,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	CreateBy       int64  `protobuf:"varint,10,opt,name=CreateBy,proto3" json:"CreateBy,omitempty"`
	CreateTime     int64  `protobuf:"varint,11,opt,name=CreateTime,proto3" json:"CreateTime,omitempty"`
	UpdatedAt      int64  `protobuf:"varint,12,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
	LastModifyBy   int64  `protobuf:"varint,13,opt,name=LastModifyBy,proto3" json:"LastModifyBy,omitempty"`
	LastModifyTime int64  `protobuf:"varint,14,opt,name=LastModifyTime,proto3" json:"LastModifyTime,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_dto_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_dto_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_dto_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *User) GetTelnumber() string {
	if x != nil {
		return x.Telnumber
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetSalt() string {
	if x != nil {
		return x.Salt
	}
	return ""
}

func (x *User) GetGender() int32 {
	if x != nil {
		return x.Gender
	}
	return 0
}

func (x *User) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *User) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *User) GetCreateBy() int64 {
	if x != nil {
		return x.CreateBy
	}
	return 0
}

func (x *User) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *User) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *User) GetLastModifyBy() int64 {
	if x != nil {
		return x.LastModifyBy
	}
	return 0
}

func (x *User) GetLastModifyTime() int64 {
	if x != nil {
		return x.LastModifyTime
	}
	return 0
}

var File_user_dto_proto protoreflect.FileDescriptor

var file_user_dto_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x2d, 0x64, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x76, 0x31, 0x22, 0x8a, 0x03, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a,
	0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x65, 0x6c, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x54, 0x65, 0x6c, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x61, 0x6c,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x53, 0x61, 0x6c, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x47,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c, 0x0a,
	0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x4d, 0x6f, 0x64,
	0x69, 0x66, 0x79, 0x42, 0x79, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x4c, 0x61, 0x73,
	0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x42, 0x79, 0x12, 0x26, 0x0a, 0x0e, 0x4c, 0x61, 0x73,
	0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0e, 0x4c, 0x61, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d,
	0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6e, 0x6d, 0x2f,
	0x6a, 0x69, 0x6e, 0x76, 0x65, 0x69, 0x2f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_dto_proto_rawDescOnce sync.Once
	file_user_dto_proto_rawDescData = file_user_dto_proto_rawDesc
)

func file_user_dto_proto_rawDescGZIP() []byte {
	file_user_dto_proto_rawDescOnce.Do(func() {
		file_user_dto_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_dto_proto_rawDescData)
	})
	return file_user_dto_proto_rawDescData
}

var file_user_dto_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_user_dto_proto_goTypes = []interface{}{
	(*User)(nil), // 0: v1.User
}
var file_user_dto_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_user_dto_proto_init() }
func file_user_dto_proto_init() {
	if File_user_dto_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_dto_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_dto_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_user_dto_proto_goTypes,
		DependencyIndexes: file_user_dto_proto_depIdxs,
		MessageInfos:      file_user_dto_proto_msgTypes,
	}.Build()
	File_user_dto_proto = out.File
	file_user_dto_proto_rawDesc = nil
	file_user_dto_proto_goTypes = nil
	file_user_dto_proto_depIdxs = nil
}

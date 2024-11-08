// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v4.25.2
// source: stats/v1/stats.proto

package v1

import (
	_ "github.com/alta/protopatch/patch/gopb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// AvailabilityState состоиние статистики
type AvailabilityState int32

const (
	// статистика не задана
	Unspecified AvailabilityState = 0
	// статистика доступна
	Available AvailabilityState = 1
	// ошибка при инициализации провайдера статистики
	Error AvailabilityState = 2
	// статистика отключена на сервере
	Disabled AvailabilityState = 3
)

// Enum value maps for AvailabilityState.
var (
	AvailabilityState_name = map[int32]string{
		0: "AVAILABILITY_STATE_UNSPECIFIED",
		1: "AVAILABILITY_STATE_AVAILABLE",
		2: "AVAILABILITY_STATE_ERROR",
		3: "AVAILABILITY_STATE_DISABLED",
	}
	AvailabilityState_value = map[string]int32{
		"AVAILABILITY_STATE_UNSPECIFIED": 0,
		"AVAILABILITY_STATE_AVAILABLE":   1,
		"AVAILABILITY_STATE_ERROR":       2,
		"AVAILABILITY_STATE_DISABLED":    3,
	}
)

func (x AvailabilityState) Enum() *AvailabilityState {
	p := new(AvailabilityState)
	*p = x
	return p
}

func (x AvailabilityState) pbString() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AvailabilityState) Descriptor() protoreflect.EnumDescriptor {
	return file_stats_v1_stats_proto_enumTypes[0].Descriptor()
}

func (AvailabilityState) Type() protoreflect.EnumType {
	return &file_stats_v1_stats_proto_enumTypes[0]
}

func (x AvailabilityState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AvailabilityState.Descriptor instead.
func (AvailabilityState) EnumDescriptor() ([]byte, []int) {
	return file_stats_v1_stats_proto_rawDescGZIP(), []int{0}
}

// AvailabilityDetails детальная информация о состоянии и причинах недоступности статистики.
type AvailabilityDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// State состояние статистики
	State AvailabilityState `protobuf:"varint,1,opt,name=state,proto3,enum=stats.v1.AvailabilityState" json:"state,omitempty"`
	// Details причина по которой статистика не доступна (возможно пустое).
	Details string `protobuf:"bytes,2,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *AvailabilityDetails) Reset() {
	*x = AvailabilityDetails{}
	mi := &file_stats_v1_stats_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AvailabilityDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AvailabilityDetails) ProtoMessage() {}

func (x *AvailabilityDetails) ProtoReflect() protoreflect.Message {
	mi := &file_stats_v1_stats_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AvailabilityDetails.ProtoReflect.Descriptor instead.
func (*AvailabilityDetails) Descriptor() ([]byte, []int) {
	return file_stats_v1_stats_proto_rawDescGZIP(), []int{0}
}

func (x *AvailabilityDetails) GetState() AvailabilityState {
	if x != nil {
		return x.State
	}
	return Unspecified
}

func (x *AvailabilityDetails) GetDetails() string {
	if x != nil {
		return x.Details
	}
	return ""
}

// Provider информация о провайдере статистики.
type Provider struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ProviderID идентификатор провайдера статистики для записи
	ProviderID string `protobuf:"bytes,1,opt,name=provider_id,json=providerId,proto3" json:"provider_id,omitempty"`
	// ProviderName название провайдера статистики для записи.
	ProviderName string `protobuf:"bytes,2,opt,name=provider_name,json=providerName,proto3" json:"provider_name,omitempty"`
	// Platform название платформы, на которой работает провайдер статистики
	Platform string `protobuf:"bytes,3,opt,name=platform,proto3" json:"platform,omitempty"`
	// AvailabilityDetails информация о доступности провайдера.
	AvailabilityDetails *AvailabilityDetails `protobuf:"bytes,4,opt,name=availability_details,json=availabilityDetails,proto3" json:"availability_details,omitempty"`
}

func (x *Provider) Reset() {
	*x = Provider{}
	mi := &file_stats_v1_stats_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Provider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Provider) ProtoMessage() {}

func (x *Provider) ProtoReflect() protoreflect.Message {
	mi := &file_stats_v1_stats_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Provider.ProtoReflect.Descriptor instead.
func (*Provider) Descriptor() ([]byte, []int) {
	return file_stats_v1_stats_proto_rawDescGZIP(), []int{1}
}

func (x *Provider) GetProviderID() string {
	if x != nil {
		return x.ProviderID
	}
	return ""
}

func (x *Provider) GetProviderName() string {
	if x != nil {
		return x.ProviderName
	}
	return ""
}

func (x *Provider) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *Provider) GetAvailabilityDetails() *AvailabilityDetails {
	if x != nil {
		return x.AvailabilityDetails
	}
	return nil
}

// Record представляет собой запись статистики на указанную дату-время для произвольного провайдера статистики.
type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Provider - провайдер, к которому относится запись статистики.
	Provider *Provider `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	// RecordValue список значений. Провайдеры могут отдавать несколько значений за раз.
	// Например, loadavg отдаёт три значения:
	//   - ср.загрузка за 1мин
	//   - ср.загрузка за 5мин
	//   - ср.загрузка за 15мин
	Value []*RecordValue `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	// Time время соотвестствующее записи статистики.
	Time *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	mi := &file_stats_v1_stats_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_stats_v1_stats_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_stats_v1_stats_proto_rawDescGZIP(), []int{2}
}

func (x *Record) GetProvider() *Provider {
	if x != nil {
		return x.Provider
	}
	return nil
}

func (x *Record) GetValue() []*RecordValue {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Record) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

// RecordValue представляет собой значение соответствующей записи.
type RecordValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID идентификатор значения
	ID string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name название значения
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Value само значение
	Value string `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *RecordValue) Reset() {
	*x = RecordValue{}
	mi := &file_stats_v1_stats_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecordValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordValue) ProtoMessage() {}

func (x *RecordValue) ProtoReflect() protoreflect.Message {
	mi := &file_stats_v1_stats_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordValue.ProtoReflect.Descriptor instead.
func (*RecordValue) Descriptor() ([]byte, []int) {
	return file_stats_v1_stats_proto_rawDescGZIP(), []int{3}
}

func (x *RecordValue) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *RecordValue) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RecordValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_stats_v1_stats_proto protoreflect.FileDescriptor

var file_stats_v1_stats_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x76, 0x31,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0e, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x62, 0x0a, 0x13, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74,
	0x79, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x31, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0xe4, 0x01, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x12, 0x31, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x10, 0xca, 0xb5, 0x03, 0x0c, 0x0a, 0x0a, 0x50,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x49, 0x44, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x37, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x12, 0xca, 0xb5,
	0x03, 0x0e, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x52, 0x0c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x50, 0x0a, 0x14, 0x61, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x13, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x95, 0x01, 0x0a,
	0x06, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x2e, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x74, 0x61, 0x74,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x08, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04,
	0x74, 0x69, 0x6d, 0x65, 0x22, 0x51, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x08, 0xca, 0xb5, 0x03, 0x04, 0x0a, 0x02, 0x49, 0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0xea, 0x01, 0x0a, 0x11, 0x41, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x35, 0x0a,
	0x1e, 0x41, 0x56, 0x41, 0x49, 0x4c, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x1a, 0x11, 0xca, 0xb5, 0x03, 0x0d, 0x0a, 0x0b, 0x55, 0x6e, 0x73, 0x70, 0x65, 0x63, 0x69,
	0x66, 0x69, 0x65, 0x64, 0x12, 0x31, 0x0a, 0x1c, 0x41, 0x56, 0x41, 0x49, 0x4c, 0x41, 0x42, 0x49,
	0x4c, 0x49, 0x54, 0x59, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x56, 0x41, 0x49, 0x4c,
	0x41, 0x42, 0x4c, 0x45, 0x10, 0x01, 0x1a, 0x0f, 0xca, 0xb5, 0x03, 0x0b, 0x0a, 0x09, 0x41, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x29, 0x0a, 0x18, 0x41, 0x56, 0x41, 0x49, 0x4c,
	0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x10, 0x02, 0x1a, 0x0b, 0xca, 0xb5, 0x03, 0x07, 0x0a, 0x05, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x2f, 0x0a, 0x1b, 0x41, 0x56, 0x41, 0x49, 0x4c, 0x41, 0x42, 0x49, 0x4c, 0x49,
	0x54, 0x59, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x44, 0x49, 0x53, 0x41, 0x42, 0x4c, 0x45,
	0x44, 0x10, 0x03, 0x1a, 0x0e, 0xca, 0xb5, 0x03, 0x0a, 0x0a, 0x08, 0x44, 0x69, 0x73, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x1a, 0x0f, 0xca, 0xb5, 0x03, 0x0b, 0xf2, 0x01, 0x08, 0x70, 0x62, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x64, 0x69, 0x6d, 0x61, 0x2d, 0x73, 0x74, 0x75, 0x64, 0x79, 0x2f, 0x6d, 0x6f,
	0x6e, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_stats_v1_stats_proto_rawDescOnce sync.Once
	file_stats_v1_stats_proto_rawDescData = file_stats_v1_stats_proto_rawDesc
)

func file_stats_v1_stats_proto_rawDescGZIP() []byte {
	file_stats_v1_stats_proto_rawDescOnce.Do(func() {
		file_stats_v1_stats_proto_rawDescData = protoimpl.X.CompressGZIP(file_stats_v1_stats_proto_rawDescData)
	})
	return file_stats_v1_stats_proto_rawDescData
}

var file_stats_v1_stats_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_stats_v1_stats_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_stats_v1_stats_proto_goTypes = []any{
	(AvailabilityState)(0),        // 0: stats.v1.AvailabilityState
	(*AvailabilityDetails)(nil),   // 1: stats.v1.AvailabilityDetails
	(*Provider)(nil),              // 2: stats.v1.Provider
	(*Record)(nil),                // 3: stats.v1.Record
	(*RecordValue)(nil),           // 4: stats.v1.RecordValue
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_stats_v1_stats_proto_depIdxs = []int32{
	0, // 0: stats.v1.AvailabilityDetails.state:type_name -> stats.v1.AvailabilityState
	1, // 1: stats.v1.Provider.availability_details:type_name -> stats.v1.AvailabilityDetails
	2, // 2: stats.v1.Record.provider:type_name -> stats.v1.Provider
	4, // 3: stats.v1.Record.value:type_name -> stats.v1.RecordValue
	5, // 4: stats.v1.Record.time:type_name -> google.protobuf.Timestamp
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_stats_v1_stats_proto_init() }
func file_stats_v1_stats_proto_init() {
	if File_stats_v1_stats_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stats_v1_stats_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_stats_v1_stats_proto_goTypes,
		DependencyIndexes: file_stats_v1_stats_proto_depIdxs,
		EnumInfos:         file_stats_v1_stats_proto_enumTypes,
		MessageInfos:      file_stats_v1_stats_proto_msgTypes,
	}.Build()
	File_stats_v1_stats_proto = out.File
	file_stats_v1_stats_proto_rawDesc = nil
	file_stats_v1_stats_proto_goTypes = nil
	file_stats_v1_stats_proto_depIdxs = nil
}

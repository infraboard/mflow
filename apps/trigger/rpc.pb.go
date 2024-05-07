// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.0
// source: mflow/apps/trigger/pb/rpc.proto

package trigger

import (
	request "github.com/infraboard/mcube/v2/http/request"
	resource "github.com/infraboard/mcube/v2/pb/resource"
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

type DELETE_BY int32

const (
	// id
	DELETE_BY_RECORD_ID DELETE_BY = 0
	// task id
	DELETE_BY_PIPELINE_TASK_ID DELETE_BY = 1
)

// Enum value maps for DELETE_BY.
var (
	DELETE_BY_name = map[int32]string{
		0: "RECORD_ID",
		1: "PIPELINE_TASK_ID",
	}
	DELETE_BY_value = map[string]int32{
		"RECORD_ID":        0,
		"PIPELINE_TASK_ID": 1,
	}
)

func (x DELETE_BY) Enum() *DELETE_BY {
	p := new(DELETE_BY)
	*p = x
	return p
}

func (x DELETE_BY) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DELETE_BY) Descriptor() protoreflect.EnumDescriptor {
	return file_mflow_apps_trigger_pb_rpc_proto_enumTypes[0].Descriptor()
}

func (DELETE_BY) Type() protoreflect.EnumType {
	return &file_mflow_apps_trigger_pb_rpc_proto_enumTypes[0]
}

func (x DELETE_BY) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DELETE_BY.Descriptor instead.
func (DELETE_BY) EnumDescriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_rpc_proto_rawDescGZIP(), []int{0}
}

type DeleteRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @gotags: json:"delete_by"
	DeleteBy DELETE_BY `protobuf:"varint,1,opt,name=delete_by,json=deleteBy,proto3,enum=infraboard.mflow.trigger.DELETE_BY" json:"delete_by"`
	// @gotags: json:"values"
	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values"`
}

func (x *DeleteRecordRequest) Reset() {
	*x = DeleteRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_trigger_pb_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRecordRequest) ProtoMessage() {}

func (x *DeleteRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_trigger_pb_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRecordRequest.ProtoReflect.Descriptor instead.
func (*DeleteRecordRequest) Descriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteRecordRequest) GetDeleteBy() DELETE_BY {
	if x != nil {
		return x.DeleteBy
	}
	return DELETE_BY_RECORD_ID
}

func (x *DeleteRecordRequest) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type DescribeRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @gotags: json:"id"
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
}

func (x *DescribeRecordRequest) Reset() {
	*x = DescribeRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_trigger_pb_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeRecordRequest) ProtoMessage() {}

func (x *DescribeRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_trigger_pb_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeRecordRequest.ProtoReflect.Descriptor instead.
func (*DescribeRecordRequest) Descriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *DescribeRecordRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type QueryRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 分页请求
	// @gotags: json:"page"
	Page *request.PageRequest `protobuf:"bytes,1,opt,name=page,proto3" json:"page"`
	// 服务Id, 查询某个服务的事件
	// @gotags: json:"service_id"
	ServiceId string `protobuf:"bytes,2,opt,name=service_id,json=serviceId,proto3" json:"service_id"`
	// 查询PipelineTask关联的事件
	// @gotags: json:"pipeline_task_id"
	PipelineTaskId string `protobuf:"bytes,3,opt,name=pipeline_task_id,json=pipelineTaskId,proto3" json:"pipeline_task_id"`
	// 相关构建配置关联的记录
	// @gotags: json:"build_conf_ids"
	BuildConfIds []string `protobuf:"bytes,6,rep,name=build_conf_ids,json=buildConfIds,proto3" json:"build_conf_ids"`
	// 资源范围
	// @gotags: json:"scope"
	Scope *resource.Scope `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope"`
	// 是否查询出管理的pipeline task对象
	// @gotags: json:"with_pipeline_task"
	WithPipelineTask bool `protobuf:"varint,7,opt,name=with_pipeline_task,json=withPipelineTask,proto3" json:"with_pipeline_task"`
	// 资源标签过滤
	// @gotags: json:"filters"
	Filters []*resource.LabelRequirement `protobuf:"bytes,5,rep,name=filters,proto3" json:"filters"`
}

func (x *QueryRecordRequest) Reset() {
	*x = QueryRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_trigger_pb_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRecordRequest) ProtoMessage() {}

func (x *QueryRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_trigger_pb_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRecordRequest.ProtoReflect.Descriptor instead.
func (*QueryRecordRequest) Descriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *QueryRecordRequest) GetPage() *request.PageRequest {
	if x != nil {
		return x.Page
	}
	return nil
}

func (x *QueryRecordRequest) GetServiceId() string {
	if x != nil {
		return x.ServiceId
	}
	return ""
}

func (x *QueryRecordRequest) GetPipelineTaskId() string {
	if x != nil {
		return x.PipelineTaskId
	}
	return ""
}

func (x *QueryRecordRequest) GetBuildConfIds() []string {
	if x != nil {
		return x.BuildConfIds
	}
	return nil
}

func (x *QueryRecordRequest) GetScope() *resource.Scope {
	if x != nil {
		return x.Scope
	}
	return nil
}

func (x *QueryRecordRequest) GetWithPipelineTask() bool {
	if x != nil {
		return x.WithPipelineTask
	}
	return false
}

func (x *QueryRecordRequest) GetFilters() []*resource.LabelRequirement {
	if x != nil {
		return x.Filters
	}
	return nil
}

type EventQueueTaskCompleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 执行完成的PipelineTask任务Id
	// @gotags: json:"pipeline_task_id"
	PipelineTaskId string `protobuf:"bytes,1,opt,name=pipeline_task_id,json=pipelineTaskId,proto3" json:"pipeline_task_id"`
}

func (x *EventQueueTaskCompleteRequest) Reset() {
	*x = EventQueueTaskCompleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_trigger_pb_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventQueueTaskCompleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventQueueTaskCompleteRequest) ProtoMessage() {}

func (x *EventQueueTaskCompleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_trigger_pb_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventQueueTaskCompleteRequest.ProtoReflect.Descriptor instead.
func (*EventQueueTaskCompleteRequest) Descriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *EventQueueTaskCompleteRequest) GetPipelineTaskId() string {
	if x != nil {
		return x.PipelineTaskId
	}
	return ""
}

var File_mflow_apps_trigger_pb_rpc_proto protoreflect.FileDescriptor

var file_mflow_apps_trigger_pb_rpc_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x18, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66,
	0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x1a, 0x21, 0x6d, 0x66, 0x6c,
	0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2f,
	0x70, 0x62, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18,
	0x6d, 0x63, 0x75, 0x62, 0x65, 0x2f, 0x70, 0x62, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x2f, 0x70, 0x61,
	0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x63, 0x75, 0x62, 0x65, 0x2f,
	0x70, 0x62, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x6d, 0x65, 0x74, 0x61,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x6d, 0x63, 0x75, 0x62, 0x65, 0x2f, 0x70, 0x62,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x40, 0x0a, 0x09,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x23, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c,
	0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x45, 0x4c, 0x45, 0x54,
	0x45, 0x5f, 0x42, 0x59, 0x52, 0x08, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x27, 0x0a, 0x15, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22,
	0xe8, 0x02, 0x0a, 0x12, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x2e, 0x6d, 0x63, 0x75, 0x62, 0x65, 0x2e, 0x70, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x28, 0x0a,
	0x10, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e,
	0x65, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x62, 0x75, 0x69, 0x6c, 0x64,
	0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0c, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x49, 0x64, 0x73, 0x12, 0x36, 0x0a,
	0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69,
	0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x63, 0x75, 0x62, 0x65, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x52, 0x05,
	0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x70, 0x69,
	0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x10, 0x77, 0x69, 0x74, 0x68, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x54,
	0x61, 0x73, 0x6b, 0x12, 0x45, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x2e, 0x6d, 0x63, 0x75, 0x62, 0x65, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x22, 0x49, 0x0a, 0x1d, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x70,
	0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x54,
	0x61, 0x73, 0x6b, 0x49, 0x64, 0x2a, 0x30, 0x0a, 0x09, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f,
	0x42, 0x59, 0x12, 0x0d, 0x0a, 0x09, 0x52, 0x45, 0x43, 0x4f, 0x52, 0x44, 0x5f, 0x49, 0x44, 0x10,
	0x00, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x49, 0x50, 0x45, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x54, 0x41,
	0x53, 0x4b, 0x5f, 0x49, 0x44, 0x10, 0x01, 0x32, 0x9e, 0x02, 0x0a, 0x03, 0x52, 0x50, 0x43, 0x12,
	0x50, 0x0a, 0x0b, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1f,
	0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a,
	0x20, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c,
	0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x12, 0x60, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x12, 0x2c, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66,
	0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23,
	0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x53, 0x65, 0x74, 0x12, 0x63, 0x0a, 0x0e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x2f, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72,
	0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f,
	0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x2f, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mflow_apps_trigger_pb_rpc_proto_rawDescOnce sync.Once
	file_mflow_apps_trigger_pb_rpc_proto_rawDescData = file_mflow_apps_trigger_pb_rpc_proto_rawDesc
)

func file_mflow_apps_trigger_pb_rpc_proto_rawDescGZIP() []byte {
	file_mflow_apps_trigger_pb_rpc_proto_rawDescOnce.Do(func() {
		file_mflow_apps_trigger_pb_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_mflow_apps_trigger_pb_rpc_proto_rawDescData)
	})
	return file_mflow_apps_trigger_pb_rpc_proto_rawDescData
}

var file_mflow_apps_trigger_pb_rpc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_mflow_apps_trigger_pb_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_mflow_apps_trigger_pb_rpc_proto_goTypes = []interface{}{
	(DELETE_BY)(0),                        // 0: infraboard.mflow.trigger.DELETE_BY
	(*DeleteRecordRequest)(nil),           // 1: infraboard.mflow.trigger.DeleteRecordRequest
	(*DescribeRecordRequest)(nil),         // 2: infraboard.mflow.trigger.DescribeRecordRequest
	(*QueryRecordRequest)(nil),            // 3: infraboard.mflow.trigger.QueryRecordRequest
	(*EventQueueTaskCompleteRequest)(nil), // 4: infraboard.mflow.trigger.EventQueueTaskCompleteRequest
	(*request.PageRequest)(nil),           // 5: infraboard.mcube.page.PageRequest
	(*resource.Scope)(nil),                // 6: infraboard.mcube.resource.Scope
	(*resource.LabelRequirement)(nil),     // 7: infraboard.mcube.resource.LabelRequirement
	(*Event)(nil),                         // 8: infraboard.mflow.trigger.Event
	(*Record)(nil),                        // 9: infraboard.mflow.trigger.Record
	(*RecordSet)(nil),                     // 10: infraboard.mflow.trigger.RecordSet
}
var file_mflow_apps_trigger_pb_rpc_proto_depIdxs = []int32{
	0,  // 0: infraboard.mflow.trigger.DeleteRecordRequest.delete_by:type_name -> infraboard.mflow.trigger.DELETE_BY
	5,  // 1: infraboard.mflow.trigger.QueryRecordRequest.page:type_name -> infraboard.mcube.page.PageRequest
	6,  // 2: infraboard.mflow.trigger.QueryRecordRequest.scope:type_name -> infraboard.mcube.resource.Scope
	7,  // 3: infraboard.mflow.trigger.QueryRecordRequest.filters:type_name -> infraboard.mcube.resource.LabelRequirement
	8,  // 4: infraboard.mflow.trigger.RPC.HandleEvent:input_type -> infraboard.mflow.trigger.Event
	3,  // 5: infraboard.mflow.trigger.RPC.QueryRecord:input_type -> infraboard.mflow.trigger.QueryRecordRequest
	2,  // 6: infraboard.mflow.trigger.RPC.DescribeRecord:input_type -> infraboard.mflow.trigger.DescribeRecordRequest
	9,  // 7: infraboard.mflow.trigger.RPC.HandleEvent:output_type -> infraboard.mflow.trigger.Record
	10, // 8: infraboard.mflow.trigger.RPC.QueryRecord:output_type -> infraboard.mflow.trigger.RecordSet
	9,  // 9: infraboard.mflow.trigger.RPC.DescribeRecord:output_type -> infraboard.mflow.trigger.Record
	7,  // [7:10] is the sub-list for method output_type
	4,  // [4:7] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_mflow_apps_trigger_pb_rpc_proto_init() }
func file_mflow_apps_trigger_pb_rpc_proto_init() {
	if File_mflow_apps_trigger_pb_rpc_proto != nil {
		return
	}
	file_mflow_apps_trigger_pb_event_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_mflow_apps_trigger_pb_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRecordRequest); i {
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
		file_mflow_apps_trigger_pb_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeRecordRequest); i {
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
		file_mflow_apps_trigger_pb_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryRecordRequest); i {
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
		file_mflow_apps_trigger_pb_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventQueueTaskCompleteRequest); i {
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
			RawDescriptor: file_mflow_apps_trigger_pb_rpc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mflow_apps_trigger_pb_rpc_proto_goTypes,
		DependencyIndexes: file_mflow_apps_trigger_pb_rpc_proto_depIdxs,
		EnumInfos:         file_mflow_apps_trigger_pb_rpc_proto_enumTypes,
		MessageInfos:      file_mflow_apps_trigger_pb_rpc_proto_msgTypes,
	}.Build()
	File_mflow_apps_trigger_pb_rpc_proto = out.File
	file_mflow_apps_trigger_pb_rpc_proto_rawDesc = nil
	file_mflow_apps_trigger_pb_rpc_proto_goTypes = nil
	file_mflow_apps_trigger_pb_rpc_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.0
// source: mflow/apps/trigger/pb/event.proto

package trigger

import (
	build "github.com/infraboard/mflow/apps/build"
	task "github.com/infraboard/mflow/apps/task"
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

type STAGE int32

const (
	// 等待执行
	STAGE_PENDDING STAGE = 0
	// 已经放入执行队列
	STAGE_ENQUEUE STAGE = 1
	// 执行中
	STAGE_RUNNING STAGE = 2
	// 取消执行
	STAGE_CANCELED STAGE = 3
	// 触发执行成功
	STAGE_SUCCESS STAGE = 4
	// 触发执行失败
	STAGE_FAILED STAGE = 5
)

// Enum value maps for STAGE.
var (
	STAGE_name = map[int32]string{
		0: "PENDDING",
		1: "ENQUEUE",
		2: "RUNNING",
		3: "CANCELED",
		4: "SUCCESS",
		5: "FAILED",
	}
	STAGE_value = map[string]int32{
		"PENDDING": 0,
		"ENQUEUE":  1,
		"RUNNING":  2,
		"CANCELED": 3,
		"SUCCESS":  4,
		"FAILED":   5,
	}
)

func (x STAGE) Enum() *STAGE {
	p := new(STAGE)
	*p = x
	return p
}

func (x STAGE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (STAGE) Descriptor() protoreflect.EnumDescriptor {
	return file_mflow_apps_trigger_pb_event_proto_enumTypes[0].Descriptor()
}

func (STAGE) Type() protoreflect.EnumType {
	return &file_mflow_apps_trigger_pb_event_proto_enumTypes[0]
}

func (x STAGE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use STAGE.Descriptor instead.
func (STAGE) EnumDescriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_event_proto_rawDescGZIP(), []int{0}
}

type EVENT_PROVIDER int32

const (
	// 来自gitlab的事件
	EVENT_PROVIDER_GITLAB EVENT_PROVIDER = 0
)

// Enum value maps for EVENT_PROVIDER.
var (
	EVENT_PROVIDER_name = map[int32]string{
		0: "GITLAB",
	}
	EVENT_PROVIDER_value = map[string]int32{
		"GITLAB": 0,
	}
)

func (x EVENT_PROVIDER) Enum() *EVENT_PROVIDER {
	p := new(EVENT_PROVIDER)
	*p = x
	return p
}

func (x EVENT_PROVIDER) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EVENT_PROVIDER) Descriptor() protoreflect.EnumDescriptor {
	return file_mflow_apps_trigger_pb_event_proto_enumTypes[1].Descriptor()
}

func (EVENT_PROVIDER) Type() protoreflect.EnumType {
	return &file_mflow_apps_trigger_pb_event_proto_enumTypes[1]
}

func (x EVENT_PROVIDER) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EVENT_PROVIDER.Descriptor instead.
func (EVENT_PROVIDER) EnumDescriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_event_proto_rawDescGZIP(), []int{1}
}

type RecordSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 总数
	// @gotags: bson:"total" json:"total"
	Total int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total" bson:"total"`
	// 列表
	// @gotags: bson:"items" json:"items"
	Items []*Record `protobuf:"bytes,2,rep,name=items,proto3" json:"items" bson:"items"`
}

func (x *RecordSet) Reset() {
	*x = RecordSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_trigger_pb_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordSet) ProtoMessage() {}

func (x *RecordSet) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_trigger_pb_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordSet.ProtoReflect.Descriptor instead.
func (*RecordSet) Descriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_event_proto_rawDescGZIP(), []int{0}
}

func (x *RecordSet) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *RecordSet) GetItems() []*Record {
	if x != nil {
		return x.Items
	}
	return nil
}

// 事件记录
type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// event相关定义
	// @gotags: bson:",inline" json:"event"
	Event *Event `protobuf:"bytes,2,opt,name=event,proto3" json:"event" bson:",inline"`
	// 构建状态
	// @gotags: bson:"build_status" json:"build_status"
	BuildStatus []*BuildStatus `protobuf:"bytes,3,rep,name=build_status,json=buildStatus,proto3" json:"build_status" bson:"build_status"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_trigger_pb_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_trigger_pb_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_mflow_apps_trigger_pb_event_proto_rawDescGZIP(), []int{1}
}

func (x *Record) GetEvent() *Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *Record) GetBuildStatus() []*BuildStatus {
	if x != nil {
		return x.BuildStatus
	}
	return nil
}

type BuildStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 构建信息, 由于可能修改, 因此需要保存整个对象
	// @gotags: bson:"build_config" json:"build_config"
	BuildConfig *build.BuildConfig `protobuf:"bytes,1,opt,name=build_config,json=buildConfig,proto3" json:"build_config" bson:"build_config"`
	// 触发状态
	// @gotags: bson:"stage" json:"stage"
	Stage STAGE `protobuf:"varint,5,opt,name=stage,proto3,enum=infraboard.mflow.trigger.STAGE" json:"stage" bson:"stage"`
	// 构建信息
	// @gotags: bson:"pipline_task_id" json:"pipline_task_id"
	PiplineTaskId string `protobuf:"bytes,2,opt,name=pipline_task_id,json=piplineTaskId,proto3" json:"pipline_task_id" bson:"pipline_task_id"`
	// pipline具体信息, 动态查询，不保持
	// @gotags: bson:"-" json:"pipline_task"
	PiplineTask *task.PipelineTask `protobuf:"bytes,3,opt,name=pipline_task,json=piplineTask,proto3" json:"pipline_task" bson:"-"`
	// 如果流水线运行报错的报错信息
	// @gotags: bson:"error_message" json:"error_message"
	ErrorMessage string `protobuf:"bytes,4,opt,name=error_message,json=errorMessage,proto3" json:"error_message" bson:"error_message"`
}

func (x *BuildStatus) Reset() {
	*x = BuildStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_trigger_pb_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildStatus) ProtoMessage() {}

func (x *BuildStatus) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_trigger_pb_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildStatus.ProtoReflect.Descriptor instead.
func (*BuildStatus) Descriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_event_proto_rawDescGZIP(), []int{2}
}

func (x *BuildStatus) GetBuildConfig() *build.BuildConfig {
	if x != nil {
		return x.BuildConfig
	}
	return nil
}

func (x *BuildStatus) GetStage() STAGE {
	if x != nil {
		return x.Stage
	}
	return STAGE_PENDDING
}

func (x *BuildStatus) GetPiplineTaskId() string {
	if x != nil {
		return x.PiplineTaskId
	}
	return ""
}

func (x *BuildStatus) GetPiplineTask() *task.PipelineTask {
	if x != nil {
		return x.PiplineTask
	}
	return nil
}

func (x *BuildStatus) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 事件id
	// @gotags: json:"id" bson:"_id"
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id" bson:"_id"`
	// 不执行Pipeline, 用于调试
	// @gotags: json:"skip_run_pipeline" bson:"skip_run_pipeline"
	SkipRunPipeline bool `protobuf:"varint,2,opt,name=skip_run_pipeline,json=skipRunPipeline,proto3" json:"skip_run_pipeline" bson:"skip_run_pipeline"`
	// 模拟的事件, 用于手动触发测试
	// @gotags: json:"is_mock" bson:"is_mock"
	IsMock bool `protobuf:"varint,3,opt,name=is_mock,json=isMock,proto3" json:"is_mock" bson:"is_mock"`
	// 事件来源, 对于gitlab而言有 发送事件的gitlab服务地址, 比如 "https://gitlab.com"
	// @gotags: bson:"from" json:"from"
	From string `protobuf:"bytes,4,opt,name=from,proto3" json:"from" bson:"from"`
	// 事件时间
	// @gotags: bson:"time" json:"time"
	Time int64 `protobuf:"varint,12,opt,name=time,proto3" json:"time" bson:"time"`
	// 发送者的user agent信息, 比如 "GitLab/15.5.0-pre"
	// @gotags: bson:"user_agent" json:"user_agent"
	UserAgent string `protobuf:"bytes,5,opt,name=user_agent,json=userAgent,proto3" json:"user_agent" bson:"user_agent"`
	// 事件提供方
	// @gotags: json:"provider" bson:"provider"
	Provider EVENT_PROVIDER `protobuf:"varint,6,opt,name=provider,proto3,enum=infraboard.mflow.trigger.EVENT_PROVIDER" json:"provider" bson:"provider"`
	// 提供方的版本, 比如gitlab 15.5.0, 不通版本有同步版本的API, 因此需要记录
	// 这个值从user agent中获取
	// @gotags: bson:"provider_version" json:"provider_version"
	ProviderVersion string `protobuf:"bytes,11,opt,name=provider_version,json=providerVersion,proto3" json:"provider_version" bson:"provider_version"`
	// 事件名称
	// @gotags: bson:"name" json:"name"
	Name string `protobuf:"bytes,7,opt,name=name,proto3" json:"name" bson:"name"`
	// 事件子名称, 子名称支持正则匹配
	// @gotags: bson:"sub_name" json:"sub_name"
	SubName string `protobuf:"bytes,8,opt,name=sub_name,json=subName,proto3" json:"sub_name" bson:"sub_name"`
	// 事件Token, 一般固定为 服务Id
	// @gotags: json:"token" bson:"token"
	Token string `protobuf:"bytes,9,opt,name=token,proto3" json:"token" bson:"token"`
	// 服务相关信息
	// @gotags: bson:"service_info" json:"service_info"
	ServiceInfo string `protobuf:"bytes,13,opt,name=service_info,json=serviceInfo,proto3" json:"service_info" bson:"service_info"`
	// 事件原始数据
	// @gotags: json:"raw" bson:"raw"
	Raw string `protobuf:"bytes,10,opt,name=raw,proto3" json:"raw" bson:"raw"`
	// 原始数据解析异常
	// @gotags: bson:"parse_error" json:"parse_error"
	ParseError string `protobuf:"bytes,14,opt,name=parse_error,json=parseError,proto3" json:"parse_error" bson:"parse_error"`
	// 其他属性
	// @gotags: bson:"labels" json:"labels"
	Extra map[string]string `protobuf:"bytes,15,rep,name=extra,proto3" json:"labels" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" bson:"labels"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_trigger_pb_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_trigger_pb_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_mflow_apps_trigger_pb_event_proto_rawDescGZIP(), []int{3}
}

func (x *Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Event) GetSkipRunPipeline() bool {
	if x != nil {
		return x.SkipRunPipeline
	}
	return false
}

func (x *Event) GetIsMock() bool {
	if x != nil {
		return x.IsMock
	}
	return false
}

func (x *Event) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Event) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Event) GetUserAgent() string {
	if x != nil {
		return x.UserAgent
	}
	return ""
}

func (x *Event) GetProvider() EVENT_PROVIDER {
	if x != nil {
		return x.Provider
	}
	return EVENT_PROVIDER_GITLAB
}

func (x *Event) GetProviderVersion() string {
	if x != nil {
		return x.ProviderVersion
	}
	return ""
}

func (x *Event) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Event) GetSubName() string {
	if x != nil {
		return x.SubName
	}
	return ""
}

func (x *Event) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Event) GetServiceInfo() string {
	if x != nil {
		return x.ServiceInfo
	}
	return ""
}

func (x *Event) GetRaw() string {
	if x != nil {
		return x.Raw
	}
	return ""
}

func (x *Event) GetParseError() string {
	if x != nil {
		return x.ParseError
	}
	return ""
}

func (x *Event) GetExtra() map[string]string {
	if x != nil {
		return x.Extra
	}
	return nil
}

var File_mflow_apps_trigger_pb_event_proto protoreflect.FileDescriptor

var file_mflow_apps_trigger_pb_event_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x18, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e,
	0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x1a, 0x1f, 0x6d,
	0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f,
	0x70, 0x62, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x26,
	0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x2f,
	0x70, 0x62, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x74, 0x61, 0x73, 0x6b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x09, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x53, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x36, 0x0a, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x89, 0x01, 0x0a, 0x06, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x35, 0x0a, 0x05,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x69, 0x6e,
	0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x48, 0x0a, 0x0c, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x69, 0x6e, 0x66, 0x72,
	0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x0b, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xa1, 0x02,
	0x0a, 0x0b, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x46, 0x0a,
	0x0c, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x42, 0x75, 0x69,
	0x6c, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0b, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x35, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e,
	0x53, 0x54, 0x41, 0x47, 0x45, 0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x12, 0x26, 0x0a, 0x0f,
	0x70, 0x69, 0x70, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x69, 0x70, 0x6c, 0x69, 0x6e, 0x65, 0x54, 0x61,
	0x73, 0x6b, 0x49, 0x64, 0x12, 0x46, 0x0a, 0x0c, 0x70, 0x69, 0x70, 0x6c, 0x69, 0x6e, 0x65, 0x5f,
	0x74, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x69, 0x6e, 0x66,
	0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x61,
	0x73, 0x6b, 0x2e, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x0b, 0x70, 0x69, 0x70, 0x6c, 0x69, 0x6e, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x23, 0x0a, 0x0d,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0xab, 0x04, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x73,
	0x6b, 0x69, 0x70, 0x5f, 0x72, 0x75, 0x6e, 0x5f, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x73, 0x6b, 0x69, 0x70, 0x52, 0x75, 0x6e, 0x50,
	0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x6d, 0x6f,
	0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4d, 0x6f, 0x63, 0x6b,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x66, 0x72, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73,
	0x65, 0x72, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x44, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x69, 0x6e, 0x66, 0x72,
	0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x2e, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x56, 0x49,
	0x44, 0x45, 0x52, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x29, 0x0a,
	0x10, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08,
	0x73, 0x75, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x75, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x21, 0x0a,
	0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x10, 0x0a, 0x03, 0x72, 0x61, 0x77, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72,
	0x61, 0x77, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x72, 0x73, 0x65, 0x5f, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x73, 0x65, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x12, 0x40, 0x0a, 0x05, 0x65, 0x78, 0x74, 0x72, 0x61, 0x18, 0x0f, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e,
	0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05,
	0x65, 0x78, 0x74, 0x72, 0x61, 0x1a, 0x38, 0x0a, 0x0a, 0x45, 0x78, 0x74, 0x72, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a,
	0x56, 0x0a, 0x05, 0x53, 0x54, 0x41, 0x47, 0x45, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x45, 0x4e, 0x44,
	0x44, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x4e, 0x51, 0x55, 0x45, 0x55,
	0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x02,
	0x12, 0x0c, 0x0a, 0x08, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x12, 0x0b,
	0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x05, 0x2a, 0x1c, 0x0a, 0x0e, 0x45, 0x56, 0x45, 0x4e, 0x54,
	0x5f, 0x50, 0x52, 0x4f, 0x56, 0x49, 0x44, 0x45, 0x52, 0x12, 0x0a, 0x0a, 0x06, 0x47, 0x49, 0x54,
	0x4c, 0x41, 0x42, 0x10, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x6d,
	0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mflow_apps_trigger_pb_event_proto_rawDescOnce sync.Once
	file_mflow_apps_trigger_pb_event_proto_rawDescData = file_mflow_apps_trigger_pb_event_proto_rawDesc
)

func file_mflow_apps_trigger_pb_event_proto_rawDescGZIP() []byte {
	file_mflow_apps_trigger_pb_event_proto_rawDescOnce.Do(func() {
		file_mflow_apps_trigger_pb_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_mflow_apps_trigger_pb_event_proto_rawDescData)
	})
	return file_mflow_apps_trigger_pb_event_proto_rawDescData
}

var file_mflow_apps_trigger_pb_event_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_mflow_apps_trigger_pb_event_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_mflow_apps_trigger_pb_event_proto_goTypes = []interface{}{
	(STAGE)(0),                // 0: infraboard.mflow.trigger.STAGE
	(EVENT_PROVIDER)(0),       // 1: infraboard.mflow.trigger.EVENT_PROVIDER
	(*RecordSet)(nil),         // 2: infraboard.mflow.trigger.RecordSet
	(*Record)(nil),            // 3: infraboard.mflow.trigger.Record
	(*BuildStatus)(nil),       // 4: infraboard.mflow.trigger.BuildStatus
	(*Event)(nil),             // 5: infraboard.mflow.trigger.Event
	nil,                       // 6: infraboard.mflow.trigger.Event.ExtraEntry
	(*build.BuildConfig)(nil), // 7: infraboard.mflow.build.BuildConfig
	(*task.PipelineTask)(nil), // 8: infraboard.mflow.task.PipelineTask
}
var file_mflow_apps_trigger_pb_event_proto_depIdxs = []int32{
	3, // 0: infraboard.mflow.trigger.RecordSet.items:type_name -> infraboard.mflow.trigger.Record
	5, // 1: infraboard.mflow.trigger.Record.event:type_name -> infraboard.mflow.trigger.Event
	4, // 2: infraboard.mflow.trigger.Record.build_status:type_name -> infraboard.mflow.trigger.BuildStatus
	7, // 3: infraboard.mflow.trigger.BuildStatus.build_config:type_name -> infraboard.mflow.build.BuildConfig
	0, // 4: infraboard.mflow.trigger.BuildStatus.stage:type_name -> infraboard.mflow.trigger.STAGE
	8, // 5: infraboard.mflow.trigger.BuildStatus.pipline_task:type_name -> infraboard.mflow.task.PipelineTask
	1, // 6: infraboard.mflow.trigger.Event.provider:type_name -> infraboard.mflow.trigger.EVENT_PROVIDER
	6, // 7: infraboard.mflow.trigger.Event.extra:type_name -> infraboard.mflow.trigger.Event.ExtraEntry
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_mflow_apps_trigger_pb_event_proto_init() }
func file_mflow_apps_trigger_pb_event_proto_init() {
	if File_mflow_apps_trigger_pb_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mflow_apps_trigger_pb_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordSet); i {
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
		file_mflow_apps_trigger_pb_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
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
		file_mflow_apps_trigger_pb_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildStatus); i {
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
		file_mflow_apps_trigger_pb_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
			RawDescriptor: file_mflow_apps_trigger_pb_event_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mflow_apps_trigger_pb_event_proto_goTypes,
		DependencyIndexes: file_mflow_apps_trigger_pb_event_proto_depIdxs,
		EnumInfos:         file_mflow_apps_trigger_pb_event_proto_enumTypes,
		MessageInfos:      file_mflow_apps_trigger_pb_event_proto_msgTypes,
	}.Build()
	File_mflow_apps_trigger_pb_event_proto = out.File
	file_mflow_apps_trigger_pb_event_proto_rawDesc = nil
	file_mflow_apps_trigger_pb_event_proto_goTypes = nil
	file_mflow_apps_trigger_pb_event_proto_depIdxs = nil
}

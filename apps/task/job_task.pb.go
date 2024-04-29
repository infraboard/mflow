// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.0
// source: mflow/apps/task/pb/job_task.proto

package task

import (
	job "github.com/infraboard/mflow/apps/job"
	pipeline "github.com/infraboard/mflow/apps/pipeline"
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
	// 调度中,或者队列中
	STAGE_SCHEDULING STAGE = 4
	// 创建中, Job创建时更新
	STAGE_CREATING STAGE = 8
	// 运行中, Job对应的Pod创建完成后更新, 可以查看Pod日志
	STAGE_ACTIVE STAGE = 12
	// 取消中
	STAGE_CANCELING STAGE = 20
	// 任务被取消
	STAGE_CANCELED STAGE = 30
	// 运行失败
	STAGE_FAILED STAGE = 34
	// 运行成功
	STAGE_SUCCEEDED STAGE = 38
	// 跳过执行, 比如Pipeline中的一些可选流程
	STAGE_SKIPPED STAGE = 45
)

// Enum value maps for STAGE.
var (
	STAGE_name = map[int32]string{
		0:  "PENDDING",
		4:  "SCHEDULING",
		8:  "CREATING",
		12: "ACTIVE",
		20: "CANCELING",
		30: "CANCELED",
		34: "FAILED",
		38: "SUCCEEDED",
		45: "SKIPPED",
	}
	STAGE_value = map[string]int32{
		"PENDDING":   0,
		"SCHEDULING": 4,
		"CREATING":   8,
		"ACTIVE":     12,
		"CANCELING":  20,
		"CANCELED":   30,
		"FAILED":     34,
		"SUCCEEDED":  38,
		"SKIPPED":    45,
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
	return file_mflow_apps_task_pb_job_task_proto_enumTypes[0].Descriptor()
}

func (STAGE) Type() protoreflect.EnumType {
	return &file_mflow_apps_task_pb_job_task_proto_enumTypes[0]
}

func (x STAGE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use STAGE.Descriptor instead.
func (STAGE) EnumDescriptor() ([]byte, []int) {
	return file_mflow_apps_task_pb_job_task_proto_rawDescGZIP(), []int{0}
}

type JobTaskSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 总数量
	// @gotags: bson:"total" json:"total"
	Total int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total" bson:"total"`
	// 清单
	// @gotags: bson:"items" json:"items"
	Items []*JobTask `protobuf:"bytes,2,rep,name=items,proto3" json:"items" bson:"items"`
}

func (x *JobTaskSet) Reset() {
	*x = JobTaskSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobTaskSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobTaskSet) ProtoMessage() {}

func (x *JobTaskSet) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobTaskSet.ProtoReflect.Descriptor instead.
func (*JobTaskSet) Descriptor() ([]byte, []int) {
	return file_mflow_apps_task_pb_job_task_proto_rawDescGZIP(), []int{0}
}

func (x *JobTaskSet) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *JobTaskSet) GetItems() []*JobTask {
	if x != nil {
		return x.Items
	}
	return nil
}

type JobTask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 元信息
	// @gotags: bson:",inline" json:"meta"
	Meta *Meta `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta" bson:",inline"`
	// task定义, job运行时参数
	// @gotags: bson:",inline" json:"spec"
	Spec *pipeline.Task `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec" bson:",inline"`
	// 任务当前状态
	// @gotags: bson:"status" json:"status"
	Status *JobTaskStatus `protobuf:"bytes,3,opt,name=status,proto3" json:"status" bson:"status"`
	// 关联Job
	// @gotags: bson:"job" json:"job"
	Job *job.Job `protobuf:"bytes,4,opt,name=job,proto3" json:"job" bson:"job"`
}

func (x *JobTask) Reset() {
	*x = JobTask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobTask) ProtoMessage() {}

func (x *JobTask) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobTask.ProtoReflect.Descriptor instead.
func (*JobTask) Descriptor() ([]byte, []int) {
	return file_mflow_apps_task_pb_job_task_proto_rawDescGZIP(), []int{1}
}

func (x *JobTask) GetMeta() *Meta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *JobTask) GetSpec() *pipeline.Task {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *JobTask) GetStatus() *JobTaskStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *JobTask) GetJob() *job.Job {
	if x != nil {
		return x.Job
	}
	return nil
}

type Meta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 创建时间
	// @gotags: bson:"create_at" json:"create_at"
	CreateAt int64 `protobuf:"varint,1,opt,name=create_at,json=createAt,proto3" json:"create_at" bson:"create_at"`
	// 更新时间
	// @gotags: bson:"update_at" json:"update_at"
	UpdateAt int64 `protobuf:"varint,2,opt,name=update_at,json=updateAt,proto3" json:"update_at" bson:"update_at"`
	// 更新人
	// @gotags: bson:"update_by" json:"update_by"
	UpdateBy string `protobuf:"bytes,3,opt,name=update_by,json=updateBy,proto3" json:"update_by" bson:"update_by"`
}

func (x *Meta) Reset() {
	*x = Meta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Meta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Meta) ProtoMessage() {}

func (x *Meta) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Meta.ProtoReflect.Descriptor instead.
func (*Meta) Descriptor() ([]byte, []int) {
	return file_mflow_apps_task_pb_job_task_proto_rawDescGZIP(), []int{2}
}

func (x *Meta) GetCreateAt() int64 {
	if x != nil {
		return x.CreateAt
	}
	return 0
}

func (x *Meta) GetUpdateAt() int64 {
	if x != nil {
		return x.UpdateAt
	}
	return 0
}

func (x *Meta) GetUpdateBy() string {
	if x != nil {
		return x.UpdateBy
	}
	return ""
}

type JobTaskStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 合并后的运行时参数(job默认参数, 系统默认参数, 用户传入参数)
	// @gotags: bson:"run_params" json:"run_params"
	RunParams *job.RunParamSet `protobuf:"bytes,12,opt,name=run_params,json=runParams,proto3" json:"run_params" bson:"run_params"`
	// 任务当前状态
	// @gotags: bson:"stage" json:"stage"
	Stage STAGE `protobuf:"varint,1,opt,name=stage,proto3,enum=infraboard.mflow.task.STAGE" json:"stage" bson:"stage"`
	// 任务开始时间
	// @gotags: bson:"start_at" json:"start_at"
	StartAt int64 `protobuf:"varint,2,opt,name=start_at,json=startAt,proto3" json:"start_at" bson:"start_at"`
	// 任务结束时间
	// @gotags: bson:"end_at" json:"end_at"
	EndAt int64 `protobuf:"varint,3,opt,name=end_at,json=endAt,proto3" json:"end_at" bson:"end_at"`
	// 状态描述
	// @gotags: bson:"message" json:"message"
	Message string `protobuf:"bytes,4,opt,name=message,proto3" json:"message" bson:"message"`
	// 任务状态详细描述, 用于Debug
	// @gotags: bson:"detail" json:"detail"
	Detail string `protobuf:"bytes,5,opt,name=detail,proto3" json:"detail" bson:"detail"`
	// Job Task运行时环境变量, 大写开头的变量会更新到pipline中, 注入到后续执行的任务中
	// 小写开头的变量, 作为Task运行后的输出, task自己保留
	// @gotags: bson:"runtime_envs" json:"runtime_envs"
	RuntimeEnvs *job.RunParamSet `protobuf:"bytes,6,opt,name=runtime_envs,json=runtimeEnvs,proto3" json:"runtime_envs" bson:"runtime_envs"`
	// 任务运行后的产物信息, 用于界面展示
	// @gotags: bson:"markdown_output" json:"markdown_output"
	MarkdownOutput string `protobuf:"bytes,7,opt,name=markdown_output,json=markdownOutput,proto3" json:"markdown_output" bson:"markdown_output"`
	// 任务所需的临时资源
	// @gotags: bson:"temporary_resources" json:"temporary_resources"
	TemporaryResources []*TemporaryResource `protobuf:"bytes,8,rep,name=temporary_resources,json=temporaryResources,proto3" json:"temporary_resources" bson:"temporary_resources"`
	// 任务事件
	// @gotags: bson:"events" json:"events"
	Events []*pipeline.Event `protobuf:"bytes,9,rep,name=events,proto3" json:"events" bson:"events"`
	// 事件回调状态
	// @gotags: bson:"webhooks" json:"webhooks"
	Webhooks []*pipeline.WebHook `protobuf:"bytes,10,rep,name=webhooks,proto3" json:"webhooks" bson:"webhooks"`
	// Jot Task 关注人通知状态
	// @gotags: bson:"mention_users" json:"mention_users"
	MentionUsers []*pipeline.MentionUser `protobuf:"bytes,11,rep,name=mention_users,json=mentionUsers,proto3" json:"mention_users" bson:"mention_users"`
	// 运行时需不需审核
	// @gotags: bson:"audit" json:"audit"
	Audit *pipeline.Audit `protobuf:"bytes,13,opt,name=audit,proto3" json:"audit" bson:"audit"`
	// 状态扩展属性
	// @gotags: bson:"extension" json:"extension"
	Extension map[string]string `protobuf:"bytes,15,rep,name=extension,proto3" json:"extension" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" bson:"extension"`
}

func (x *JobTaskStatus) Reset() {
	*x = JobTaskStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobTaskStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobTaskStatus) ProtoMessage() {}

func (x *JobTaskStatus) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobTaskStatus.ProtoReflect.Descriptor instead.
func (*JobTaskStatus) Descriptor() ([]byte, []int) {
	return file_mflow_apps_task_pb_job_task_proto_rawDescGZIP(), []int{3}
}

func (x *JobTaskStatus) GetRunParams() *job.RunParamSet {
	if x != nil {
		return x.RunParams
	}
	return nil
}

func (x *JobTaskStatus) GetStage() STAGE {
	if x != nil {
		return x.Stage
	}
	return STAGE_PENDDING
}

func (x *JobTaskStatus) GetStartAt() int64 {
	if x != nil {
		return x.StartAt
	}
	return 0
}

func (x *JobTaskStatus) GetEndAt() int64 {
	if x != nil {
		return x.EndAt
	}
	return 0
}

func (x *JobTaskStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *JobTaskStatus) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *JobTaskStatus) GetRuntimeEnvs() *job.RunParamSet {
	if x != nil {
		return x.RuntimeEnvs
	}
	return nil
}

func (x *JobTaskStatus) GetMarkdownOutput() string {
	if x != nil {
		return x.MarkdownOutput
	}
	return ""
}

func (x *JobTaskStatus) GetTemporaryResources() []*TemporaryResource {
	if x != nil {
		return x.TemporaryResources
	}
	return nil
}

func (x *JobTaskStatus) GetEvents() []*pipeline.Event {
	if x != nil {
		return x.Events
	}
	return nil
}

func (x *JobTaskStatus) GetWebhooks() []*pipeline.WebHook {
	if x != nil {
		return x.Webhooks
	}
	return nil
}

func (x *JobTaskStatus) GetMentionUsers() []*pipeline.MentionUser {
	if x != nil {
		return x.MentionUsers
	}
	return nil
}

func (x *JobTaskStatus) GetAudit() *pipeline.Audit {
	if x != nil {
		return x.Audit
	}
	return nil
}

func (x *JobTaskStatus) GetExtension() map[string]string {
	if x != nil {
		return x.Extension
	}
	return nil
}

// 临时资源, 在Pipline允许结束时,调用释放
type TemporaryResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 资源的类型, 比如 configmap
	// @gotags: bson:"kind" json:"kind"
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind" bson:"kind"`
	// 资源的名字, 资源的集群和Namespace 由job维护
	// @gotags: bson:"name" json:"name"
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name" bson:"name"`
	// 资源的详情数据
	// @gotags: bson:"detail" json:"detail"
	Detail string `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail" bson:"detail"`
	// 创建时间
	// @gotags: bson:"create_at" json:"create_at"
	CreateAt int64 `protobuf:"varint,4,opt,name=create_at,json=createAt,proto3" json:"create_at" bson:"create_at"`
	// 创建时间
	// @gotags: bson:"release_at" json:"release_at"
	ReleaseAt int64 `protobuf:"varint,5,opt,name=release_at,json=releaseAt,proto3" json:"release_at" bson:"release_at"`
}

func (x *TemporaryResource) Reset() {
	*x = TemporaryResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemporaryResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemporaryResource) ProtoMessage() {}

func (x *TemporaryResource) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemporaryResource.ProtoReflect.Descriptor instead.
func (*TemporaryResource) Descriptor() ([]byte, []int) {
	return file_mflow_apps_task_pb_job_task_proto_rawDescGZIP(), []int{4}
}

func (x *TemporaryResource) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *TemporaryResource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TemporaryResource) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *TemporaryResource) GetCreateAt() int64 {
	if x != nil {
		return x.CreateAt
	}
	return 0
}

func (x *TemporaryResource) GetReleaseAt() int64 {
	if x != nil {
		return x.ReleaseAt
	}
	return 0
}

type RunTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 任务名称
	// @gotags: bson:"name" json:"name"
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name" bson:"name"`
	// job名称定义
	// @gotags: bson:"job_spec" json:"job_spec"
	JobSpec string `protobuf:"bytes,2,opt,name=job_spec,json=jobSpec,proto3" json:"job_spec" bson:"job_spec"`
	// 试运行, 并不会真正执行Job
	// @gotags: bson:"dry_run" json:"dry_run"
	DryRun bool `protobuf:"varint,3,opt,name=dry_run,json=dryRun,proto3" json:"dry_run" bson:"dry_run"`
	// 手动更新Job的状态, 默认由job runner的operator更新
	// @gotags: bson:"manual_update_status" json:"manual_update_status"
	ManualUpdateStatus bool `protobuf:"varint,4,opt,name=manual_update_status,json=manualUpdateStatus,proto3" json:"manual_update_status" bson:"manual_update_status"`
	// job运行时参数
	// @gotags: bson:"params" json:"params"
	Params *job.RunParamSet `protobuf:"bytes,5,opt,name=params,proto3" json:"params" bson:"params"`
	// 标签
	// @gotags: bson:"labels" json:"labels"
	Labels map[string]string `protobuf:"bytes,15,rep,name=labels,proto3" json:"labels" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" bson:"labels"`
}

func (x *RunTaskRequest) Reset() {
	*x = RunTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunTaskRequest) ProtoMessage() {}

func (x *RunTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mflow_apps_task_pb_job_task_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunTaskRequest.ProtoReflect.Descriptor instead.
func (*RunTaskRequest) Descriptor() ([]byte, []int) {
	return file_mflow_apps_task_pb_job_task_proto_rawDescGZIP(), []int{5}
}

func (x *RunTaskRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RunTaskRequest) GetJobSpec() string {
	if x != nil {
		return x.JobSpec
	}
	return ""
}

func (x *RunTaskRequest) GetDryRun() bool {
	if x != nil {
		return x.DryRun
	}
	return false
}

func (x *RunTaskRequest) GetManualUpdateStatus() bool {
	if x != nil {
		return x.ManualUpdateStatus
	}
	return false
}

func (x *RunTaskRequest) GetParams() *job.RunParamSet {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *RunTaskRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

var File_mflow_apps_task_pb_job_task_proto protoreflect.FileDescriptor

var file_mflow_apps_task_pb_job_task_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x74, 0x61, 0x73,
	0x6b, 0x2f, 0x70, 0x62, 0x2f, 0x6a, 0x6f, 0x62, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e,
	0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x1a, 0x1b, 0x6d, 0x66, 0x6c, 0x6f,
	0x77, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x2f, 0x70, 0x62, 0x2f, 0x6a, 0x6f,
	0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x61,
	0x70, 0x70, 0x73, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2f, 0x70, 0x62, 0x2f,
	0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x58,
	0x0a, 0x0a, 0x4a, 0x6f, 0x62, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x12, 0x34, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d,
	0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x4a, 0x6f, 0x62, 0x54, 0x61, 0x73,
	0x6b, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xda, 0x01, 0x0a, 0x07, 0x4a, 0x6f, 0x62,
	0x54, 0x61, 0x73, 0x6b, 0x12, 0x2f, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e,
	0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x52,
	0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x33, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x12, 0x3c, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x69, 0x6e, 0x66,
	0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x61,
	0x73, 0x6b, 0x2e, 0x4a, 0x6f, 0x62, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2b, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x4a, 0x6f, 0x62,
	0x52, 0x03, 0x6a, 0x6f, 0x62, 0x22, 0x5d, 0x0a, 0x04, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x1b, 0x0a,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x5f, 0x62, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x42, 0x79, 0x22, 0xc3, 0x06, 0x0a, 0x0d, 0x4a, 0x6f, 0x62, 0x54, 0x61, 0x73, 0x6b,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x40, 0x0a, 0x0a, 0x72, 0x75, 0x6e, 0x5f, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x6e, 0x66,
	0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x6a, 0x6f,
	0x62, 0x2e, 0x52, 0x75, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x53, 0x65, 0x74, 0x52, 0x09, 0x72,
	0x75, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x32, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e,
	0x53, 0x54, 0x41, 0x47, 0x45, 0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x65, 0x6e, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x65, 0x6e, 0x64, 0x41, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x12, 0x44, 0x0a, 0x0c, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x76, 0x73,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f,
	0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x52, 0x75,
	0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x53, 0x65, 0x74, 0x52, 0x0b, 0x72, 0x75, 0x6e, 0x74, 0x69,
	0x6d, 0x65, 0x45, 0x6e, 0x76, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x6d, 0x61, 0x72, 0x6b, 0x64, 0x6f,
	0x77, 0x6e, 0x5f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x6d, 0x61, 0x72, 0x6b, 0x64, 0x6f, 0x77, 0x6e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12,
	0x59, 0x0a, 0x13, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x72, 0x79, 0x5f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69,
	0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e,
	0x74, 0x61, 0x73, 0x6b, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x72, 0x79, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x12, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x72,
	0x79, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x06, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x6e, 0x66,
	0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x69,
	0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x3e, 0x0a, 0x08, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73,
	0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f,
	0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x2e, 0x57, 0x65, 0x62, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x08, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x73, 0x12, 0x4b, 0x0a, 0x0d, 0x6d, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x69, 0x6e,
	0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70,
	0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x4d, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x0c, 0x6d, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x73, 0x65, 0x72,
	0x73, 0x12, 0x36, 0x0a, 0x05, 0x61, 0x75, 0x64, 0x69, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66,
	0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x41, 0x75, 0x64,
	0x69, 0x74, 0x52, 0x05, 0x61, 0x75, 0x64, 0x69, 0x74, 0x12, 0x51, 0x0a, 0x09, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x69,
	0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e,
	0x74, 0x61, 0x73, 0x6b, 0x2e, 0x4a, 0x6f, 0x62, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x3c, 0x0a, 0x0e,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x8f, 0x01, 0x0a, 0x11, 0x54,
	0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6b, 0x69, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x12, 0x1b, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x41, 0x74, 0x22, 0xcb, 0x02, 0x0a,
	0x0e, 0x52, 0x75, 0x6e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x53, 0x70, 0x65, 0x63, 0x12, 0x17,
	0x0a, 0x07, 0x64, 0x72, 0x79, 0x5f, 0x72, 0x75, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x64, 0x72, 0x79, 0x52, 0x75, 0x6e, 0x12, 0x30, 0x0a, 0x14, 0x6d, 0x61, 0x6e, 0x75, 0x61,
	0x6c, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x6d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x0a, 0x06, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x6e, 0x66, 0x72,
	0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x6a, 0x6f, 0x62,
	0x2e, 0x52, 0x75, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x53, 0x65, 0x74, 0x52, 0x06, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x12, 0x49, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x0f,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x2e, 0x6d, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x52, 0x75, 0x6e,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a,
	0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x84, 0x01, 0x0a, 0x05, 0x53,
	0x54, 0x41, 0x47, 0x45, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x45, 0x4e, 0x44, 0x44, 0x49, 0x4e, 0x47,
	0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x43, 0x48, 0x45, 0x44, 0x55, 0x4c, 0x49, 0x4e, 0x47,
	0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x52, 0x45, 0x41, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x08,
	0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x0c, 0x12, 0x0d, 0x0a, 0x09,
	0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x49, 0x4e, 0x47, 0x10, 0x14, 0x12, 0x0c, 0x0a, 0x08, 0x43,
	0x41, 0x4e, 0x43, 0x45, 0x4c, 0x45, 0x44, 0x10, 0x1e, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x22, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x55, 0x43, 0x43, 0x45, 0x45, 0x44,
	0x45, 0x44, 0x10, 0x26, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x4b, 0x49, 0x50, 0x50, 0x45, 0x44, 0x10,
	0x2d, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x69, 0x6e, 0x66, 0x72, 0x61, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x6d, 0x66, 0x6c, 0x6f, 0x77,
	0x2f, 0x61, 0x70, 0x70, 0x73, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_mflow_apps_task_pb_job_task_proto_rawDescOnce sync.Once
	file_mflow_apps_task_pb_job_task_proto_rawDescData = file_mflow_apps_task_pb_job_task_proto_rawDesc
)

func file_mflow_apps_task_pb_job_task_proto_rawDescGZIP() []byte {
	file_mflow_apps_task_pb_job_task_proto_rawDescOnce.Do(func() {
		file_mflow_apps_task_pb_job_task_proto_rawDescData = protoimpl.X.CompressGZIP(file_mflow_apps_task_pb_job_task_proto_rawDescData)
	})
	return file_mflow_apps_task_pb_job_task_proto_rawDescData
}

var file_mflow_apps_task_pb_job_task_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_mflow_apps_task_pb_job_task_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_mflow_apps_task_pb_job_task_proto_goTypes = []interface{}{
	(STAGE)(0),                   // 0: infraboard.mflow.task.STAGE
	(*JobTaskSet)(nil),           // 1: infraboard.mflow.task.JobTaskSet
	(*JobTask)(nil),              // 2: infraboard.mflow.task.JobTask
	(*Meta)(nil),                 // 3: infraboard.mflow.task.Meta
	(*JobTaskStatus)(nil),        // 4: infraboard.mflow.task.JobTaskStatus
	(*TemporaryResource)(nil),    // 5: infraboard.mflow.task.TemporaryResource
	(*RunTaskRequest)(nil),       // 6: infraboard.mflow.task.RunTaskRequest
	nil,                          // 7: infraboard.mflow.task.JobTaskStatus.ExtensionEntry
	nil,                          // 8: infraboard.mflow.task.RunTaskRequest.LabelsEntry
	(*pipeline.Task)(nil),        // 9: infraboard.mflow.pipeline.Task
	(*job.Job)(nil),              // 10: infraboard.mflow.job.Job
	(*job.RunParamSet)(nil),      // 11: infraboard.mflow.job.RunParamSet
	(*pipeline.Event)(nil),       // 12: infraboard.mflow.pipeline.Event
	(*pipeline.WebHook)(nil),     // 13: infraboard.mflow.pipeline.WebHook
	(*pipeline.MentionUser)(nil), // 14: infraboard.mflow.pipeline.MentionUser
	(*pipeline.Audit)(nil),       // 15: infraboard.mflow.pipeline.Audit
}
var file_mflow_apps_task_pb_job_task_proto_depIdxs = []int32{
	2,  // 0: infraboard.mflow.task.JobTaskSet.items:type_name -> infraboard.mflow.task.JobTask
	3,  // 1: infraboard.mflow.task.JobTask.meta:type_name -> infraboard.mflow.task.Meta
	9,  // 2: infraboard.mflow.task.JobTask.spec:type_name -> infraboard.mflow.pipeline.Task
	4,  // 3: infraboard.mflow.task.JobTask.status:type_name -> infraboard.mflow.task.JobTaskStatus
	10, // 4: infraboard.mflow.task.JobTask.job:type_name -> infraboard.mflow.job.Job
	11, // 5: infraboard.mflow.task.JobTaskStatus.run_params:type_name -> infraboard.mflow.job.RunParamSet
	0,  // 6: infraboard.mflow.task.JobTaskStatus.stage:type_name -> infraboard.mflow.task.STAGE
	11, // 7: infraboard.mflow.task.JobTaskStatus.runtime_envs:type_name -> infraboard.mflow.job.RunParamSet
	5,  // 8: infraboard.mflow.task.JobTaskStatus.temporary_resources:type_name -> infraboard.mflow.task.TemporaryResource
	12, // 9: infraboard.mflow.task.JobTaskStatus.events:type_name -> infraboard.mflow.pipeline.Event
	13, // 10: infraboard.mflow.task.JobTaskStatus.webhooks:type_name -> infraboard.mflow.pipeline.WebHook
	14, // 11: infraboard.mflow.task.JobTaskStatus.mention_users:type_name -> infraboard.mflow.pipeline.MentionUser
	15, // 12: infraboard.mflow.task.JobTaskStatus.audit:type_name -> infraboard.mflow.pipeline.Audit
	7,  // 13: infraboard.mflow.task.JobTaskStatus.extension:type_name -> infraboard.mflow.task.JobTaskStatus.ExtensionEntry
	11, // 14: infraboard.mflow.task.RunTaskRequest.params:type_name -> infraboard.mflow.job.RunParamSet
	8,  // 15: infraboard.mflow.task.RunTaskRequest.labels:type_name -> infraboard.mflow.task.RunTaskRequest.LabelsEntry
	16, // [16:16] is the sub-list for method output_type
	16, // [16:16] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_mflow_apps_task_pb_job_task_proto_init() }
func file_mflow_apps_task_pb_job_task_proto_init() {
	if File_mflow_apps_task_pb_job_task_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mflow_apps_task_pb_job_task_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobTaskSet); i {
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
		file_mflow_apps_task_pb_job_task_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobTask); i {
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
		file_mflow_apps_task_pb_job_task_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Meta); i {
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
		file_mflow_apps_task_pb_job_task_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobTaskStatus); i {
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
		file_mflow_apps_task_pb_job_task_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TemporaryResource); i {
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
		file_mflow_apps_task_pb_job_task_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunTaskRequest); i {
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
			RawDescriptor: file_mflow_apps_task_pb_job_task_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mflow_apps_task_pb_job_task_proto_goTypes,
		DependencyIndexes: file_mflow_apps_task_pb_job_task_proto_depIdxs,
		EnumInfos:         file_mflow_apps_task_pb_job_task_proto_enumTypes,
		MessageInfos:      file_mflow_apps_task_pb_job_task_proto_msgTypes,
	}.Build()
	File_mflow_apps_task_pb_job_task_proto = out.File
	file_mflow_apps_task_pb_job_task_proto_rawDesc = nil
	file_mflow_apps_task_pb_job_task_proto_goTypes = nil
	file_mflow_apps_task_pb_job_task_proto_depIdxs = nil
}

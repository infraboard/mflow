package task

import (
	context "context"
	"fmt"
	"strings"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/websocket"
	"github.com/infraboard/mcenter/apps/policy"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/tools/pretty"
	job "github.com/infraboard/mflow/apps/job"
	pipeline "github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mpaas/common/terminal"
	"github.com/infraboard/mpaas/provider/k8s/workload"
	v1 "k8s.io/api/core/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

const (
	AppName = "tasks"
)

type Service interface {
	JobService
	PipelineService
}

type JobService interface {
	// 执行Job
	RunJob(context.Context, *pipeline.Task) (*JobTask, error)
	// 删除任务
	DeleteJobTask(context.Context, *DeleteJobTaskRequest) (*JobTask, error)
	// 任务Debug
	DebugJobTask(context.Context, *DebugJobTaskRequest)
	// 审核任务
	AuditJobTask(context.Context, *AuditJobTaskRequest) (*JobTask, error)

	JobRPCServer
}

func NewDebugJobTaskRequest(taskId string) *DebugJobTaskRequest {
	return &DebugJobTaskRequest{
		TaskId: taskId,
	}
}

type DebugJobTaskRequest struct {
	// 任务Id
	TaskId string `json:"task_id"`
	// 容器名称
	ContainerName string `json:"container_name"`
	// 终端配置
	Terminal terminal.TerminalSize `json:"terminal_size"`
	// Debug容器终端
	terminal *terminal.WebSocketTerminal
}

func (r *DebugJobTaskRequest) CopyPodRunRequest(namespace, podName string) *workload.CopyPodRunRequest {
	req := workload.NewCopyPodRunRequest()
	req.SourcePod.Name = podName
	req.SourcePod.Namespace = namespace
	req.TargetPodMeta.Namespace = namespace
	req.TargetPodMeta.Name = fmt.Sprintf("%s-d%s", req.SourcePod.Name, utilrand.String(5))
	req.ExecContainer = r.ContainerName
	req.ExecHoldCmd = workload.HoldContaienrCmd(1 * time.Hour)
	req.TerminationGracePeriodSeconds = 1
	return req
}

func (r *DebugJobTaskRequest) WebTerminal() *terminal.WebSocketTerminal {
	return r.terminal
}

func (r *DebugJobTaskRequest) SetWebTerminal(term *terminal.WebSocketTerminal) {
	r.terminal = term
	r.terminal.SetSize(r.Terminal)
}

func NewQueryTaskRequestFromHttp(r *restful.Request) *QueryJobTaskRequest {
	req := NewQueryTaskRequest()
	req.Page = request.NewPageRequestFromHTTP(r.Request)
	req.Scope = token.GetTokenFromRequest(r).GenScope()
	req.Filters = policy.GetScopeFilterFromRequest(r)
	req.JobId = r.QueryParameter("job_id")
	req.PipelineTaskId = r.QueryParameter("pipeline_task_id")
	req.Auditor = r.QueryParameter("auditor")
	auditEnable := r.QueryParameter("audit_enable")
	if auditEnable != "" {
		req.WithAuditEnable(auditEnable == "true")
	}
	auditStage := r.QueryParameter("audit_stages")
	if auditStage != "" {
		for _, stage := range strings.Split(auditStage, ",") {
			v, err := pipeline.ParseAUDIT_STAGEFromString(stage)
			if err != nil {
				continue
			}
			req.AuditStages = append(req.AuditStages, v)
		}
	}
	return req
}

func NewQueryTaskRequest() *QueryJobTaskRequest {
	return &QueryJobTaskRequest{
		Page:        request.NewDefaultPageRequest(),
		Ids:         []string{},
		AuditStages: []pipeline.AUDIT_STAGE{},
	}
}

func (req *QueryJobTaskRequest) HasLabel() bool {
	return req.Labels != nil && len(req.Labels) > 0
}

func (req *QueryJobTaskRequest) WithAuditEnable(v bool) *QueryJobTaskRequest {
	req.AuditEnable = &v
	return req
}

func (req *QueryJobTaskRequest) AddAuditStage(ss ...pipeline.AUDIT_STAGE) *QueryJobTaskRequest {
	req.AuditStages = append(req.AuditStages, ss...)
	return req
}

func NewRunTaskRequest(name, spec string, params *job.RunParamSet) *RunTaskRequest {
	return &RunTaskRequest{
		Name:    name,
		JobSpec: spec,
		Params:  params,
	}
}

// Job运行时需要的注解
func (r *RunTaskRequest) Annotations() map[string]string {
	annotations := map[string]string{}
	if !r.ManualUpdateStatus {
		annotations[ANNOTATION_TASK_ID] = r.Params.GetParamValue(job.SYSTEM_VARIABLE_JOB_TASK_ID)
	}
	return annotations
}

func (r *RunTaskRequest) RenderJobSpec() string {
	renderedSpec := r.JobSpec
	vars := r.Params.TemplateVars()
	for i := range vars {
		t := vars[i]
		renderedSpec = strings.ReplaceAll(renderedSpec, t.RefName(), t.Value)
	}
	return renderedSpec
}

func NewJobTaskEnvConfigMapName(TaskId string) string {
	return fmt.Sprintf("task-%s", TaskId)
}

// Job Task 挂载一个空的config 用于收集运行时的 环境变量
func (s *RunTaskRequest) RuntimeEnvConfigMap(mountPath string) *v1.ConfigMap {
	cm := new(v1.ConfigMap)
	s.Params.GetDeploymentId()
	cm.Name = NewJobTaskEnvConfigMapName(s.Params.GetJobTaskId())
	cm.BinaryData = map[string][]byte{
		CONFIG_MAP_RUNTIME_ENV_KEY: []byte(""),
	}
	cm.Annotations = map[string]string{
		workload.ANNOTATION_CONFIGMAP_MOUNT: mountPath,
	}
	return cm
}

type PipelineService interface {
	PipelineRPCServer
}

func NewDescribeJobTaskRequest(id string) *DescribeJobTaskRequest {
	return &DescribeJobTaskRequest{
		Id: id,
	}
}

func (r *DescribeJobTaskRequest) Validate() error {
	if r.Id == "" {
		return fmt.Errorf("job task id required")
	}
	return nil
}

func (r *DescribeJobTaskRequest) SetWithJob(v bool) *DescribeJobTaskRequest {
	r.WithJob = v
	return r
}

func NewDeleteJobTaskRequest(id string) *DeleteJobTaskRequest {
	return &DeleteJobTaskRequest{
		Id: id,
	}
}

func NewUpdateJobTaskOutputRequest(id string) *UpdateJobTaskOutputRequest {
	return &UpdateJobTaskOutputRequest{
		Id:          id,
		RuntimeEnvs: []*job.RunParam{},
	}
}

func NewAuditJobTaskRequest(id string) *AuditJobTaskRequest {
	return &AuditJobTaskRequest{
		TaskId: id,
		Status: pipeline.NewAuditStatus(),
	}
}

func (req *AuditJobTaskRequest) Validate() error {
	switch req.Status.Stage {
	case pipeline.AUDIT_STAGE_DENY:
		if req.Status.Comment == "" {
			return fmt.Errorf("审核失败, 需要填写失败原因")
		}
	}

	return nil
}

func (req *UpdateJobTaskOutputRequest) AddRuntimeEnv(name, value string) {
	req.RuntimeEnvs = append(req.RuntimeEnvs, job.NewRunParam(name, value))
}

func NewUpdateJobTaskStatusRequest(id string) *UpdateJobTaskStatusRequest {
	return &UpdateJobTaskStatusRequest{
		Id:        id,
		Extension: map[string]string{},
	}
}

func (r *UpdateJobTaskStatusRequest) MarkError(err error) {
	r.Stage = STAGE_FAILED
	r.Message = err.Error()
}

func NewQueryPipelineTaskRequest() *QueryPipelineTaskRequest {
	return &QueryPipelineTaskRequest{
		Page: request.NewDefaultPageRequest(),
	}
}

func NewQueryPipelineTaskRequestFromHttp(r *restful.Request) *QueryPipelineTaskRequest {
	req := NewQueryPipelineTaskRequest()
	req.Page = request.NewPageRequestFromHTTP(r.Request)
	req.Scope = token.GetTokenFromRequest(r).GenScope()
	req.Filters = policy.GetScopeFilterFromRequest(r)
	return req
}

func NewDeletePipelineTaskRequest(id string) *DeletePipelineTaskRequest {
	return &DeletePipelineTaskRequest{
		Id: id,
	}
}

func NewDescribePipelineTaskRequest(id string) *DescribePipelineTaskRequest {
	return &DescribePipelineTaskRequest{
		Id: id,
	}
}

// IM 消息通知, 独立适配
type TaskMessage interface {
	// 消息来自的源
	GetDomain() string
	// 消息来源的空间
	GetNamespace() string
	// IM消息事件的Title
	ShowTitle() string
	// Markdown格式的消息内容
	MarkdownContent() string
	// HTML格式的消息内容
	HTMLContent() string
	// 消息事件状态
	GetStatusStage() STAGE
	// 通知过程中的事件
	AddErrorEvent(format string, a ...any)
}

type WebHookMessage interface {
	TaskMessage
	// Web事件回调
	AddWebhookStatus(items ...*pipeline.WebHook)
}

// Task状态变更用户通知
type MentionUserMessage interface {

	// Task状态变化通知消息
	TaskMessage
	// 通知回调, 是否通知成功
	AddNotifyStatus(items ...*pipeline.MentionUser)
}

func NewJobTaskStreamReponse() *JobTaskStreamReponse {
	return &JobTaskStreamReponse{
		Data: make([]byte, 0, 512),
	}
}

func (r *JobTaskStreamReponse) ReSet() {
	r.Data = r.Data[:0]
}

func NewWatchJobTaskLogRequest(taskId string) *WatchJobTaskLogRequest {
	return &WatchJobTaskLogRequest{
		TaskId: taskId,
	}
}

func (req *WatchJobTaskLogRequest) ToJSON() string {
	return pretty.ToJSON(req)
}

func NewTaskLogWebsocketTerminal(conn *websocket.Conn) *TaskLogWebsocketTerminal {
	return &TaskLogWebsocketTerminal{
		terminal.NewWebSocketWriter(conn),
	}
}

type TaskLogWebsocketTerminal struct {
	*terminal.WebSocketWriter
}

func (r *TaskLogWebsocketTerminal) Send(in *JobTaskStreamReponse) (err error) {
	_, err = r.Write(in.Data)
	return
}

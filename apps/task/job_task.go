package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/infraboard/mcenter/apps/notify"
	"github.com/infraboard/mflow/apps/job"
	pipeline "github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/common/format"
	"github.com/infraboard/mpaas/provider/k8s/workload"
	core_v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func NewJobTaskSet() *JobTaskSet {
	return &JobTaskSet{
		Items: []*JobTask{},
	}
}

func (s *JobTaskSet) Add(tasks ...*JobTask) {
	s.Items = append(s.Items, tasks...)
}

func (s *JobTaskSet) Len() int {
	return len(s.Items)
}

// 继承来自Pipeline Task的属性
func (s *JobTaskSet) UpdateFromPipelineTask(p *PipelineTask) *JobTaskSet {
	for i := range s.Items {
		item := s.Items[i]
		item.Spec.RunParams.DryRun = p.Params.DryRun
		item.Spec.Domain = p.Params.Domain
		item.Spec.Namespace = p.Params.Namespace
	}
	return s
}

func (s *JobTaskSet) ToDocs() (docs []any) {
	for i := range s.Items {
		docs = append(docs, s.Items[i])
	}
	return
}

func (s *JobTaskSet) TaskNames() (names []string) {
	for i := range s.Items {
		names = append(names, s.Items[i].Spec.TaskName)
	}
	return
}

func (s *JobTaskSet) HasStage(stage STAGE) bool {
	for i := range s.Items {
		item := s.Items[i]
		if item.Status != nil && item.Status.Stage.Equal(stage) {
			return true
		}
	}
	return false
}

// 查询Stage中 等待执行的Job Task
func (s *JobTaskSet) GetJobTaskByStage(stage STAGE) (tasks []*JobTask) {
	for i := range s.Items {
		item := s.Items[i]
		if item.Status != nil && item.Status.Stage.Equal(stage) {
			tasks = append(tasks, item)
		}
	}
	return
}

func NewDefaultJobTask() *JobTask {
	req := pipeline.NewTask("")
	return NewJobTask(req)
}

func NewMeta() *Meta {
	return &Meta{
		CreateAt: time.Now().Unix(),
	}
}

func NewJobTask(req *pipeline.Task) *JobTask {
	req.SetDefault()
	t := &JobTask{
		Meta:   NewMeta(),
		Spec:   req,
		Job:    nil,
		Status: NewJobTaskStatus(),
	}

	if t.Spec.SkipRun {
		t.Status.Stage = STAGE_SKIPPED
		t.Status.Message = "skip run"
	}
	return t
}

var (
	// 关于Go模版语法可以参考: https://www.tizi365.com/archives/85.html
	JOB_TASK_MARKDOWN_TEMPLATE = `
**开始时间: **
{{ .Status.StartAtFormat }}
**结束时间: **
{{ .Status.EndAtAtFormat }}
**任务参数: **
{{ range .Spec.RunParams.Params -}}
▫ *{{.Name}}:  {{.Value}}*
{{end}}
`
)

func (p *JobTask) MarkdownContent() string {
	buf := bytes.NewBuffer([]byte{})
	t := template.New("job task")
	tmpl, err := t.Parse(JOB_TASK_MARKDOWN_TEMPLATE)
	if err != nil {
		return err.Error()
	}

	err = tmpl.Execute(buf, p)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

var (
	// 关于Go模版语法可以参考: https://www.tizi365.com/archives/85.html
	JOB_TASK_HTML_TEMPLATE = `
开始时间: 
{{ .Status.StartAtFormat }}
结束时间: 
{{ .Status.EndAtAtFormat }}
任务参数: 
{{ range .Spec.RunParams.Params -}}
▫ *{{.Name}}:  {{.Value}}*
{{end}}
`
)

func (p *JobTask) HTMLContent() string {
	buf := bytes.NewBuffer([]byte{})
	t := template.New("job task")
	tmpl, err := t.Parse(JOB_TASK_HTML_TEMPLATE)
	if err != nil {
		return err.Error()
	}

	err = tmpl.Execute(buf, p)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func (p *JobTask) GetDomain() string {
	return p.Spec.Domain
}

func (p *JobTask) GetNamespace() string {
	return p.Spec.Namespace
}

func (p *JobTask) AddNotifyStatus(items ...*pipeline.MentionUser) {
	p.Status.AddNotifyStatus(items...)
}

func (p *JobTask) AddErrorEvent(format string, a ...any) {
	p.Status.AddErrorEvent(format, a...)
}

func (p *JobTask) AddWebhookStatus(items ...*pipeline.WebHook) {
	p.Status.AddWebhookStatus(items...)
}

func (p *JobTask) BuildSearchLabel() {
	if p.Job != nil && p.Job.Spec != nil {
		if p.Job.Spec.Labels == nil {
			p.Job.Spec.Labels = map[string]string{}
		}
		for k, v := range p.Job.Spec.Labels {
			p.Spec.Labels[k] = v
		}
	}

	p.Spec.BuildSearchLabel()
}

func (p *JobTask) GetRunParamSet() *job.RunParamSet {
	return p.Spec.RunParams
}

func (p *JobTask) GetStatusRunParam() []*job.RunParam {
	if p.Status == nil ||
		p.Status.RunParams == nil ||
		p.Status.RunParams.Params == nil {
		return []*job.RunParam{}
	}
	return p.Status.RunParams.Params
}

func (p *JobTask) Desense() {
	if p.Status != nil && p.Status.RunParams != nil {
		p.Status.RunParams.Densense()
	}
	if p.Job != nil && p.Job.Spec.RunParams != nil {
		p.Job.Spec.RunParams.Densense()
	}
}

func (p *JobTask) SystemRunParam() (items []*job.RunParam) {
	items = append(items,
		job.NewRunParam(
			job.SYSTEM_VARIABLE_JOB_TASK_ID,
			p.Spec.TaskId,
		).SetReadOnly(true).
			SetRequired(true),
		job.NewRunParam(
			job.SYSTEM_VARIABLE_JOB_TASK_UPDATE_TOKEN,
			p.Spec.UpdateToken,
		).SetReadOnly(true).
			SetRequired(true),
	)

	if p.Job != nil {
		items = append(items,
			job.NewRunParam(
				job.SYSTEM_VARIABLE_JOB_ID,
				p.Job.Meta.Id,
			).SetReadOnly(true).
				SetRequired(true),
		)
	}

	if p.Spec.PipelineTask != "" {
		items = append(items,
			job.NewRunParam(
				job.SYSTEM_VARIABLE_PIPELINE_TASK_ID,
				p.Spec.PipelineTask,
			).SetReadOnly(true).
				SetRequired(true),
		)
	}
	return
}

func (p *JobTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*pipeline.Task
		*Meta
		Status *JobTaskStatus `json:"status"`
		Job    *job.Job       `json:"job"`
	}{p.Spec, p.Meta, p.Status, p.Job})
}

func (p *JobTask) RuntimeRunParams() []*job.RunParam {
	params := []*job.RunParam{}
	if p == nil ||
		p.Status == nil ||
		p.Status.RuntimeEnvs == nil ||
		p.Status.RuntimeEnvs.Params == nil {
		return params
	}
	return p.Status.RuntimeEnvs.Params
}

func (t *JobTask) GetStatusDetail() string {
	if t.Status != nil {
		return t.Status.Detail
	}

	return ""
}

func (s *JobTask) GetStatusStage() STAGE {
	return s.Status.GetStage()
}

func (s *JobTask) HasJobSpec() bool {
	if s.Job != nil && s.Job.Spec != nil {
		return true
	}

	return false
}

func (t *JobTask) WorkLoad() (*workload.WorkLoad, error) {
	if t.Status == nil {
		return nil, fmt.Errorf("")
	}
	switch t.Job.Spec.RunnerType {
	case job.RUNNER_TYPE_K8S_JOB:
		return workload.ParseWorkloadFromYaml("Job", t.Status.Detail)
	default:
		return nil, fmt.Errorf("unkonw runner type :%s", t.Job.Spec.RunnerType)
	}
}

func (t *JobTask) Update(job *job.Job, status *JobTaskStatus) {
	t.Job = job
	t.Status = status
}

func (t *JobTask) ToJson() string {
	return format.Prettify(t)
}

func (t *JobTask) ValidateToken(token string) error {
	if t.Spec.UpdateToken != token {
		return fmt.Errorf("update token invalidate")
	}
	return nil
}

func (s *JobTask) ShowTitle() string {
	return fmt.Sprintf("任务[%s]当前状态: %s", s.Spec.JobName, s.Status.Stage.String())
}

func (t *JobTask) AuditPass() bool {
	// 没有开启 直接通过
	if !t.Spec.Audit.Enable {
		return true
	}

	return t.IsAuditStatus(pipeline.AUDIT_STAGE_PASS)
}

// 兼容空数据
func (t *JobTask) AuditStatus() pipeline.AUDIT_STAGE {
	if t.Status.Audit == nil || t.Status.Audit.Status == nil {
		return pipeline.AUDIT_STAGE_PENDDING
	}
	return t.Status.Audit.Status.Stage
}

// 兼容空数据
func (t *JobTask) IsAuditStatus(stage pipeline.AUDIT_STAGE) bool {
	return t.AuditStatus().Equal(stage)
}

// 兼容空数据
func (t *JobTask) SetAuditStatus(stage pipeline.AUDIT_STAGE) {
	if t.Status.Audit == nil {
		t.Status.Audit = pipeline.NewAudit()
	}
	if t.Status.Audit.Status == nil {
		t.Status.Audit.Status = pipeline.NewAuditStatus()
	}
	t.Status.Audit.Status.Stage = stage
}

// 兼容空数据
func (t *JobTask) CheckUpdateStage(updatedStage pipeline.AUDIT_STAGE) error {
	current := t.AuditStatus()
	if updatedStage <= current {
		return fmt.Errorf("更新状态[%s], 落后或者等于当前状态[%s]", updatedStage, current)
	}
	return nil
}

func NewJobTaskStatus() *JobTaskStatus {
	return &JobTaskStatus{
		StartAt:            time.Now().Unix(),
		TemporaryResources: []*TemporaryResource{},
		Events:             []*pipeline.Event{},
		ImRobotNotify:      []*pipeline.WebHook{},
		Webhooks:           []*pipeline.WebHook{},
		MentionUsers:       []*pipeline.MentionUser{},
		Audit:              pipeline.NewAudit(),
		Extension:          map[string]string{},
	}
}

func (t *JobTaskStatus) MessageToError() error {
	return fmt.Errorf(t.Message)
}

func (t *JobTaskStatus) MarkedError(err error) {
	t.EndAt = time.Now().Unix()
	t.Stage = STAGE_FAILED
	t.Message = err.Error()
}

func (t *JobTaskStatus) MarkedCreating() {
	t.StartAt = time.Now().Unix()
	t.Stage = STAGE_CREATING
}

func (t *JobTaskStatus) MarkedSuccess() {
	t.Stage = STAGE_SUCCEEDED
	t.EndAt = time.Now().Unix()
}

func (t *JobTaskStatus) StartAtFormat() string {
	start := time.Unix(t.StartAt, 0)
	return start.Format("2006-01-02 03:04:05")
}

func (t *JobTaskStatus) EndAtAtFormat() string {
	start := time.Unix(t.EndAt, 0)
	return start.Format("2006-01-02 03:04:05")
}

func (t *JobTaskStatus) IsComplete() bool {
	return t.Stage >= STAGE_CANCELED
}

func (p *JobTaskStatus) AddWebhookStatus(items ...*pipeline.WebHook) {
	p.Webhooks = append(p.Webhooks, items...)
}

func (p *JobTaskStatus) AddNotifyStatus(items ...*pipeline.MentionUser) {
	p.MentionUsers = append(p.MentionUsers, items...)
}

func (t *JobTaskStatus) UpdateStatus(req *UpdateJobTaskStatusRequest) {
	t.Stage = req.Stage
	t.Message = req.Message

	// 更新扩展属性
	for k, v := range req.Extension {
		t.Extension[k] = v
	}

	// 取消的任务 不需要更新detail详情
	if !t.Stage.Equal(STAGE_CANCELED) {
		t.Detail = req.Detail
	}

	// 如果没传递结束时间, 则自动生成结束时间
	if t.IsComplete() && t.EndAt == 0 {
		t.EndAt = time.Now().Unix()
	}
}

func (s *JobTaskStatus) PodCount() int {
	return len(s.PodKeys())
}

func (s *JobTaskStatus) GetOrNewPodKey(podName string) string {
	key, _ := s.GetPodKey(podName)
	if key != "" {
		return key
	}
	return TaskPodArrayKey(s.PodCount())
}

func (s *JobTaskStatus) PodKeys() (keys []string) {
	for k := range s.Extension {
		if strings.HasPrefix(k, EXTENSION_FOR_TASK_POD_ARRAY) {
			keys = append(keys, k)
		}
	}
	return
}

func (s *JobTaskStatus) GetLatestPodValue() string {
	keys := s.PodKeys()
	if len(keys) == 0 {
		return ""
	}

	sortKeys := sort.StringSlice(keys)
	sort.Sort(sortKeys)
	latestKey := len(sortKeys) - 1
	return s.Extension[sortKeys[latestKey]]
}

func (s *JobTaskStatus) GetLatestPod() (*core_v1.Pod, error) {
	podStr := s.GetLatestPodValue()
	if podStr == "" {
		return nil, nil
	}
	pod := core_v1.Pod{}
	err := yaml.Unmarshal([]byte(podStr), &pod)
	if err != nil {
		return nil, fmt.Errorf("unmarshal pod to obj error, %s", err)
	}
	return &pod, nil
}

func (s *JobTaskStatus) GetPodKey(podName string) (string, error) {
	for k, v := range s.Extension {
		if strings.HasPrefix(k, EXTENSION_FOR_TASK_POD_ARRAY) {
			pod := core_v1.Pod{}
			err := yaml.Unmarshal([]byte(v), &pod)
			if err != nil {
				return "", err
			}
			if pod.Name == podName {
				return k, nil
			}
		}
	}

	return "", nil
}

func (t *JobTaskStatus) UpdateOutput(req *UpdateJobTaskOutputRequest) {
	if req.MarkdownOutput != "" {
		t.MarkdownOutput = req.MarkdownOutput
	}
}

func (t *JobTaskStatus) AddTemporaryResource(items ...*TemporaryResource) {
	t.TemporaryResources = append(t.TemporaryResources, items...)
}

func (t *JobTaskStatus) AddErrorEvent(format string, a ...any) {
	t.Events = append(t.Events, pipeline.NewEvent(pipeline.EVENT_LEVEL_ERROR, fmt.Sprintf(format, a...)))
}

func (t *JobTaskStatus) AddEvent(level pipeline.EVENT_LEVEL, format string, a ...any) {
	t.Events = append(t.Events, pipeline.NewEvent(level, fmt.Sprintf(format, a...)))
}

func (t *JobTaskStatus) GetTemporaryResource(kind, name string) *TemporaryResource {
	for i := range t.TemporaryResources {
		tr := t.TemporaryResources[i]
		if tr.Kind == kind && tr.Name == name {
			return tr
		}
	}
	return nil
}

func NewTemporaryResource(kind, name string) *TemporaryResource {
	return &TemporaryResource{
		Kind:     kind,
		Name:     name,
		CreateAt: time.Now().Unix(),
	}
}

func (r *TemporaryResource) IsReleased() bool {
	return r.ReleaseAt != 0
}

func NewMentionUser(username string, nt notify.NOTIFY_TYPE) *pipeline.MentionUser {
	m := pipeline.NewMentionUser(username)
	m.AddEvent(STAGE_SUCCEEDED.String())
	m.AddNotifyType(nt)
	return m
}

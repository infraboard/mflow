package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/infraboard/mcenter/apps/notify"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/validator"
	"github.com/infraboard/mcube/v2/pb/resource"
	job "github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/common/format"
	"github.com/rs/xid"
	"sigs.k8s.io/yaml"
)

func NewPipelineSet() *PipelineSet {
	return &PipelineSet{
		Items: []*Pipeline{},
	}
}

func (s *PipelineSet) Add(item *Pipeline) {
	s.Items = append(s.Items, item)
}

func (s *PipelineSet) GetJobs() (ids []string) {
	m := map[string]struct{}{}
	for i := range s.Items {
		p := s.Items[i]
		for _, id := range p.GetJobs() {
			m[id] = struct{}{}
		}
	}
	for k := range m {
		ids = append(ids, k)
	}
	return
}

type TaskForEatchHandler func(*Task)

func (s *PipelineSet) TaskForEach(h TaskForEatchHandler) {
	for i := range s.Items {
		for stageIndex := range s.Items[i].Spec.Stages {
			stage := s.Items[i].Spec.Stages[stageIndex]
			for taskIndex := range stage.Tasks {
				h(stage.Tasks[taskIndex])
			}
		}
	}
}

func (s *PipelineSet) UpdateJobs(jobs ...*job.Job) {
	for i := range s.Items {
		p := s.Items[i]
		p.UpdateJobs(jobs...)
	}
}

func NewDefaultPipeline() *Pipeline {
	return &Pipeline{
		Spec: NewCreatePipelineRequest(),
	}
}

// New 新建一个部署配置
func New(req *CreatePipelineRequest) (*Pipeline, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := req.CheckStageDuplicate(); err != nil {
		return nil, err
	}

	req.BuildNumber()
	d := &Pipeline{
		Meta: resource.NewMeta(),
		Spec: req,
	}
	return d, nil
}

// 注入编号
func (req *CreatePipelineRequest) BuildNumber() {
	for m := range req.Stages {
		stage := req.Stages[m]
		stage.Number = int32(m) + 1
		for n := range stage.Tasks {
			t := stage.Tasks[n]
			t.Number = int32(n) + 1
			t.StageName = stage.Name
			t.StageNumber = stage.Number
			t.Job = nil
			t.Init()
		}
	}
}

func (req *CreatePipelineRequest) CheckStageDuplicate() error {
	m := map[string]int64{}
	for i := range req.Stages {
		stage := req.Stages[i]
		m[stage.Name]++
	}
	for k, v := range m {
		if v > 1 {
			return fmt.Errorf("阶段: %s 名称重复", k)
		}
	}
	return nil
}

func (p *Pipeline) Update(req *UpdatePipelineRequest) {
	p.Meta.UpdateAt = time.Now().Unix()
	p.Meta.UpdateBy = req.UpdateBy
	p.Spec = req.Spec
}

func (p *Pipeline) TriggerNext() bool {
	return !p.Spec.IsParallel
}

func (p *Pipeline) Patch(req *UpdatePipelineRequest) error {
	p.Meta.UpdateAt = time.Now().Unix()
	p.Meta.UpdateBy = req.UpdateBy
	return mergo.MergeWithOverwrite(p.Spec, req.Spec)
}

func (p *Pipeline) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*resource.Meta
		*CreatePipelineRequest
	}{p.Meta, p.Spec})
}

func (p *Pipeline) GetStage(name string) *Stage {
	if p.Spec == nil {
		return nil
	}

	for i := range p.Spec.Stages {
		stage := p.Spec.Stages[i]
		if stage.Name == name {
			return stage
		}
	}

	return nil
}

func (p *Pipeline) GetJobs() (ids []string) {
	if p.Spec == nil {
		return nil
	}

	for i := range p.Spec.Stages {
		stage := p.Spec.Stages[i]
		ids = append(ids, stage.JobIds()...)
	}

	return
}

func (p *Pipeline) UpdateJobs(jobs ...*job.Job) {
	for ij := range jobs {
		jobIns := jobs[ij]
		for i := range p.Spec.Stages {
			stage := p.Spec.Stages[i]
			for it := range stage.Tasks {
				task := stage.Tasks[it]
				if task.JobId == jobIns.Meta.Id {
					task.Job = jobIns
					break
				}
			}
		}
	}
}

func NewCreatePipelineRequestFromYAML(yml string) (*CreatePipelineRequest, error) {
	req := NewCreatePipelineRequest()

	err := yaml.Unmarshal([]byte(yml), req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func NewCreatePipelineRequest() *CreatePipelineRequest {
	return &CreatePipelineRequest{
		With:         []*job.RunParam{},
		Stages:       []*Stage{},
		Webhooks:     []*WebHook{},
		MentionUsers: []*MentionUser{},
		Labels:       map[string]string{},
	}
}

func (req *CreatePipelineRequest) UpdateFromToken(tk *token.Token) {
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace
	req.CreateBy = tk.UserId
}

func (req *CreatePipelineRequest) ToYAML() string {
	yml, err := yaml.Marshal(req)
	if err != nil {
		panic(err)
	}
	return string(yml)
}

func (req *CreatePipelineRequest) MatchedWebHooks(event string) []*WebHook {
	hooks := []*WebHook{}
	for i := range req.Webhooks {
		h := req.Webhooks[i]
		if h.IsMatch(event) {
			hooks = append(hooks, h)
		}
	}

	return hooks
}

func (req *CreatePipelineRequest) AddStage(stages ...*Stage) {
	req.Stages = append(req.Stages, stages...)
}

func NewTask(jobName string) *Task {
	return &Task{
		JobName:       jobName,
		RunParams:     job.NewRunParamSet(),
		Webhooks:      []*WebHook{},
		MentionUsers:  []*MentionUser{},
		ImRobotNotify: []*WebHook{},
		Audit:         NewAudit(),
		Extension:     map[string]string{},
		Labels:        map[string]string{},
	}
}

func (r *Task) UpdateFromToken(tk *token.Token) {
	if tk == nil {
		return
	}

	r.Domain = tk.Domain
	r.Namespace = tk.Namespace
	r.RunBy = tk.UserId
}

func (r *Task) VersionName(version string) string {
	if strings.Contains(r.JobName, job.UNIQ_VERSION_SPLITER) {
		return r.JobName
	}
	return fmt.Sprintf("%s%s%s", r.JobName, job.UNIQ_VERSION_SPLITER, version)
}

func (r *Task) Enabled() bool {
	return !r.SkipRun
}

func (r *Task) IsAuditor(auditor string) bool {
	for _, aid := range r.Audit.Auditors {
		if aid == auditor {
			return true
		}
	}

	return false
}

func (r *Task) Init() {
	if r.Audit == nil {
		r.Audit = NewAudit()
	}
	if r.Webhooks == nil {
		r.Webhooks = []*WebHook{}
	}
	if r.ImRobotNotify == nil {
		r.ImRobotNotify = []*WebHook{}
	}
}

func (r *Task) StageNumberString() string {
	return fmt.Sprintf("%d", r.StageNumber)
}

func (r *Task) TaskNumberString() string {
	return fmt.Sprintf("%d", r.Number)
}

func (r *Task) MatchedWebHooks(event string) []*WebHook {
	hooks := []*WebHook{}
	for i := range r.Webhooks {
		h := r.Webhooks[i]
		if h.IsMatch(event) {
			hooks = append(hooks, h)
		}
	}

	return hooks
}

func (r *Task) MatchedImRobotNotify(event string) []*WebHook {
	hooks := []*WebHook{}
	for i := range r.ImRobotNotify {
		h := r.ImRobotNotify[i]
		if h.IsMatch(event) {
			hooks = append(hooks, h)
		}
	}

	return hooks
}

func (r *Task) AddWebhook(items ...*WebHook) {
	for i := range items {
		item := items[i]
		if !r.HasWebhook(item) {
			r.Webhooks = append(r.Webhooks, items...)
		}
	}
}

func (r *Task) HasWebhook(hook *WebHook) bool {
	for i := range r.Webhooks {
		if r.Webhooks[i].Url == hook.Url {
			return true
		}
	}
	return false
}

func (r *Task) AddMentionUser(items ...*MentionUser) {
	r.MentionUsers = append(r.MentionUsers, items...)
}

func (r *Task) BuildSearchLabel() {
	if r.Labels == nil {
		r.Labels = map[string]string{}
	}
	if r.RunParams == nil {
		return
	}

	lables := r.RunParams.SearchLabels()
	for k, v := range lables {
		r.Labels[k] = v
	}
}

func (r *Task) SetDefault() {
	if r.TaskId == "" {
		r.TaskId = "task-" + xid.New().String()
	}
	if r.UpdateToken == "" {
		r.UpdateToken = xid.New().String()
	}
	if r.RunParams == nil {
		r.RunParams = job.NewRunParamSet()
	}
	if r.RollbackParams == nil {
		r.RollbackParams = job.NewRunParamSet()
	}
	if r.Webhooks == nil {
		r.Webhooks = []*WebHook{}
	}
	if r.Labels == nil {
		r.Labels = map[string]string{}
	}
}

func (r *Task) GetJobShortName() string {
	nl := strings.Split(r.JobName, "@")
	if len(nl) > 0 && nl[0] != "" {
		return nl[0]
	}

	return r.JobName
}

func (req *RunPipelineRequest) UpdateFromToken(tk *token.Token) {
	if tk == nil {
		return
	}

	req.Domain = tk.Domain
	req.Namespace = tk.Namespace
	req.RunBy = tk.UserId
}

func (req *RunPipelineRequest) AddRunParam(params ...*job.RunParam) {
	req.RunParams = append(req.RunParams, params...)
}

func (req *RunPipelineRequest) GetTaskParams(stageNumber, taskNumber string) (
	params []*job.RunParam) {
	log := log.L()

	for i := range req.RunParams {
		param := req.RunParams[i]
		// scope 全局参数, 适配所有Stage
		if param.ParamScope == nil ||
			param.ParamScope.Stage == "" ||
			param.ParamScope.Stage == "*" ||
			param.ParamScope.Stage == "0" {
			params = append(params, param)
			continue
		}

		// 匹配具体的Stage
		if param.ParamScope.Stage == stageNumber {
			// Stage通用变量
			if param.ParamScope.Task == "" ||
				param.ParamScope.Task == "*" ||
				param.ParamScope.Task == "0" {
				params = append(params, param)
				continue
			}

			// Task专有变量
			if param.ParamScope.Task == taskNumber {
				log.Debug().Msgf("[task param] %s param scope: %s", param.Name, param.ParamScope)
				params = append(params, param)
				continue
			}
		}
	}
	return
}

func (req *RunPipelineRequest) Validate() error {
	return validator.Validate(req)
}

func (req *RunPipelineRequest) ToJson() string {
	return format.Prettify(req)
}

func (req *RunPipelineRequest) RunParamsKVMap() map[string]string {
	m := map[string]string{}
	for i := range req.RunParams {
		p := req.RunParams[i]
		m[p.Name] = p.Value
	}
	return m
}

func (req *RunPipelineRequest) MergeLabels(labels map[string]string, overwrite bool) {
	if req.Labels == nil {
		req.Labels = map[string]string{}
	}
	for k, v := range labels {
		if req.Labels[k] != "" && !overwrite {
			continue
		}
		req.Labels[k] = v
	}
}

func NewWebHook(url string) *WebHook {
	return &WebHook{
		Url:    url,
		Header: map[string]string{},
		Events: []string{},
		Scope:  job.NewParamScope(),
	}
}

func (h *WebHook) IsMatch(event string) bool {
	if len(h.Events) == 0 {
		return true
	}

	for _, e := range h.Events {
		if strings.EqualFold(e, event) {
			return true
		}
	}

	return false
}

// 显示名称
func (h *WebHook) ShowName() string {
	return fmt.Sprintf("%s[%s]", h.Description, h.Url)
}

func NewMentionUser(username string) *MentionUser {
	return &MentionUser{
		UserName:    username,
		Events:      []string{},
		NotifyTypes: []notify.NOTIFY_TYPE{},
		Scope:       job.NewParamScope(),
	}
}

func (m *MentionUser) AddEvent(events ...string) {
	m.Events = append(m.Events, events...)
}

func (m *MentionUser) AddNotifyType(nts ...notify.NOTIFY_TYPE) {
	m.NotifyTypes = append(m.NotifyTypes, nts...)
}

func (m *MentionUser) IsMatch(event string) bool {
	if len(m.Events) == 0 {
		return true
	}

	for _, e := range m.Events {
		if strings.EqualFold(e, event) {
			return true
		}
	}

	return false
}

func (s *Stage) JobIds() (ids []string) {
	for i := range s.Tasks {
		task := s.Tasks[i]
		ids = append(ids, task.JobId)
	}
	return
}

func (s *Stage) GetTaskByNumber(number int32) *Task {
	if s.Tasks == nil {
		return nil
	}

	for i := range s.Tasks {
		task := s.Tasks[i]
		if task.Number == number {
			return task
		}
	}

	return nil
}

func NewAuditStatus() *AuditStatus {
	return &AuditStatus{
		Extension: map[string]string{},
	}
}

func NewAudit() *Audit {
	return &Audit{
		Auditors: []string{},
		Status:   NewAuditStatus(),
		Scope:    job.NewParamScope(),
	}
}

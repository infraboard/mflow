package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/pb/resource"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/pipeline"
)

func NewPipelineTaskSet() *PipelineTaskSet {
	return &PipelineTaskSet{
		Items: []*PipelineTask{},
	}
}

func (s *PipelineTaskSet) Add(item *PipelineTask) {
	s.Items = append(s.Items, item)
}

func (s *PipelineTaskSet) Len() int {
	return len(s.Items)
}

func NewPipelineTask(p *pipeline.Pipeline, in *pipeline.RunPipelineRequest) *PipelineTask {
	pt := NewDefaultPipelineTask()
	pt.Pipeline = p
	pt.Params = in

	// 合并来自于Pipeline的标签
	pt.Params.MergeLabels(p.Spec.Labels, false)

	// 如果传入了id 则使用传入的id
	if in.PipelineTaskId != "" {
		pt.Meta.Id = in.PipelineTaskId
		in.PipelineTaskId = ""
	}

	// 初始化所有的JobTask
	for i := range p.Spec.Stages {
		spec := p.Spec.Stages[i]
		ss := NewStageStatus(spec, pt.Meta.Id)
		pt.Status.AddStage(ss)
	}

	return pt
}

func NewDefaultPipelineTask() *PipelineTask {
	return &PipelineTask{
		Meta:   resource.NewMeta(),
		Params: pipeline.NewRunPipelineRequest(""),
		Status: NewPipelineTaskStatus(),
	}
}

func (p *PipelineTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*resource.Meta
		*pipeline.RunPipelineRequest
		*PipelineTaskStatus
		Pipeline *pipeline.Pipeline `json:"pipeline"`
	}{p.Meta, p.Params, p.Status, p.Pipeline})
}

func (p *PipelineTask) IsActive() bool {
	if p.Status != nil && p.Status.Stage.Equal(STAGE_ACTIVE) {
		return true
	}

	return false
}

func (p *PipelineTask) HasJobSpec() bool {
	if p.Pipeline != nil && p.Pipeline.Spec != nil {
		return true
	}

	return false
}

func (p *PipelineTask) MarkedRunning() {
	p.Status.Stage = STAGE_ACTIVE
	p.Status.StartAt = time.Now().Unix()
}

func (p *PipelineTask) MarkedSucceed() {
	p.Status.Stage = STAGE_SUCCEEDED
	p.Status.EndAt = time.Now().Unix()
}

func (p *PipelineTask) GetFirstJobTask() *JobTask {
	for i := range p.Status.StageStatus {
		s := p.Status.StageStatus[i]
		if len(s.Tasks) > 0 {
			task := s.Tasks[0]
			task.Spec.RunParams.DryRun = p.Params.DryRun
			task.Spec.RollbackParams.DryRun = p.Params.DryRun
			task.Spec.Domain = p.Params.Domain
			task.Spec.Namespace = p.Params.Namespace
			return task
		}
	}
	return nil
}

func (p *PipelineTask) JobTasks() *JobTaskSet {
	set := NewJobTaskSet()
	if p.Status == nil {
		return set
	}

	return p.Status.JobTasks().UpdateFromPipelineTask(p)
}

// 返回下个需要执行的JobTask, 允许一次并行执行多个(批量执行)
// Task DryRun属性要继承PipelineTask
func (p *PipelineTask) NextRun() (*JobTaskSet, error) {
	set := NewJobTaskSet()
	var stage *StageStatus

	if p.Status == nil || p.Pipeline == nil {
		return set, nil
	}

	// 需要未执行完成的Job Tasks
	stages := p.Status.StageStatus
	for i := range stages {
		stage = stages[i]

		// 找出Stage中未执行完的Job Task
		set = stage.UnCompleteJobTask(p.Params.DryRun)
		set.UpdateFromPipelineTask(p)
		// 如果找到 直接Break
		if set.Len() > 0 {
			break
		}
	}

	// 如果所有Stage寻找完，都没找到, 表示PipelineTask执行完成
	if set.Len() == 0 {
		return set, nil
	}

	// TODO: 如果这些未执行当中的Job Task 有处于运行中的, 不会执行下个一个任务
	if set.HasStage(STAGE_ACTIVE) {
		return set, exception.NewConflict("Stage 还处于运行中")
	}

	// 当未执行的任务中，没有运行中的时，剩下的就是需要被执行的任务
	tasks := set.GetJobTaskByStage(STAGE_PENDDING)

	nextTasks := NewJobTaskSet()
	stageSpec := p.Pipeline.GetStage(stage.Name)
	if stageSpec.IsParallel {
		// 并行任务 返回该Stage所有等待执行的job
		nextTasks.Add(tasks...)
	} else {
		// 串行任务取第一个
		nextTasks.Add(tasks[0])
	}

	return nextTasks.UpdateFromPipelineTask(p), nil
}

func (p *PipelineTask) GetStage(name string) *StageStatus {
	if p.Status == nil || p.Pipeline == nil {
		return nil
	}

	stages := p.Status.StageStatus
	// 先查 StageStatus 是否有
	for i := range stages {
		stage := stages[i]
		if stage.Name == name {
			return stage
		}
	}

	// 如果没有动态创建
	stageSpec := p.Pipeline.GetStage(name)
	if stageSpec == nil {
		return nil
	}

	stage := NewStageStatus(stageSpec, p.Meta.Id)
	p.Status.AddStage(stage)

	return stage
}

func (p *PipelineTask) GetJobTask(id string) *JobTask {
	return p.Status.GetJobTask(id)
}

// Pipeline执行成功
func (p *PipelineTask) UUID() string {
	return fmt.Sprintf("jobtask-%s", p.Meta.Id)
}

// Pipeline执行成功
func (p *PipelineTask) MarkedSuccess() {
	p.Status.Stage = STAGE_SUCCEEDED
	p.Status.EndAt = time.Now().Unix()
}

// Pipeline执行失败
func (p *PipelineTask) MarkedFailed(err error) {
	p.Status.Stage = STAGE_FAILED
	p.Status.EndAt = time.Now().Unix()
	p.Status.Message = err.Error()
}

// Pipeline执行取消
func (p *PipelineTask) MarkedCanceled() {
	p.Status.Stage = STAGE_CANCELED
	p.Status.EndAt = time.Now().Unix()
}

// Pipeline执行取消
func (p *PipelineTask) IsComplete() bool {
	if p.Status != nil && p.Status.Stage >= STAGE_CANCELED {
		return true
	}

	return false
}

// Pipeline执行取消
func (p *PipelineTask) IsRunning() bool {
	if p.Status == nil {
		return false
	}

	return p.Status.Stage > STAGE_PENDDING && p.Status.Stage < STAGE_CANCELED
}

// IM消息的标题
func (p *PipelineTask) ShowTitle() string {
	return fmt.Sprintf("流水线[%s]当前状态: %s", p.Pipeline.Spec.Name, p.Status.Stage.String())
}

var (
	// 关于Go模版语法可以参考: https://www.tizi365.com/archives/85.html
	PIPELINE_TASK_MARKDOWN_TEMPLATE = `
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

func (p *PipelineTask) MarkdownContent() string {
	buf := bytes.NewBuffer([]byte{})
	t := template.New("pipeline task")
	tmpl, err := t.Parse(PIPELINE_TASK_MARKDOWN_TEMPLATE)
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
	PIPELINE_TASK_HTML_TEMPLATE = `
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

func (p *PipelineTask) HTMLContent() string {
	buf := bytes.NewBuffer([]byte{})
	t := template.New("pipeline task")
	tmpl, err := t.Parse(PIPELINE_TASK_HTML_TEMPLATE)
	if err != nil {
		return err.Error()
	}

	err = tmpl.Execute(buf, p)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func (p *PipelineTask) GetStatusStage() STAGE {
	return p.Status.GetStage()
}

func (p *PipelineTask) GetDomain() string {
	return p.Params.Domain
}

func (p *PipelineTask) GetNamespace() string {
	return p.Params.Namespace
}

func (p *PipelineTask) AddErrorEvent(format string, a ...any) {
	p.Status.AddErrorEvent(format, a...)
}

func (p *PipelineTask) AddSuccessEvent(format string, a ...any) *pipeline.Event {
	return p.Status.AddSuccessEvent(format, a...)
}

func (p *PipelineTask) GetBuildConfId() string {
	return p.Params.Labels[build.PIPELINE_TASK_BUILD_CONFIG_ID_LABLE_KEY]
}

func (p *PipelineTask) GetEventId() string {
	return p.Params.Labels[build.PIPELINE_TASK_EVENT_ID_LABLE_KEY]
}

func (p *PipelineTask) HasBuildEvent() bool {
	return p.GetEventId() != "" && p.GetBuildConfId() != ""
}

func (p *PipelineTask) AddWebhookStatus(items ...*pipeline.WebHook) {
	p.Status.AddWebhookStatus(items...)
}

func (p *PipelineTask) AddNotifyStatus(items ...*pipeline.MentionUser) {
	p.Status.AddNotifyStatus(items...)
}

// 大写导出
func (s *PipelineTask) RuntimeRunParams() (envs []*job.RunParam) {
	if s.Status == nil || s.Status.RuntimeEnvs == nil {
		return
	}
	return s.Status.RuntimeEnvs.Params
}

// 执行参数 回填回Pipeline, 方便前端直接展示
func (s *PipelineTask) SetParamValueToPipeline() {
	if s.Pipeline == nil || s.Status == nil {
		return
	}

	for stageIndex := range s.Status.StageStatus {
		stage := s.Status.StageStatus[stageIndex]
		for taskIndex := range stage.Tasks {
			task := stage.Tasks[taskIndex]
			targetTask := s.Pipeline.GetStage(stage.Name).GetTaskByNumber(task.Spec.Number)
			if targetTask == nil {
				continue
			}
			params := task.GetStatusRunParam()
			for paramIndex := range params {
				param := params[paramIndex]
				targetParm := targetTask.RunParams.GetParam(param.Name)
				if targetParm == nil {
					continue
				}
				targetParm.Value = param.Value
			}
		}
	}
}

func NewPipelineTaskStatus() *PipelineTaskStatus {
	return &PipelineTaskStatus{
		RuntimeEnvs:  job.NewRunParamSet(),
		StageStatus:  []*StageStatus{},
		Webhooks:     []*pipeline.WebHook{},
		MentionUsers: []*pipeline.MentionUser{},
		Events:       []*pipeline.Event{},
	}
}

func (p *PipelineTaskStatus) JobTasks() *JobTaskSet {
	set := NewJobTaskSet()
	for i := range p.StageStatus {
		stage := p.StageStatus[i]
		set.Add(stage.Tasks...)
	}
	return set
}

func (s *PipelineTaskStatus) AddStage(item *StageStatus) {
	s.StageStatus = append(s.StageStatus, item)
}

func (t *PipelineTaskStatus) AddErrorEvent(format string, a ...any) {
	t.Events = append(t.Events, pipeline.NewEvent(pipeline.EVENT_LEVEL_ERROR, fmt.Sprintf(format, a...)))
}

func (t *PipelineTaskStatus) AddSuccessEvent(format string, a ...any) *pipeline.Event {
	e := pipeline.NewEvent(pipeline.EVENT_LEVEL_INFO, fmt.Sprintf(format, a...))
	t.Events = append(t.Events, e)
	return e
}

func (t *PipelineTaskStatus) AddEvent(level pipeline.EVENT_LEVEL, format string, a ...any) {
	t.Events = append(t.Events, pipeline.NewEvent(level, fmt.Sprintf(format, a...)))
}

func (p *PipelineTaskStatus) AddWebhookStatus(items ...*pipeline.WebHook) {
	p.Webhooks = append(p.Webhooks, items...)
}

func (s *PipelineTaskStatus) GetJobTask(id string) *JobTask {
	for i := range s.StageStatus {
		stage := s.StageStatus[i]
		jobTask := stage.GetJobTask(id)
		if jobTask != nil {
			return jobTask
		}
	}

	return nil
}

func (p *PipelineTaskStatus) AddNotifyStatus(items ...*pipeline.MentionUser) {
	p.MentionUsers = append(p.MentionUsers, items...)
}

func NewStageStatus(s *pipeline.Stage, pipelineTaskId string) *StageStatus {
	status := &StageStatus{
		Name:  s.Name,
		Tasks: []*JobTask{},
	}

	for i := range s.Tasks {
		req := s.Tasks[i]
		req.PipelineTask = pipelineTaskId
		req.StageName = s.Name
		jobTask := NewJobTask(req)
		status.Add(jobTask)
	}

	return status
}

// 获取Stage中未完成的任务, 包含 运行中和等待执行的
func (s *StageStatus) UnCompleteJobTask(dryRun bool) *JobTaskSet {
	set := NewJobTaskSet()
	for i := range s.Tasks {
		item := s.Tasks[i]
		// stage中任何一个job task未完成, 该stage都处于未完成
		if item.Status != nil && !item.Status.IsComplete() {
			item.Spec.RunParams.DryRun = dryRun
			set.Add(item)
		}
	}

	return set
}

func (s *StageStatus) Add(item *JobTask) {
	s.Tasks = append(s.Tasks, item)
}

// 根据Job Task id获取当前stage中的Job Task
func (s *StageStatus) GetJobTask(id string) *JobTask {
	for i := range s.Tasks {
		item := s.Tasks[i]
		if item.Spec.TaskId == id {
			return item
		}
	}
	return nil
}

package trigger

import (
	context "context"
	"fmt"
	"strings"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcenter/apps/policy"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/pb/resource"
	"github.com/infraboard/mflow/apps/build"
)

const (
	AppName = "triggers"
)

type Service interface {
	RPCServer

	// 事件队列任务执行完成通知
	EventQueueTaskComplete(context.Context, *EventQueueTaskCompleteRequest) (*BuildStatus, error)
}

type EventQueueTaskComplete struct {
}

func NewGitlabWebHookEvent() *GitlabWebHookEvent {
	return &GitlabWebHookEvent{
		Project: &Project{},
		Commits: []*Commit{},
	}
}

func (e *GitlabWebHookEvent) ShortDesc() string {
	return fmt.Sprintf("%s %s [%s]", e.Ref, e.EventName, e.ObjectKind)
}

func (e *GitlabWebHookEvent) GetBranch() string {
	switch e.EventType {
	case GITLAB_EVENT_TYPE_MERGE_REQUEST:
		return e.ObjectAttributes.TargetBranch
	case GITLAB_EVENT_TYPE_PUSH:
		return strings.TrimPrefix(e.Ref, "refs/heads/")
	case GITLAB_EVENT_TYPE_TAG:
		return e.GetTag()
	default:
		return e.Ref
	}
}

func (e *GitlabWebHookEvent) GetTag() string {
	return strings.TrimPrefix(e.Ref, "refs/tags/")
}

func (e *GitlabWebHookEvent) GetLatestCommit() *Commit {
	count := len(e.Commits)
	if count > 0 {
		return e.Commits[count-1]
	}
	return nil
}

func (e *GitlabWebHookEvent) GetLatestCommitShortId() string {
	cm := e.GetLatestCommit()
	if cm != nil {
		return cm.Short()
	}
	return ""
}

func NewGitlabEvent(raw string) *Event {
	return &Event{
		Provider: EVENT_PROVIDER_GITLAB,
		Time:     time.Now().UnixMilli(),
		Raw:      raw,
	}
}

func NewBuildStatus(bc *build.BuildConfig) *BuildStatus {
	return &BuildStatus{
		BuildConfig: bc,
	}
}

func NewQueryRecordRequestFromHTTP(r *restful.Request) *QueryRecordRequest {
	req := NewQueryRecordRequest()
	req.Page = request.NewPageRequestFromHTTP(r.Request)
	req.Scope = token.GetTokenFromRequest(r).GenScope()
	req.Filters = policy.GetScopeFilterFromRequest(r)
	req.ServiceId = r.QueryParameter("service_id")
	req.PipelineTaskId = r.QueryParameter("pipeline_task_id")
	req.WithPipelineTask = r.QueryParameter("with_pipeline_task") == "true"

	bid := r.QueryParameter("build_conf_ids")
	if bid != "" {
		req.BuildConfIds = strings.Split(bid, ",")
	}

	return req
}

func NewQueryRecordRequest() *QueryRecordRequest {
	return &QueryRecordRequest{
		Page:         request.NewDefaultPageRequest(),
		Scope:        resource.NewScope(),
		Filters:      []*resource.LabelRequirement{},
		BuildConfIds: []string{},
	}
}

func NewEventQueueTaskCompleteRequest(pid string) *EventQueueTaskCompleteRequest {
	return &EventQueueTaskCompleteRequest{
		PipelineTaskId: pid,
	}
}

package trigger

import (
	context "context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/infraboard/mcube/v2/http/request"
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

func NewQueryRecordRequestFromHTTP(r *http.Request) *QueryRecordRequest {
	return &QueryRecordRequest{
		Page: request.NewPageRequestFromHTTP(r),
	}
}

func NewQueryRecordRequest() *QueryRecordRequest {
	return &QueryRecordRequest{
		Page: request.NewDefaultPageRequest(),
	}
}

func NewEventQueueTaskCompleteRequest(pid string) *EventQueueTaskCompleteRequest {
	return &EventQueueTaskCompleteRequest{
		PipelineTaskId: pid,
	}
}

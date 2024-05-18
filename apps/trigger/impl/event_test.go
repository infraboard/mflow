package impl_test

import (
	"testing"

	"github.com/infraboard/mflow/apps/trigger"
	"github.com/infraboard/mflow/test/tools"
)

func TestHandleEvent(t *testing.T) {
	raw := tools.MustReadContentFile("test/gitlab_push.json")
	req := trigger.NewGitlabEvent(raw)
	req.SkipRunPipeline = false

	ps, err := impl.HandleEvent(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ps))
}

func TestQueryRecord(t *testing.T) {
	req := trigger.NewQueryRecordRequest()
	req.AddBuildConfId("colfqdh97i61i9lcg4jg")
	req.AddBuildStage(trigger.STAGE_ENQUEUE)
	req.IsOrderAscend = true
	req.Page.PageSize = 1
	set, err := impl.QueryRecord(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(set))
}

func TestEventQueueTaskComplete(t *testing.T) {
	req := trigger.NewEventQueueTaskCompleteRequest("", "colfqdh97i61i9lcg4jg")
	req.TriggerNext = false
	set, err := impl.EventQueueTaskComplete(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(set))
}

func TestDeleteRecord(t *testing.T) {
	req := trigger.NewDeleteRecordRequest()
	req.DeleteBy = trigger.DELETE_BY_PIPELINE_TASK_ID
	req.AddValue("coig4d197i655m4akveg")
	err := impl.DeleteRecord(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
}

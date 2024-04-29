package webhook_test

import (
	"context"
	"testing"

	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/test/conf"
	"github.com/infraboard/mflow/test/tools"

	"github.com/infraboard/mflow/apps/task/webhook"
)

var (
	ctx    = context.Background()
	sender = webhook.NewWebHook()
	jt     = testPipelineStep()
)

func TestFeishuWebHook(t *testing.T) {
	hooks := testPipelineWebHook(conf.C.FEISHU_BOT_URL)
	sender.SendTaskStatus(ctx, hooks, jt)
	t.Log(jt.ToJson())
}

func TestDingDingWebHook(t *testing.T) {
	hooks := testPipelineWebHook(conf.C.DINGDING_BOT_URL)
	sender.SendTaskStatus(ctx, hooks, jt)
	t.Log(jt)
}

func TestWechatWebHook(t *testing.T) {
	hooks := testPipelineWebHook(conf.C.WECHAT_BOT_URL)
	sender.SendTaskStatus(ctx, hooks, jt)
	t.Log(jt)
}

func testPipelineWebHook(url string) []*pipeline.WebHook {
	h1 := &pipeline.WebHook{
		Url:         url,
		Events:      []string{task.STAGE_SUCCEEDED.String()},
		Description: "测试",
	}
	return []*pipeline.WebHook{h1}
}

func testPipelineStep() *task.JobTask {
	t := task.NewJobTask(pipeline.NewTask("test"))
	t.Spec.RunParams.Add(
		job.NewRunParam("ENV1", "VALUE1"),
	)
	t.Status.MarkedSuccess()

	return t
}

func init() {
	tools.DevelopmentSetup()
}

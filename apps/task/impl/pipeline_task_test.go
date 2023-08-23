package impl_test

import (
	"testing"

	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/test/conf"
	"github.com/infraboard/mflow/test/tools"
)

func TestQueryPipelineTask(t *testing.T) {
	req := task.NewQueryPipelineTaskRequest()
	set, err := impl.QueryPipelineTask(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(set))
}

func TestRunTestPipeline(t *testing.T) {
	req := pipeline.NewRunPipelineRequest(conf.C.CICD_PIPELINE_ID)
	req.RunBy = "test"
	ins, err := impl.RunPipeline(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToYaml(ins))
}

func TestRunmflowPipeline(t *testing.T) {
	req := pipeline.NewRunPipelineRequest(conf.C.MFLOW_PIPELINE_ID)
	req.RunBy = "test"
	req.RunParams = job.NewRunParamWithKVPaire(
		"GIT_SSH_URL", "git@github.com:infraboard/mflow.git",
		"GIT_BRANCH", "master",
		"GIT_COMMIT_ID", "57953a59e0ff5c93d0596696fbf6ffef6a90b446",
		build.SYSTEM_VARIABLE_APP_VERSION, "v0.0.10",
	)

	ins, err := impl.RunPipeline(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToYaml(ins))
}

func TestDescribePipelineTask(t *testing.T) {
	req := task.NewDescribePipelineTaskRequest(conf.C.PIPELINE_TASK_ID)
	ins, err := impl.DescribePipelineTask(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

func TestDeletePipelineTask(t *testing.T) {
	req := task.NewDeletePipelineTaskRequest(conf.C.PIPELINE_TASK_ID)
	ins, err := impl.DeletePipelineTask(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToYaml(ins))
}

package k8s_test

import (
	"testing"

	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/test/conf"
	"github.com/infraboard/mflow/test/tools"
)

func TestRun(t *testing.T) {
	jobSpec := tools.MustReadContentFile("test/job.yaml")
	params := job.NewRunParamSet()
	params.Add(
		&job.RunParam{
			Name:  "cluster_id",
			Value: "k8s-test",
		},
		&job.RunParam{
			Name:  "namespace",
			Value: "default",
		},
		&job.RunParam{
			Name:  "DB",
			Value: "xxx",
		},
		&job.RunParam{
			Name:  job.SYSTEM_VARIABLE_DEPLOY_ID,
			Value: conf.C.DEPLOY_ID,
		},
	)

	req := task.NewRunTaskRequest("test-job", jobSpec, params)
	ins, err := impl.Run(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToYaml(ins))
}

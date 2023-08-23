package k8s_test

import (
	"context"

	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/task/runner"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl runner.Runner
	ctx  = context.Background()
)

func init() {
	// 需要依赖cluster服务
	tools.DevelopmentSetup()
	impl = runner.GetRunner(job.RUNNER_TYPE_K8S_JOB)
}

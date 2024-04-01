package k8s

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/apps/task/runner"
	"github.com/infraboard/mpaas/clients/rpc"
	"github.com/rs/zerolog"
)

type K8sRunner struct {
	task task.PipelineService
	log  *zerolog.Logger

	mpaas *rpc.ClientSet
}

func (r *K8sRunner) Init() error {
	r.task = ioc.Controller().Get(task.AppName).(task.PipelineService)
	r.log = log.Sub("runner.k8s")
	r.mpaas = rpc.C()
	return nil
}

func (r *K8sRunner) RunnerType() job.RUNNER_TYPE {
	return job.RUNNER_TYPE_K8S_JOB
}

func init() {
	runner.Registry(&K8sRunner{})
}

package k8s

import (
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/apps/task/runner"
)

type K8sRunner struct {
	task task.PipelineService
	log  logger.Logger
}

func (r *K8sRunner) Init() error {
	r.task = ioc.GetController(task.AppName).(task.PipelineService)
	r.log = zap.L().Named("runner.k8s")
	return nil
}

func (r *K8sRunner) RunnerType() job.RUNNER_TYPE {
	return job.RUNNER_TYPE_K8S_JOB
}

func init() {
	runner.Registry(&K8sRunner{})
}

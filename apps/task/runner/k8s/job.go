package k8s

import (
	"context"
	"fmt"

	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/common/format"
	"github.com/infraboard/mpaas/provider/k8s/workload"
	batch "k8s.io/api/batch/v1"
	"sigs.k8s.io/yaml"
)

func (r *K8sRunner) Run(ctx context.Context, in *task.RunTaskRequest) (
	*task.JobTaskStatus, error) {

	// 获取k8s集群参数
	runnerParams := in.Params.K8SJobRunnerParams()

	// 生成k8s client
	k8sClient, err := runnerParams.Client(ctx)
	if err != nil {
		return nil, err
	}
	r.log.Debug().Msgf("new k8s client success")

	// 获取Job定义
	obj := new(batch.Job)
	jobYamlSpec := in.RenderJobSpec()
	r.log.Debug().Msgf("job rendered yaml spec: %s", jobYamlSpec)
	if err := yaml.Unmarshal([]byte(jobYamlSpec), obj); err != nil {
		return nil, err
	}

	// 修改任务名称
	obj.Name = in.Name
	obj.Namespace = runnerParams.Namespace

	// Job注入标签
	workload.InjectJobLabels(obj, in.Labels)
	// Job注入注解
	workload.InjectJobAnnotations(obj, in.Annotations())
	// 给Job容器注入环境变量
	workload.InjectPodEnvVars(&obj.Spec.Template.Spec, in.Params.EnvVars())

	status := task.NewJobTaskStatus()
	status.MarkedCreating()

	// 执行Job
	if !in.DryRun {
		r.log.Debug().Msgf("run job yaml: %s", format.MustToYaml(obj))
		obj, err := k8sClient.WorkLoad().CreateJob(ctx, obj)
		if err != nil {
			return nil, fmt.Errorf("create k8s job error, %s", err)
		}
		objYaml, err := yaml.Marshal(obj)
		if err != nil {
			return nil, fmt.Errorf("marshal k8s obj to yaml error, %s", err)
		}
		status.Detail = string(objYaml)
	}

	return status, nil
}

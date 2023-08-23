package k8s

import (
	"context"

	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/common/format"
	"github.com/infraboard/mflow/provider/k8s/workload"
	v1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/yaml"
)

func (r *K8sRunner) Run(ctx context.Context, in *task.RunTaskRequest) (
	*task.JobTaskStatus, error) {

	// 获取k8s集群参数
	runnerParams := in.Params.K8SJobRunnerParams()
	k8sClient, err := runnerParams.Client()
	if err != nil {
		return nil, err
	}

	// 获取Job定义
	obj := new(v1.Job)
	jobYamlSpec := in.RenderJobSpec()
	r.log.Debugf("job rendered yaml spec: %s", jobYamlSpec)
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
		r.log.Debugf("run job yaml: %s", format.MustToYaml(obj))
		obj, err = k8sClient.WorkLoad().CreateJob(ctx, obj)
		if err != nil {
			return nil, err
		}
	}

	objYaml, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	status.Detail = string(objYaml)
	return status, nil
}

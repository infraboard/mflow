package impl

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	v1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/yaml"

	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/apps/task/runner"
	"github.com/infraboard/mpaas/provider/k8s/config"
	"github.com/infraboard/mpaas/provider/k8s/meta"
	"github.com/infraboard/mpaas/provider/k8s/workload"
)

func (i *impl) RunJob(ctx context.Context, in *pipeline.Task) (
	*task.JobTask, error) {
	ins := task.NewJobTask(in)

	// 获取之前任务的状态, 因为里面有当前任务的审核状态
	err := i.GetJotTaskStatus(ctx, ins)
	if err != nil {
		return nil, err
	}

	// 开启审核后, 执行任务则 调整审核状态为等待中
	if ins.Spec.Audit.Enable {
		auditStatus := ins.AuditStatus()
		switch auditStatus {
		case pipeline.AUDIT_STAGE_PENDDING:
			ins.AuditStatusFlowTo(pipeline.AUDIT_STAGE_WAITING)
		case pipeline.AUDIT_STAGE_DENY:
			return nil, fmt.Errorf("任务审核未通过")
		}
		i.log.Debug().Msgf("任务: %s 审核状态为, %s", ins.Spec.TaskId, auditStatus)
	}

	// 如果不忽略执行, 并且审核通过, 则执行
	if in.Enabled() && ins.AuditPass() {
		// 任务状态检查与处理
		switch ins.Status.Stage {
		case task.STAGE_PENDDING:
			ins.Status.Stage = task.STAGE_CREATING
		case task.STAGE_ACTIVE:
			return nil, exception.NewConflict("任务: %s 当前处于运行中, 需要等待运行结束后才能执行", in.TaskId)
		}

		// 查询需要执行的Job
		req := job.NewDescribeJobRequestByName(in.JobName)

		j, err := i.job.DescribeJob(ctx, req)
		if err != nil {
			return nil, err
		}
		ins.Job = j
		ins.Spec.JobId = j.Meta.Id
		i.log.Info().Msgf("describe job success, %s[%s]", j.Spec.Name, j.Meta.Id)

		// 脱敏参数动态还原
		in.RunParams.RestoreSensitive(j.Spec.RunParams)

		// 合并允许参数(Job里面有默认值), 并检查参数合法性
		// 注意Param的合并是有顺序的，也就是参数优先级(低-->高):
		// 1. 系统变量(默认禁止修改)
		// 2. job默认变量
		// 3. job运行变量
		// 4. pipeline 运行变量
		// 5. pipeline 运行时变量
		params := job.NewRunParamSet()
		params.Add(ins.SystemRunParam()...)
		params.Add(j.Spec.RunParams.Params...)
		params.Merge(in.RunParams.Params...)
		err = i.LoadPipelineRunParam(ctx, in, params)
		if err != nil {
			return nil, err
		}

		// 校验参数合法性
		err = params.Validate()
		if err != nil {
			return nil, fmt.Errorf("校验任务【%s】参数错误, %s", j.Spec.DisplayName, err)
		}
		i.log.Info().Msgf("params check ok, %s", params)

		// 获取执行器执行
		r := runner.GetRunner(j.Spec.RunnerType)
		runReq := task.NewRunTaskRequest(ins.Spec.TaskId, j.Spec.RunnerSpec, params)
		runReq.DryRun = in.RunParams.DryRun
		runReq.Labels = in.Labels
		runReq.ManualUpdateStatus = j.Spec.ManualUpdateStatus

		i.log.Debug().Msgf("[%s] start run task: %s", ins.Spec.PipelineTask, in.TaskName)
		status, err := r.Run(ctx, runReq)
		if err != nil {
			return nil, fmt.Errorf("run job error, %s", err)
		}
		status.RunParams = params
		ins.Status = status

		// 添加搜索标签
		ins.BuildSearchLabel()
	}

	// 保存任务
	updateOpt := options.Update()
	updateOpt.SetUpsert(true)
	if _, err := i.jcol.UpdateByID(ctx, ins.Spec.TaskId, bson.M{"$set": ins}, updateOpt); err != nil {
		return nil, exception.NewInternalServerError("upsert a job task document error, %s", err)
	}
	return ins, nil
}

// 加载Pipeline 提供的运行时参数
func (i *impl) LoadPipelineRunParam(ctx context.Context, in *pipeline.Task, params *job.RunParamSet) error {
	if in.PipelineTask == "" {
		return nil
	}
	// 查询出Pipeline
	pt, err := i.DescribePipelineTask(ctx, task.NewDescribePipelineTaskRequest(in.PipelineTask))
	if err != nil {
		return err
	}

	// 获取出属于改Task的变量
	pipelineRunParams := pt.Params.GetTaskParams(in.StageNumberString(), in.TaskNumberString())
	// 合并PipelineTask传入的变量参数
	params.Merge(pipelineRunParams...)
	// 合并PipelineTask的运行时参数, Task运行时更新的
	params.Merge(pt.RuntimeRunParams()...)
	return nil
}

func (i *impl) GetJotTaskStatus(ctx context.Context, t *task.JobTask) error {
	if t.Spec.TaskId == "" {
		return nil
	}

	ins, err := i.DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(t.Spec.TaskId))
	if err != nil && exception.IsNotFoundError(err) {
		return err
	}
	if ins != nil && ins.Status != nil {
		t.Status = ins.Status
	}

	return nil

}

func (i *impl) JobTaskBatchSave(ctx context.Context, in *task.JobTaskSet) error {
	if _, err := i.jcol.InsertMany(ctx, in.ToDocs()); err != nil {
		return exception.NewInternalServerError("inserted job tasks document error, %s", err)
	}
	return nil
}

func (i *impl) QueryJobTask(ctx context.Context, in *task.QueryJobTaskRequest) (
	*task.JobTaskSet, error) {
	r := newQueryRequest(in)

	filter := r.FindFilter()
	i.log.Debug().Msgf("query filter: %v", filter)
	resp, err := i.jcol.Find(ctx, filter, r.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find deploy error, error is %s", err)
	}

	set := task.NewJobTaskSet()
	// 循环
	for resp.Next(ctx) {
		ins := task.NewDefaultJobTask()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode deploy error, error is %s", err)
		}
		ins.Init()
		set.Add(ins)
	}

	// count
	count, err := i.jcol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, exception.NewInternalServerError("get deploy count error, error is %s", err)
	}
	set.Total = count
	return set, nil
}

func (i *impl) CheckAllowUpdate(ctx context.Context, ins *task.JobTask, token string, force bool) error {
	// 校验更新合法性
	err := ins.ValidateToken(token)
	if err != nil {
		return err
	}

	// 修改任务状态
	if !force && ins.Status.IsComplete() {
		return exception.NewBadRequest("状态[%s], 已经结束的任务不能更新状态", ins.Status.Stage)
	}
	return nil
}

// 更新任务运行结果
func (i *impl) UpdateJobTaskOutput(ctx context.Context, in *task.UpdateJobTaskOutputRequest) (
	*task.JobTask, error) {
	ins, err := i.DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(in.Id))
	if err != nil {
		return nil, err
	}

	// 校验更新合法性
	err = i.CheckAllowUpdate(ctx, ins, in.UpdateToken, in.Force)
	if err != nil {
		return nil, err
	}
	ins.Status.UpdateOutput(in)

	// 只更新任务状态
	if _, err := i.jcol.UpdateByID(ctx, ins.Spec.TaskId, bson.M{"$set": bson.M{"status": ins.Status}}); err != nil {
		return nil, exception.NewInternalServerError("update task(%s) document error, %s",
			in.Id, err)
	}

	return ins, nil
}

// 更新Job状态
func (i *impl) UpdateJobTaskStatus(ctx context.Context, in *task.UpdateJobTaskStatusRequest) (
	*task.JobTask, error) {
	ins, err := i.DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(in.Id))
	if err != nil {
		return nil, err
	}

	// 校验更新合法性
	err = i.CheckAllowUpdate(ctx, ins, in.UpdateToken, in.ForceUpdateStatus)
	if err != nil {
		return nil, err
	}

	// 状态更新
	preStatus := ins.Status.Stage
	ins.Status.UpdateStatus(in)

	// Job Task状态变更回调
	i.JobTaskStatusChangedCallback(ctx, ins)

	// 更新数据库
	if err := i.updateJobTaskStatus(ctx, ins); err != nil {
		return nil, err
	}

	// Pipeline Task 状态变更回调
	if ins.Spec.PipelineTask != "" {
		// 如果状态未变化, 不触发流水线更新
		if !in.ForceTriggerPipeline && preStatus.Equal(in.Stage) {
			i.log.Debug().Msgf("task %s status not changed: %s, skip update pipeline", in.Id, in.Stage)
			return ins, nil
		}
		_, err := i.PipelineTaskStatusChanged(ctx, ins)
		if err != nil {
			return nil, err
		}
	}
	return ins, nil
}

func (i *impl) JobTaskStatusChangedCallback(ctx context.Context, in *task.JobTask) {
	if in.Status == nil {
		return
	}

	i.log.Debug().Msgf("task %s 执行状态变化回调...", in.Spec.TaskId)

	// 个人通知
	for index := range in.Spec.MentionUsers {
		mu := in.Spec.MentionUsers[index]
		i.TaskMention(ctx, mu, in)
	}
	if len(in.Spec.MentionUsers) > 0 {
		i.updateJobTaskMentionUser(ctx, in.Spec.TaskId, in.Spec.MentionUsers)
	}

	// 群组通知
	imRobotHooks := in.Spec.MatchedImRobotNotify(in.Status.Stage.String())
	i.log.Debug().Msgf("task %s 群组通知: %v", in.Spec.TaskId, imRobotHooks)
	i.hook.SendTaskStatus(ctx, imRobotHooks, in)
	if len(imRobotHooks) > 0 {
		i.updateJobTaskImRobotNotify(ctx, in.Spec.TaskId, imRobotHooks)
	}

	// WebHook回调
	webhooks := in.Spec.MatchedWebHooks(in.Status.Stage.String())
	i.log.Debug().Msgf("task %s WebHook通知: %v", in.Spec.TaskId, webhooks)
	i.hook.SendTaskStatus(ctx, webhooks, in)
	if len(webhooks) > 0 {
		i.updateJobTaskWebHook(ctx, in.Spec.TaskId, webhooks)
	}
}

// 任务执行详情
func (i *impl) DescribeJobTask(ctx context.Context, in *task.DescribeJobTaskRequest) (
	*task.JobTask, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	ins := task.NewDefaultJobTask()
	if err := i.jcol.FindOne(ctx, bson.M{"_id": in.Id}).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("job task %s not found", in.Id)
		}

		return nil, exception.NewInternalServerError("find job task %s error, %s", in.Id, err)
	}

	if in.WithJob {
		j, err := i.job.DescribeJob(ctx, job.NewDescribeJobRequestById(ins.Spec.JobId))
		if err != nil {
			return nil, err
		}
		ins.Job = j
	}

	ins.Init()
	return ins, nil
}

// 删除任务
func (i *impl) DeleteJobTask(ctx context.Context, in *task.DeleteJobTaskRequest) (
	*task.JobTask, error) {
	descReq := task.NewDescribeJobTaskRequest(in.Id)
	descReq.WithJob = true
	ins, err := i.DescribeJobTask(ctx, descReq)
	if err != nil {
		return nil, err
	}

	// 清理Job关联的临时资源
	err = i.CleanTaskResource(ctx, ins)
	if err != nil {
		if !in.Force {
			return nil, err
		}
		i.log.Warn().Msgf("force delete, but has error, %s", err)
	}

	// 删除本地记录
	_, err = i.jcol.DeleteOne(ctx, bson.M{"_id": in.Id})
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (i *impl) CleanTaskResource(ctx context.Context, in *task.JobTask) error {
	if !in.HasJobSpec() {
		return nil
	}

	switch in.Job.Spec.RunnerType {
	case job.RUNNER_TYPE_K8S_JOB:
		k8sParams := in.Status.RunParams.K8SJobRunnerParams()
		k8sClient, err := k8sParams.Client(ctx)
		if err != nil {
			return err
		}

		// 清除临时挂载的configmap
		for i := range in.Status.TemporaryResources {
			resource := in.Status.TemporaryResources[i]
			if resource.IsReleased() {
				continue
			}
			switch resource.Kind {
			case config.CONFIG_KIND_CONFIG_MAP.String():
				cmDeleteReq := meta.NewDeleteRequest(resource.Name).WithNamespace(k8sParams.Namespace)
				err = k8sClient.Config().DeleteConfigMap(ctx, cmDeleteReq)
				if err != nil {
					return fmt.Errorf("delete config map error, %s", err)
				}
				in.Status.AddEvent(pipeline.EVENT_LEVEL_DEBUG, "delete job runtime env configmap: %s", resource.Name)
				resource.ReleaseAt = time.Now().Unix()
			}
		}

		// 清理Job
		detail := in.GetStatusDetail()
		if detail == "" {
			return fmt.Errorf("no k8s job found in status detail")
		}

		obj := new(v1.Job)
		if err := yaml.Unmarshal([]byte(detail), obj); err != nil {
			return err
		}

		req := meta.NewDeleteRequest(obj.Name)
		req.Namespace = obj.Namespace
		err = k8sClient.WorkLoad().DeleteJob(ctx, req)
		if err != nil {
			i.log.Error().Msgf("delete k8s job error, %s", err)
		}
	}

	return nil
}

// 查询Task日志
func (i *impl) WatchJobTaskLog(in *task.WatchJobTaskLogRequest, stream task.JobRPC_WatchJobTaskLogServer) error {
	writer := NewWatchJobTaskLogServerWriter(stream)
	writer.WriteMessagef("正在查询Task[%s]的日志 请稍等...", in.TaskId)

	// 等待Task的Pod正常启动
	t, err := i.WaitPodLogReady(stream.Context(), in, writer)
	if err != nil {
		return err
	}

	switch t.Job.Spec.RunnerType {
	case job.RUNNER_TYPE_K8S_JOB:
		k8sParams := t.Spec.RunParams.K8SJobRunnerParams()
		k8sClient, err := k8sParams.Client(stream.Context())
		if err != nil {
			return fmt.Errorf("create k8s client error, %s", err)
		}

		// 找到Job执行的Pod
		podReq := meta.NewListRequest().
			WithNamespace(k8sParams.Namespace).
			WithLabelSelector(meta.NewLabelSelector().Add("job-name", t.Spec.TaskId))
		pods, err := k8sClient.WorkLoad().ListPod(stream.Context(), podReq)
		if err != nil {
			return fmt.Errorf("list job pod error, %s", err)
		}
		if len(pods.Items) == 0 {
			return fmt.Errorf("job's pod not found by lable job-name=%s", t.Spec.TaskId)
		}

		req := workload.NewWatchConainterLogRequest()
		req.PodName = pods.Items[0].Name
		req.Namespace = k8sParams.Namespace
		req.Container = in.ContainerName
		r, err := k8sClient.WorkLoad().WatchConainterLog(stream.Context(), req)
		if err != nil {
			return fmt.Errorf("watch container log error, %s", err)
		}
		defer r.Close()

		// copy日志流
		_, err = io.Copy(writer, r)
		return err
	}

	return nil
}

func (i *impl) WaitPodLogReady(
	ctx context.Context,
	in *task.WatchJobTaskLogRequest,
	writer *WatchJobTaskLogServerWriter,
) (*task.JobTask, error) {
	maxRetryCount := 0
WAIT_TASK_ACTIVE:
	// 查询Task信息
	descReq := task.NewDescribeJobTaskRequest(in.TaskId).SetWithJob(true)
	t, err := i.DescribeJobTask(ctx, descReq)
	if err != nil {
		return nil, err
	}

	pod, err := t.Status.GetLatestPod()
	if err != nil {
		return nil, err
	}

	if t.Status.Stage.Equal(task.STAGE_FAILED) && pod == nil {
		writer.WriteMessagef("[%s]运行失败, %s", t.Spec.TaskName, t.Status.Message)
		return t, fmt.Errorf("任务没有成功创建Pod")
	}

	writer.WriteMessagef("任务当前状态: [%s]", t.Status.Stage)
	if !t.Status.IsComplete() && maxRetryCount < 30 {
		if pod != nil {
			writer.WriteMessagef("当前Pod状态: [%s], 等待任务启动中...", pod.Status.Phase)
			// Job状态运行成功，返回
			if pod.Status.Phase != "Pending" {
				return t, nil
			}
		}
		time.Sleep(2 * time.Second)
		maxRetryCount++
		goto WAIT_TASK_ACTIVE
	}

	writer.WriteMessagef("当前Pod状态: [%s]", pod.Status.Phase)

	// 容器Init Container启动失败，则默认查看失败的日志
	if in.ContainerName == "" {
		for _, initContainer := range pod.Status.InitContainerStatuses {
			if !initContainer.Ready {
				in.ContainerName = initContainer.Name
				writer.WriteMessagef("Init Container: %s 启动失败, 查看失败日志...", initContainer.Name)
				break
			}
		}
	}
	return t, nil
}

// Task Debug
func (i *impl) DebugJobTask(ctx context.Context, in *task.DebugJobTaskRequest) {
	term := in.WebTerminal()

	// 查询Task信息
	t, err := i.DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(in.TaskId))
	if err != nil {
		term.Failed(err)
		return
	}

	switch t.Job.Spec.RunnerType {
	case job.RUNNER_TYPE_K8S_JOB:
		k8sParams := t.Status.RunParams.K8SJobRunnerParams()
		k8sClient, err := k8sParams.Client(ctx)
		if err != nil {
			term.Failed(fmt.Errorf("初始化k8s客户端失败, %s", err))
			return
		}

		// 找到Job执行的Pod
		term.WriteTextln("正在查询Job Task【%s】运行的Pod", t.Spec.TaskId)
		podReq := meta.NewListRequest().
			WithNamespace(k8sParams.Namespace).
			WithLabelSelector(meta.NewLabelSelector().Add("job-name", t.Spec.TaskId))
		pods, err := k8sClient.WorkLoad().ListPod(ctx, podReq)
		if err != nil {
			term.Failed(err)
			return
		}
		if len(pods.Items) == 0 {
			term.Failed(fmt.Errorf("job's pod not found by lable job-name=%s", t.Spec.TaskId))
			return
		}

		targetCopyPod := pods.Items[0]
		term.WriteTextln("Job Task【%s】位于Namespace: %s, PodName: %s",
			t.Spec.TaskId, targetCopyPod.Namespace,
			targetCopyPod.Name,
		)

		req := in.CopyPodRunRequest(k8sParams.Namespace, targetCopyPod.Name)
		req.SetAttachTerminal(term)
		req.Remove = true

		_, err = k8sClient.WorkLoad().CopyPodRun(ctx, req)
		if err != nil {
			term.Failed(err)
			return
		}
	default:
		term.Failed(fmt.Errorf("unknonw runner type %s", t.Job.Spec.RunnerType))
		return
	}
}

func NewWatchJobTaskLogServerWriter(
	stream task.JobRPC_WatchJobTaskLogServer) *WatchJobTaskLogServerWriter {
	return &WatchJobTaskLogServerWriter{
		stream: stream,
		buf:    task.NewJobTaskStreamReponse(),
	}
}

type WatchJobTaskLogServerWriter struct {
	stream task.JobRPC_WatchJobTaskLogServer
	buf    *task.JobTaskStreamReponse
}

func (w *WatchJobTaskLogServerWriter) Write(p []byte) (n int, err error) {
	w.buf.Data = p
	err = w.stream.Send(w.buf)
	if err != nil {
		return 0, err
	}
	w.buf.ReSet()
	return len(p), nil
}

func (w *WatchJobTaskLogServerWriter) WriteMessagef(format string, a ...any) {
	_, err := w.Write([]byte(fmt.Sprintf(format+"\n", a...)))
	if err != nil {
		log.L().Error().Msgf("write message error, %s", err)
	}
}

// 审核任务
func (i *impl) AuditJobTask(
	ctx context.Context,
	in *task.AuditJobTaskRequest) (
	*task.JobTask, error) {
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest("参数校验失败: %s", err.Error())
	}

	t, err := i.DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(in.TaskId))
	if err != nil {
		return nil, err
	}
	currentAuditStaus := t.AuditStatus()

	if !t.Spec.Audit.Enable {
		return nil, exception.NewBadRequest("任务没有开启审核")
	}

	if !t.Spec.IsAuditor(in.Status.AuditBy) {
		return nil, exception.NewBadRequest("你不在审核人名单中")
	}

	if err := t.CheckUpdateStage(in.Status.Stage); err != nil {
		return nil, exception.NewBadRequest("状态更新异常, %s", err)
	}

	// 更新任务审核状态
	t.UpdateAuditStatus(in.Status)

	// 审核失败结束任务, 更新任务状态, 并且触发流水线状态
	if !t.AuditPass() {
		// 如果审核通过, 再次触发任务状态变化
		t.Status.MarkedError(fmt.Errorf("审核没通过, %s", in.Status.Comment))
		_, err = i.PipelineTaskStatusChanged(ctx, t)
		if err != nil {
			i.log.Error().Msgf("流水线状态更新异常: %s", err)
		}
	}

	if err := i.updateJobTaskStatus(ctx, t); err != nil {
		return nil, err
	}

	// 审核通过, 且审核状态为等待执行
	if t.AuditPass() && currentAuditStaus.Equal(pipeline.AUDIT_STAGE_WAITING) {
		_, err := i.RunJob(ctx, t.Spec)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

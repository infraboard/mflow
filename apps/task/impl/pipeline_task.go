package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/lock"
	"github.com/infraboard/mflow/apps/approval"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/apps/trigger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 执行Pipeline
func (i *impl) RunPipeline(ctx context.Context, in *pipeline.RunPipelineRequest) (
	*task.PipelineTask, error) {
	// 检查Pipeline请求参数
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	// 查询需要执行的Pipeline
	pReq := pipeline.NewDescribePipelineRequest(in.PipelineId)
	pReq.WithJob = true
	p, err := i.pipeline.DescribePipeline(ctx, pReq)
	if err != nil {
		return nil, err
	}

	// 如果Pipeline是窜行执行, 需要加锁
	if p.IsSequence() {
		if in.SequenceLabelKey == "" {
			return nil, exception.NewBadRequest("窜行执行时, sequence_label_key必须传递")
		}
		if in.SequenceLabelValue() == "" {
			return nil, exception.NewBadRequest("窜行执行时, label not found sequence_label_key value")
		}
		// 避免构建同时触发, 更新时加锁
		m := lock.L().New(in.SequenceLabelValue(), time.Duration(i.LockTimeoutSeconds)*time.Second)
		if err := m.Lock(ctx); err != nil {
			return nil, fmt.Errorf("lock error, %s", err)
		}
		defer m.UnLock(ctx)
	}

	// 检查Pipeline状态
	if err := i.CheckPipelineAllowRun(ctx, p, in.ApprovalId, in.SequenceLabel()); err != nil {
		return nil, err
	}

	// 从pipeline 取出需要执行的任务
	ins := task.NewPipelineTask(p, in)

	ts, err := ins.NextRun()
	if err != nil {
		return nil, fmt.Errorf("find next run task error, %s", err)
	}
	if ts == nil || ts.Len() == 0 {
		return nil, fmt.Errorf("not job task to run")
	}
	if in.DryRun {
		return ins, nil
	}

	// 保存Job Task, 所有JobTask 批量生成, 全部处于Pendding状态, 然后入库 等待状态更新
	err = i.JobTaskBatchSave(ctx, ins.JobTasks())
	if err != nil {
		return nil, err
	}

	ins.MarkedRunning()
	// 保存Pipeline状态
	if _, err := i.pcol.InsertOne(ctx, ins); err != nil {
		return nil, exception.NewInternalServerError("inserted a pipeline task document error, %s", err)
	}

	/*
		运行 第一批Task, 驱动Pipeline执行
	*/
	i.log.Debug().Msgf("run pipeline tasks: %v", ts.TaskNames())
	for _, t := range ts.Items {
		uReq := task.NewUpdateJobTaskStatusRequest(t.Spec.TaskId)
		uReq.UpdateToken = t.Spec.UpdateToken
		reqp, err := i.RunJob(ctx, t.Spec)
		if err != nil {
			uReq.MarkError(err)
		} else {
			uReq.Stage = reqp.Status.Stage
		}
		_, err = i.UpdateJobTaskStatus(ctx, uReq)
		if err != nil {
			i.log.Error().Msgf("update pipeline status form task error, %s", err)
		}
	}
	return ins, nil
}

// lables 参数: 因为一个Pipeline可能同时被多个服务的多个构建使用到, 只是不允许同一个构建的任务被同时触发，
// 但是允许不同构建的任务同时触发
func (i *impl) CheckPipelineAllowRun(
	ctx context.Context,
	ins *pipeline.Pipeline,
	approlvalId string,
	lables map[string]string,
) error {
	// 1. 检查审核状态
	if ins.Spec.RequiredApproval {
		if approlvalId == "" {
			return exception.NewBadRequest("流水线需要审核单才能运行")
		}
		a, err := i.approval.DescribeApproval(
			ctx, approval.NewDescribeApprovalRequest(approlvalId),
		)
		if err != nil {
			return err
		}

		if a.Pipeline.Meta.Id != ins.Meta.Id {
			return fmt.Errorf("审核单不属于当前执行的流水线")
		}

		if !a.Status.IsAllowPublish() {
			return fmt.Errorf("当前状态: %s 不允许发布", a.Status.Stage)
		}
	}

	// 2. 检查当前pipeline是否已经处于运行中
	if ins.IsSequence() {
		// 查询当前pipeline task有没有未完成的任务
		req := task.NewQueryPipelineTaskRequest()
		req.Stages = task.UnCompleteStage
		req.PipelineId = ins.Meta.Id
		req.Page.PageSize = 1
		if lables != nil {
			req.Labels = lables
		}
		set, err := i.QueryPipelineTask(ctx, req)
		if err != nil {
			return err
		}
		// 没有最近的任务
		if set.Len() > 0 {
			return task.ERR_PIPELINE_IS_RUNNING.WithMessagef("pipeline task %s is UnCompleted", set.Items[0].Meta.Id)
		}
	}

	return nil
}

// Pipeline中任务有变化时,
// 如果执行成功则 继续执行, 如果失败则标记Pipeline结束
// 当所有任务成功结束时标记Pipeline执行成功
func (i *impl) PipelineTaskStatusChanged(ctx context.Context, in *task.JobTask) (
	*task.PipelineTask, error) {
	if in == nil || in.Status == nil {
		return nil, exception.NewBadRequest("job task or job task status is nil")
	}

	if in.Spec.PipelineTask == "" {
		return nil, exception.NewBadRequest("Pipeline Id参数缺失")
	}

	runErrorJobTasks := []*task.UpdateJobTaskStatusRequest{}
	// 获取Pipeline Task, 因为Job Task是先保存在触发的回调, 这里获取的Pipeline Task是最新的
	descReq := task.NewDescribePipelineTaskRequest(in.Spec.PipelineTask)
	p, err := i.DescribePipelineTask(ctx, descReq)
	if err != nil {
		return nil, err
	}

	defer func() {
		// 更新当前任务的pipeline task状态
		i.mustUpdatePipelineStatus(ctx, p)

		// 如果JobTask正常执行, 则等待回调更新, 如果执行失败 则需要立即更新JobTask状态
		for index := range runErrorJobTasks {
			_, err = i.UpdateJobTaskStatus(ctx, runErrorJobTasks[index])
			if err != nil {
				p.MarkedFailed(err)
				i.mustUpdatePipelineStatus(ctx, p)
			}
		}
	}()

	// 更新Pipeline Task 运行时环境变量
	p.Status.RuntimeEnvs.Merge(in.RuntimeRunParams()...)

	switch in.Status.Stage {
	case task.STAGE_PENDDING,
		task.STAGE_SCHEDULING,
		task.STAGE_CREATING,
		task.STAGE_ACTIVE,
		task.STAGE_CANCELING:
		// Task状态无变化
		return p, nil
	case task.STAGE_CANCELED:
		// 任务取消, pipeline 取消执行
		p.MarkedCanceled()
		return p, nil
	case task.STAGE_FAILED:
		// 任务执行失败, 更新Pipeline状态为失败
		if !in.Spec.RunParams.IgnoreFailed {
			p.MarkedFailed(in.Status.MessageToError())
			return p, nil
		}
	case task.STAGE_SUCCEEDED:
		// 任务运行成功, pipeline继续执行
		i.log.Info().Msgf("task: %s run successed", in.Spec.TaskId)
	}

	/*
		task执行成功或者忽略执行失败, 此时pipeline 仍然处于运行中, 需要获取下一个任务执行
	*/
	nexts, err := p.NextRun()
	if err != nil {
		p.MarkedFailed(err)
		return p, nil
	}

	// 如果没有需要执行的任务, Pipeline执行结束, 更新Pipeline状态为成功
	if nexts == nil || nexts.Len() == 0 {
		p.MarkedSuccess()
		return p, nil
	}

	// 如果有需要执行的JobTask, 继续执行
	for index := range nexts.Items {
		item := nexts.Items[index]
		// 如果任务执行成功则等待任务的回调更新任务状态
		// 如果任务执行失败, 直接更新任务状态
		t, err := i.RunJob(ctx, item.Spec)
		if err != nil {
			updateT := task.NewUpdateJobTaskStatusRequest(item.Spec.TaskId)
			updateT.UpdateToken = item.Spec.UpdateToken
			updateT.MarkError(err)
			runErrorJobTasks = append(runErrorJobTasks, updateT)
		} else {
			item.Status = t.Status
			item.Job = t.Job
		}
	}

	return p, nil
}

// 更新Pipeline状态
func (i *impl) updatePipelineStatus(ctx context.Context, in *task.PipelineTask) (*task.PipelineTask, error) {
	// 立即更新Pipeline Task状态
	if err := i.updatePiplineTaskStatus(ctx, in); err != nil {
		return nil, err
	}

	// 执行PipelineTask状态变更回调
	i.PipelineStatusChangedCallback(ctx, in)

	// 补充回调执行状态
	if err := i.updatePiplineTaskStatus(ctx, in); err != nil {
		return nil, err
	}
	return in, nil
}

func (i *impl) PipelineStatusChangedCallback(ctx context.Context, in *task.PipelineTask) {
	if !in.HasJobSpec() {
		return
	}

	if in.Status == nil {
		return
	}

	// WebHook回调
	webhooks := in.Pipeline.Spec.MatchedWebHooks(in.Status.Stage.String())
	i.hook.SendTaskStatus(ctx, webhooks, in)

	// 关注人通知回调
	for index := range in.Pipeline.Spec.MentionUsers {
		mu := in.Pipeline.Spec.MentionUsers[index]
		i.TaskMention(ctx, mu, in)
	}

	// 是否需要运行下一个Pipeline
	if in.Pipeline.Spec.NextPipeline != "" {
		// 如果有审核单则提交审核单, 没有则直接运行
		i.log.Debug().Msgf("next pipeline: %s", in.Pipeline.Spec.NextPipeline)
	}

	// 事件队列回调通知, 通过构建任务PipelineTask执行完成
	buildEventId, buildConfId := in.GetEventId(), in.GetBuildConfId()
	if in.IsComplete() && in.HasBuildEvent() {
		i.log.Debug().Msgf("build conf %s next event", buildConfId)
		tReq := trigger.NewEventQueueTaskCompleteRequest(buildEventId, buildConfId)
		bs, err := i.trigger.EventQueueTaskComplete(ctx, tReq)
		if err != nil {
			in.AddErrorEvent("触发队列回调失败: %s", err)
		}
		if bs != nil {
			in.AddSuccessEvent("触发队列成功: %s", bs.BuildConfig.Spec.Name).SetDetail(bs.String())
		}
	}
}

// 更新Pipeline状态
func (i *impl) mustUpdatePipelineStatus(ctx context.Context, in *task.PipelineTask) {
	_, err := i.updatePipelineStatus(ctx, in)
	if err != nil {
		i.log.Err(err)
	}
}

// 查询Pipeline任务
func (i *impl) QueryPipelineTask(ctx context.Context, in *task.QueryPipelineTaskRequest) (
	*task.PipelineTaskSet, error) {
	r := newQueryPipelineTaskRequest(in)
	filter := r.FindFilter()
	i.log.Debug().Msgf("query pipeline task filter: %s", filter)
	resp, err := i.pcol.Find(ctx, filter, r.FindOptions())
	if err != nil {
		return nil, exception.NewInternalServerError("find pipeline task error, error is %s", err)
	}

	set := task.NewPipelineTaskSet()
	// 循环
	for resp.Next(ctx) {
		ins := task.NewDefaultPipelineTask()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode pipeline task  error, error is %s", err)
		}
		set.Add(ins)
	}

	// count
	count, err := i.pcol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, exception.NewInternalServerError("get pipeline task count error, error is %s", err)
	}
	set.Total = count
	return set, nil
}

// 查询Pipeline任务详情
func (i *impl) DescribePipelineTask(ctx context.Context, in *task.DescribePipelineTaskRequest) (
	*task.PipelineTask, error) {
	filter := bson.M{"_id": in.Id}

	ins := task.NewDefaultPipelineTask()
	if err := i.pcol.FindOne(ctx, filter).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("pipeline task %s not found", in.Id)
		}

		return nil, exception.NewInternalServerError("find pipeline task %s error, %s", in.Id, err)
	}

	// 补充该PipelineTask管理的JobTask
	query := task.NewQueryTaskRequest()
	query.PipelineTaskId = in.Id
	query.SortType = task.SORT_TYPE_ASCEND
	tasks, err := i.QueryJobTask(ctx, query)
	if err != nil {
		return nil, err
	}

	// 将tasks 填充给pipeline task
	for i := range tasks.Items {
		t := tasks.Items[i]
		stage := ins.GetStage(t.Spec.StageName)
		if stage != nil {
			stage.Add(t)
		}
	}

	// 执行参数 回填回Pipeline, 方便前端直接展示
	ins.SetParamValueToPipeline()
	return ins, nil
}

// 删除Pipeline任务详情
func (i *impl) DeletePipelineTask(ctx context.Context, in *task.DeletePipelineTaskRequest) (
	*task.PipelineTask, error) {
	ins, err := i.DescribePipelineTask(ctx, task.NewDescribePipelineTaskRequest(in.Id))
	if err != nil {
		return nil, err
	}
	// 运行中的流水线不运行删除, 先取消 才能删除
	if ins.IsRunning() {
		return nil, fmt.Errorf("流水线运行结束才能删除, 如果没结束, 请先取消再删除")
	}

	// 删除关联的构建记录
	rDel := trigger.NewDeleteRecordByPipelineTask(ins.Meta.Id)
	if err := i.trigger.DeleteRecord(ctx, rDel); err != nil {
		i.log.Error().Msgf("delete trigger record error %s", err)
	}

	// 删除该Pipeline下所有的Job Task
	tasks := ins.Status.JobTasks()
	for index := range tasks.Items {
		t := tasks.Items[index]

		// 没有运行过的任务不需要清理
		if t.Status.Stage.Equal(task.STAGE_PENDDING) {
			continue
		}

		deleteReq := task.NewDeleteJobTaskRequest(t.Spec.TaskId)
		deleteReq.Force = in.Force
		_, err := i.DeleteJobTask(ctx, deleteReq)
		if err != nil {
			if !exception.IsNotFoundError(err) {
				return nil, err
			}
		}
	}

	if err := i.deletePipelineTask(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}

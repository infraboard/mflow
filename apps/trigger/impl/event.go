package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcenter/apps/service"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/lock"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/apps/trigger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 应用事件处理
func (i *impl) HandleEvent(ctx context.Context, in *trigger.Event) (
	*trigger.Record, error) {
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	ins := trigger.NewRecord(in)

	// 查询服务相关信息
	i.log.Debug().Msgf("查询事件关联服务信息...")
	svc, err := i.mcenter.Service().DescribeService(
		ctx,
		service.NewDescribeServiceRequest(in.Token),
	)
	if err != nil {
		return nil, exception.NewBadRequest("查询服务%s异常, %s", in.Token, err)
	}
	ins.Event.ServiceInfo = svc.Desense().ToJSON()

	// 获取该服务对应事件的触发配置
	i.log.Debug().Msgf("查询该服务是否有关联构建配置...")
	req := build.NewQueryBuildConfigRequest()
	req.AddService(in.Token)
	req.Event = in.Name
	req.SetEnabled(true)
	set, err := i.build.QueryBuildConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	if set.Len() == 0 {
		return nil, fmt.Errorf("没有找到服务[%s]的构建规则", svc.Spec.Name)
	}

	// 子事件匹配
	i.log.Debug().Msgf("服务[%s]关联%d个构建配置", svc.Spec.Name, set.Len())
	matched := set.MatchSubEvent(in.SubName)
	for index := range matched.Items {
		// 执行构建配置匹配的流水线
		buildConf := matched.Items[index]
		bs := i.RunBuildConf(ctx, in, buildConf)
		ins.AddBuildStatus(bs)
	}

	// 保存
	if _, err := i.col.InsertOne(ctx, ins); err != nil {
		return nil, exception.NewInternalServerError("inserted a deploy document error, %s", err)
	}
	return ins, nil
}

func (i *impl) RunBuildConf(ctx context.Context, in *trigger.Event, buildConf *build.BuildConfig) *trigger.BuildStatus {
	bs := trigger.NewBuildStatus(buildConf)

	pipelineId := buildConf.Spec.PipelineId
	if pipelineId == "" {
		bs.ErrorMessage = "未配置流水线"
		return bs
	}

	runReq := pipeline.NewRunPipelineRequest(pipelineId)
	runReq.RunBy = "@" + in.UUID()
	runReq.TriggerMode = pipeline.TRIGGER_MODE_EVENT
	runReq.DryRun = in.SkipRunPipeline
	runReq.Labels[build.PIPELINE_TASK_EVENT_ID_LABLE_KEY] = in.Id
	runReq.Labels[build.PIPELINE_TASK_BUILD_CONFIG_ID_LABLE_KEY] = buildConf.Meta.Id
	runReq.Domain = buildConf.Spec.Scope.Domain
	runReq.Namespace = buildConf.Spec.Scope.Namespace
	// 如果窜行时, 不允许BuildConf并发执行
	runReq.SequenceKey = buildConf.Meta.Id

	// 补充Build用户自定义变量
	runReq.AddRunParam(buildConf.Spec.CustomParams...)

	// 补充Gitlab事件特有的变量
	switch in.Provider {
	case trigger.EVENT_PROVIDER_GITLAB:
		event, err := in.GetGitlabEvent()
		if err != nil {
			bs.Failed(err)
			return bs
		}

		// 补充Git信息
		runReq.AddRunParam(in.GitRunParams().Params...)
		// 补充版本信息
		switch buildConf.Spec.VersionNamedRule {
		case build.VERSION_NAMED_RULE_DATE_BRANCH_COMMIT:
			runReq.AddRunParam(event.DateCommitVersion(buildConf.Spec.VersionPrefix))
		case build.VERSION_NAMED_RULE_GIT_TAG:
			runReq.AddRunParam(event.TagVersion(buildConf.Spec.VersionPrefix))
		}
	}

	i.log.Debug().Msgf("run pipeline req: %s, params: %v", runReq.PipelineId, runReq.RunParamsKVMap())
	pt, err := i.task.RunPipeline(ctx, runReq)
	if err != nil {
		if exception.IsError(err, task.ERR_PIPELINE_IS_RUNNING) {
			bs.Enqueue()
		} else {
			i.log.Debug().Msgf("run pipeline error, %s", err)
			bs.Failed(err)
		}
	} else {
		i.log.Debug().Msgf("update run build conf pipeline task id: %s", pt.Meta.Id)
		bs.Running(pt)
	}
	return bs
}

// 查询事件
func (i *impl) QueryRecord(ctx context.Context, in *trigger.QueryRecordRequest) (
	*trigger.RecordSet, error) {
	r := newQueryRequest(in)
	filter := r.FindFilter()
	i.log.Debug().Msgf("query filter: %s", filter)
	resp, err := i.col.Find(ctx, filter, r.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find event record error, error is %s", err)
	}

	set := trigger.NewRecordSet()
	// 循环
	for resp.Next(ctx) {
		ins := trigger.NewDefaultRecord()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode event record error, error is %s", err)
		}
		set.Add(ins)
	}

	if in.WithPipelineTask && len(set.Items) > 0 {
		set.PiplineTaskIds()
		tReq := task.NewQueryPipelineTaskRequest()
		ts, err := i.task.QueryPipelineTask(ctx, tReq)
		if err != nil {
			return nil, err
		}
		set.UpdatePiplineTask(ts.Items...)
	}

	// count
	count, err := i.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, exception.NewInternalServerError("get event record count error, error is %s", err)
	}
	set.Total = count
	return set, nil
}

// 事件队列任务执行完成通知
func (i *impl) EventQueueTaskComplete(ctx context.Context, in *trigger.EventQueueTaskCompleteRequest) (
	*trigger.BuildStatus, error) {
	// 请求校验
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate param error, %s", err)
	}

	// 查询该Pipeline Task关联的构建记录
	req := trigger.NewQueryRecordRequest()
	req.AddEventId(in.EventId)
	rs, err := i.QueryRecord(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(rs.Items) == 0 {
		return nil, exception.NewBadRequest("event %s record not found", in.EventId)
	}

	// 查询构建记录对应的BuildConf
	currentRecord := rs.Items[0]
	currentBuildConf, currentBuildConfIndex := currentRecord.GetBuildStatusByBuildConfId(in.BuildConfId)
	if currentBuildConf == nil || currentBuildConf.BuildConfig == nil {
		return nil, exception.NewBadRequest("find current record build config %s not found", in.BuildConfId)
	}

	// 查询PipelineTask
	pt, err := i.task.DescribePipelineTask(ctx, task.NewDescribePipelineTaskRequest(currentBuildConf.PiplineTaskId))
	if err != nil {
		return nil, err
	}

	if !in.ForceUpdate && currentBuildConf.IsComplete() {
		return nil, exception.NewBadRequest("task is complete")
	}

	// 更新构建记录完成状态
	switch pt.Status.Stage {
	case task.STAGE_CANCELED:
		currentBuildConf.Failed(fmt.Errorf("流水线任务执行取消"))
	case task.STAGE_FAILED:
		currentBuildConf.Failed(fmt.Errorf(pt.Status.Message))
	case task.STAGE_SUCCEEDED:
		currentBuildConf.Success()
	}

	if err := i.UpdateRecordBuildConf(ctx, currentRecord.Event.Id, currentBuildConfIndex, currentBuildConf); err != nil {
		return nil, err
	}

	// 从BuildConf的队列中获取最近需要执行的一个
	if pt.Pipeline.TriggerNext() {
		// 避免构建同时触发, 更新时加锁, RunPipeline是 也加有BuildConfId锁，避免死锁
		m := lock.L().New("trigger_"+in.BuildConfId, 15*time.Second)
		if err := m.Lock(ctx); err != nil {
			return nil, exception.NewInternalServerError("lock error, %s", err)
		}
		defer m.UnLock(ctx)

		queueReq := trigger.NewQueryRecordRequest()
		req.AddBuildConfId(currentBuildConf.BuildConfig.Meta.Id)
		req.AddBuildStage(trigger.STAGE_ENQUEUE)
		req.IsOrderAscend = true
		req.Page.PageSize = 1
		rs, err = i.QueryRecord(ctx, queueReq)
		if err != nil {
			return nil, err
		}

		if len(rs.Items) == 0 {
			i.log.Debug().Msgf("build config %s not found enqueu record", currentBuildConf.BuildConfig.Meta.Id)
			return currentBuildConf, nil
		}

		// 获取Record里面 Next需要运行的 BuildConf
		nextRecord := rs.Items[0]
		nextBuildConf, nextBuildConfIndex := nextRecord.GetBuildStatusByBuildConfId(currentBuildConf.BuildConfig.Meta.Id)
		bs := i.RunBuildConf(ctx, nextRecord.Event, nextBuildConf.BuildConfig)
		if err := i.UpdateRecordBuildConf(ctx, nextRecord.Event.Id, nextBuildConfIndex, bs); err != nil {
			return nil, err
		}
	}

	return currentBuildConf, nil
}

// 删除执行记录
func (i *impl) DeleteRecord(ctx context.Context, in *trigger.DeleteRecordRequest) error {
	if err := in.Validate(); err != nil {
		return err
	}

	resp, err := i.col.DeleteMany(ctx, DeleteRecordFilter(in))
	if err != nil {
		return err
	}
	i.log.Info().Msgf("delete %d record", resp.DeletedCount)
	return nil
}

func DeleteRecordFilter(r *trigger.DeleteRecordRequest) bson.M {
	filter := bson.M{}

	switch r.DeleteBy {
	case trigger.DELETE_BY_PIPELINE_TASK_ID:
		filter["build_status.pipline_task_id"] = bson.M{"$in": r.Values}
	default:
		filter["_id"] = bson.M{"$in": r.Values}
	}
	return filter
}

// 查询事件详情
func (i *impl) DescribeRecord(
	ctx context.Context,
	in *trigger.DescribeRecordRequest) (
	*trigger.Record, error) {
	ins := trigger.NewDefaultRecord()
	if err := i.col.FindOne(ctx, bson.M{"_id": in.Id}).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("trigger event record %s not found", in.Id)
		}

		return nil, exception.NewInternalServerError("find trigger event record %s error, %s", in.Id, err)
	}
	return ins, nil
}

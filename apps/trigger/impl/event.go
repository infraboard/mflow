package impl

import (
	"context"
	"fmt"

	"github.com/infraboard/mcenter/apps/service"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/trigger"
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
	runReq.Labels[trigger.PIPELINE_TASK_EVENT_LABLE_KEY] = in.Id
	runReq.Labels[build.PIPELINE_TASK_BUILD_CONFIG_LABLE_KEY] = buildConf.Meta.Id
	runReq.Domain = buildConf.Spec.Scope.Domain
	runReq.Namespace = buildConf.Spec.Scope.Namespace

	// 补充Build用户自定义变量
	runReq.AddRunParam(buildConf.Spec.CustomParams...)

	// 补充Gitlab事件特有的变量
	switch in.Provider {
	case trigger.EVENT_PROVIDER_GITLAB:
		event, err := in.GetGitlabEvent()
		if err != nil {
			bs.ErrorMessage = err.Error()
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
		bs.ErrorMessage = err.Error()
	}
	if pt != nil {
		i.log.Debug().Msgf("update run build conf pipeline task id: %s", pt.Meta.Id)
		bs.PiplineTaskId = pt.Meta.Id
		bs.PiplineTask = pt
	}
	return bs
}

// 查询事件
func (i *impl) QueryRecord(ctx context.Context, in *trigger.QueryRecordRequest) (
	*trigger.RecordSet, error) {
	r := newQueryRequest(in)
	resp, err := i.col.Find(ctx, r.FindFilter(), r.FindOptions())

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

	// count
	count, err := i.col.CountDocuments(ctx, r.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get event record count error, error is %s", err)
	}
	set.Total = count
	return set, nil
}

// 事件队列任务执行完成通知
func (i *impl) EventQueueTaskComplete(ctx context.Context, in *trigger.EventQueueTaskCompleteRequest) (
	*trigger.BuildStatus, error) {
	req := trigger.NewQueryRecordRequest()
	req.PipelineTaskId = in.PipelineTaskId
	rs, err := i.QueryRecord(ctx, req)
	if err != nil {
		return nil, err
	}

	i.log.Debug().Msgf("%s", rs)
	return nil, nil
}

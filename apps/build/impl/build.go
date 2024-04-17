package impl

import (
	"context"
	"time"

	"dario.cat/mergo"
	"github.com/infraboard/mcenter/apps/service"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/pb/request"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (i *impl) CreateBuildConfig(ctx context.Context, in *build.CreateBuildConfigRequest) (
	*build.BuildConfig, error) {
	ins, err := build.New(in)
	if err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	if err := i.CheckBuildConfig(ctx, ins); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	if _, err := i.col.InsertOne(ctx, ins); err != nil {
		return nil, exception.NewInternalServerError("inserted a build document error, %s", err)
	}
	return ins, nil
}

func (i *impl) CheckBuildConfig(ctx context.Context, ins *build.BuildConfig) error {
	// 检查service id是否存在
	svc, err := i.mcenter.Service().DescribeService(ctx, service.NewDescribeServiceRequest(ins.Spec.ServiceId))
	if err != nil {
		return err
	}
	i.log.Debug().Msgf("found service: %s", svc.Spec.Name)
	return nil
}

func (i *impl) QueryBuildConfig(ctx context.Context, in *build.QueryBuildConfigRequest) (
	*build.BuildConfigSet, error) {
	r := newQueryRequest(in)
	resp, err := i.col.Find(ctx, r.FindFilter(), r.FindOptions())
	i.log.Debug().Msgf("query build filter: %s", r.FindFilter())

	if err != nil {
		return nil, exception.NewInternalServerError("find build error, error is %s", err)
	}

	set := build.NewBuildConfigSet()
	// 循环
	for resp.Next(ctx) {
		ins := build.NewDefaultBuildConfig()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode build error, error is %s", err)
		}
		set.Add(ins)
	}

	if in.WithPipeline {
		pReq := pipeline.NewQueryPipelineRequest()
		pReq.Ids = set.PipelineIds()
		pReq.Scope = in.Scope
		pset, err := i.pipeline.QueryPipeline(ctx, pReq)
		if err != nil {
			return nil, err
		}
		set.UpdatePipeline(pset.Items...)
	}

	// count
	count, err := i.col.CountDocuments(ctx, r.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get build count error, error is %s", err)
	}
	set.Total = count
	return set, nil
}

func (i *impl) DescribeBuildConfig(ctx context.Context, in *build.DescribeBuildConfigRequst) (
	*build.BuildConfig, error) {
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	ins := build.NewDefaultBuildConfig()
	if err := i.col.FindOne(ctx, bson.M{"_id": in.Id}).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("build config %s not found", in.Id)
		}

		return nil, exception.NewInternalServerError("find build config %s error, %s", in.Id, err)
	}
	return ins, nil
}

func (i *impl) UpdateBuildConfig(ctx context.Context, in *build.UpdateBuildConfigRequest) (
	*build.BuildConfig, error) {
	req := build.NewDescribeBuildConfigRequst(in.Id)
	ins, err := i.DescribeBuildConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	switch in.UpdateMode {
	case request.UpdateMode_PUT:
		ins.Spec = in.Spec
	case request.UpdateMode_PATCH:
		if err := mergo.MergeWithOverwrite(ins.Spec, in.Spec); err != nil {
			return nil, err
		}
	default:
		return nil, exception.NewBadRequest("unknown update mode: %s", in.UpdateMode)
	}

	if err := ins.Spec.Validate(); err != nil {
		return nil, err
	}

	ins.Meta.UpdateAt = time.Now().Unix()
	ins.Meta.UpdateBy = in.UpdateBy
	_, err = i.col.UpdateOne(ctx, bson.M{"_id": ins.Meta.Id}, bson.M{"$set": ins})
	if err != nil {
		return nil, exception.NewInternalServerError("update build config(%s) error, %s", ins.Meta.Id, err)
	}

	return ins, nil
}

func (i *impl) DeleteBuildConfig(ctx context.Context, in *build.DeleteBuildConfigRequest) (
	*build.BuildConfig, error) {
	req := build.NewDescribeBuildConfigRequst(in.Id)
	ins, err := i.DescribeBuildConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	_, err = i.col.DeleteOne(ctx, bson.M{"_id": in.Id})
	if err != nil {
		return nil, exception.NewInternalServerError("delete build config(%s) error, %s", in.Id, err)
	}
	return ins, nil
}

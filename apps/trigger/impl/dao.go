package impl

import (
	"context"

	"github.com/infraboard/mflow/apps/trigger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newQueryRequest(r *trigger.QueryRecordRequest) *queryRequest {
	return &queryRequest{
		r,
	}
}

type queryRequest struct {
	*trigger.QueryRecordRequest
}

func (r *queryRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.Page.PageSize)
	skip := int64(r.Page.PageSize) * int64(r.Page.PageNumber-1)

	opt := &options.FindOptions{
		Sort: bson.D{
			{Key: "time", Value: -1},
		},
		Limit: &pageSize,
		Skip:  &skip,
	}
	return opt
}

func (r *queryRequest) FindFilter() bson.M {
	filter := bson.M{}

	if r.ServiceId != "" {
		filter["token"] = r.ServiceId
	}

	if r.PipelineTaskId != "" {
		filter["build_status.pipeline_task_id"] = r.PipelineTaskId
	}

	if len(r.BuildConfIds) > 0 {
		filter["build_status.build_config._id"] = bson.M{"$in": r.BuildConfIds}
	}

	return filter
}

// 查询事件详情
func (i *impl) UpdateRecordBuildConf(
	ctx context.Context, recordId, buildConfId string, stage trigger.STAGE) error {
	filter := bson.M{"_id": recordId, "build_status.build_config._id": buildConfId}
	if _, err := i.col.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"build_status.$.stage": stage,
	}}); err != nil {
		return err
	}
	return nil
}

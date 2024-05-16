package impl

import (
	"context"
	"fmt"

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

	r.AddBuildStage()

	opt := &options.FindOptions{
		Sort: bson.D{
			{Key: "time", Value: r.SortValue()},
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
		filter["build_status.pipline_task_id"] = r.PipelineTaskId
	}
	if len(r.BuildConfIds) > 0 {
		filter["build_status.build_config._id"] = bson.M{"$in": r.BuildConfIds}
	}
	if len(r.BuildStages) > 0 {
		filter["build_status.stage"] = bson.M{"$in": r.BuildStages}
	}

	return filter
}

// 查询事件详情
func (i *impl) UpdateRecordBuildConf(
	ctx context.Context, recordId string, buildIndex int, bc *trigger.BuildStatus) error {
	filter := bson.M{"_id": recordId}
	if _, err := i.col.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		fmt.Sprintf("build_status.%d", buildIndex): bc,
	}}); err != nil {
		return err
	}
	return nil
}

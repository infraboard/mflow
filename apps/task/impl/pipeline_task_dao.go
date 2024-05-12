package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcenter/apps/policy"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mflow/apps/task"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newQueryPipelineTaskRequest(r *task.QueryPipelineTaskRequest) *queryPipelineTaskRequest {
	return &queryPipelineTaskRequest{
		r,
	}
}

type queryPipelineTaskRequest struct {
	*task.QueryPipelineTaskRequest
}

func (r *queryPipelineTaskRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.Page.PageSize)
	skip := int64(r.Page.PageSize) * int64(r.Page.PageNumber-1)

	opt := &options.FindOptions{
		Sort: bson.D{
			{Key: "create_at", Value: -1},
		},
		Limit: &pageSize,
		Skip:  &skip,
	}

	return opt
}

func (r *queryPipelineTaskRequest) FindFilter() bson.M {
	filter := bson.M{}
	token.MakeMongoFilter(filter, r.Scope)
	policy.MakeMongoFilter(filter, "labels", r.Filters)

	if len(r.Ids) > 0 {
		filter["_id"] = bson.M{"$in": r.Ids}
	}

	if r.PipelineId != "" {
		filter["pipeline._id"] = r.PipelineId
	}

	if len(r.Stages) > 0 {
		filter["status.stage"] = bson.M{"$in": r.Stages}
	}
	// 补充标签
	for k, v := range r.Labels {
		filter[fmt.Sprintf("labels.%s", k)] = v
	}

	return filter
}

func (i *impl) deletePipelineTask(ctx context.Context, ins *task.PipelineTask) error {
	if ins == nil || ins.Meta.Id == "" {
		return fmt.Errorf("pipeline is nil")
	}

	result, err := i.pcol.DeleteOne(ctx, bson.M{"_id": ins.Meta.Id})
	if err != nil {
		return exception.NewInternalServerError("delete pipeline task(%s) error, %s", ins.Meta.Id, err)
	}

	if result.DeletedCount == 0 {
		return exception.NewNotFound("pipeline task %s not found", ins.Meta.Id)
	}

	return nil
}

func (i *impl) updatePiplineTaskStatus(ctx context.Context, in *task.PipelineTask) error {
	in.Meta.UpdateAt = time.Now().Unix()
	if _, err := i.pcol.UpdateByID(ctx, in.Meta.Id, bson.M{"$set": bson.M{"status": in.Status}}); err != nil {
		return exception.NewInternalServerError("update task(%s) document error, %s",
			in.Meta.Id, err)
	}

	return nil
}

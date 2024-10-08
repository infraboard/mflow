package impl

import (
	"context"

	"github.com/infraboard/mcenter/apps/policy"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newQueryRequest(r *task.QueryJobTaskRequest) *queryRequest {
	return &queryRequest{
		r,
	}
}

type queryRequest struct {
	*task.QueryJobTaskRequest
}

func (r *queryRequest) SortValue() int {
	switch r.SortType {
	case task.SORT_TYPE_ASCEND:
		return 1
	case task.SORT_TYPE_DESCEND:
		return -1
	}

	return -1
}

func (r *queryRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.Page.PageSize)
	skip := int64(r.Page.PageSize) * int64(r.Page.PageNumber-1)

	opt := &options.FindOptions{
		Sort: bson.D{
			{Key: "create_at", Value: r.SortValue()},
		},
		Limit: &pageSize,
		Skip:  &skip,
	}

	return opt
}

func (r *queryRequest) FindFilter() bson.M {
	filter := bson.M{}
	token.MakeMongoFilter(filter, r.Scope)
	policy.MakeMongoFilter(filter, "labels", r.Filters)

	if len(r.Ids) > 0 {
		filter["_id"] = bson.M{"$in": r.Ids}
	}

	if r.PipelineTaskId != "" {
		filter["pipeline_task"] = r.PipelineTaskId
	}

	if r.JobId != "" {
		filter["job_id"] = r.JobId
	}

	if r.Stage != nil {
		filter["status.stage"] = *r.Stage
	}
	if r.Auditor != "" {
		filter["audit.auditors"] = r.Auditor
	}
	if r.AuditEnable != nil {
		filter["audit.enable"] = *r.AuditEnable
	}
	if len(r.AuditStages) > 0 {
		filter["audit.status.stage"] = bson.M{"$in": r.AuditStages}
	}

	if r.HasLabel() {
		for k, v := range r.Labels {
			filter["labels."+k] = v
		}
	}

	return filter
}

func (i *impl) updateJobTaskStatus(ctx context.Context, t *task.JobTask) error {
	i.log.Debug().Msgf("任务[%s] 执行状态: %s 审核状态: %s", t.Spec.TaskId, t.Status.Stage, t.AuditStatus())

	// 更新数据库
	if _, err := i.jcol.UpdateByID(ctx, t.Spec.TaskId, bson.M{"$set": bson.M{"status": t.Status}}); err != nil {
		return exception.NewInternalServerError("update task(%s) document error, %s",
			t.Spec.TaskId, err)
	}
	return nil
}

func (i *impl) updateJobTaskImRobotNotify(ctx context.Context, taskId string, t []*pipeline.WebHook) {
	if _, err := i.jcol.UpdateByID(ctx, taskId, bson.M{"$set": bson.M{"im_robot_notify": t}}); err != nil {
		i.log.Error().Msgf("update task %s im_robot_notify error, %s", taskId, err)
	}
}

func (i *impl) updateJobTaskWebHook(ctx context.Context, taskId string, t []*pipeline.WebHook) {
	if _, err := i.jcol.UpdateByID(ctx, taskId, bson.M{"$set": bson.M{"webhooks": t}}); err != nil {
		i.log.Error().Msgf("update task %s webhooks error, %s", taskId, err)
	}
}

func (i *impl) updateJobTaskMentionUser(ctx context.Context, taskId string, t []*pipeline.MentionUser) {
	if _, err := i.jcol.UpdateByID(ctx, taskId, bson.M{"$set": bson.M{"mention_users": t}}); err != nil {
		i.log.Error().Msgf("update task %s mention_users error, %s", taskId, err)
	}
}

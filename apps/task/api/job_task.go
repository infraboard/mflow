package api

import (
	"time"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
)

var (
	JOB_TASK_RESOURCE_NAME = "JobTask"
)

// 用户自己手动管理任务状态相关接口
func (h *JobTaskHandler) RegistryUserHandler() {
	tags := []string{"Job任务管理"}

	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.GET("/").
		To(h.QueryJobTask).
		Doc("查询JobTask运行任务列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, JOB_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(task.QueryJobTaskRequest{}).
		Writes(task.JobTaskSet{}))

	ws.Route(ws.POST("/").
		To(h.RunJob).
		Doc("运行JobTask").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, JOB_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Create.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(pipeline.Task{}).
		Writes(task.JobTask{}))

	ws.Route(ws.GET("/{id}").
		To(h.DescribeJobTask).
		Doc("查询JobTask运行任务详情").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, JOB_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(task.DescribeJobTaskRequest{}).
		Writes(task.JobTask{}))

	// 内部通过审核人列表判断权限
	ws.Route(ws.POST("/{id}/audit").
		To(h.AuditJobTask).
		Doc("任务审核").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, JOB_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Update.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Disable).
		Reads(task.AuditJobTaskRequest{}).
		Writes(task.JobTask{}))

	// 通过Job自身的Token进行认证
	ws.Route(ws.POST("/{id}/status").
		To(h.UpdateJobTaskStatus).
		Doc("更新任务状态").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, JOB_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Create.Value()).
		Reads(task.UpdateJobTaskStatusRequest{}).
		Writes(task.JobTask{}))

	// 通过Job自身的Token进行认证
	ws.Route(ws.POST("/{id}/output").
		To(h.UpdateJobTaskOutput).
		Doc("更新任务输出").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, JOB_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Create.Value()).
		Reads(task.UpdateJobTaskOutputRequest{}).
		Writes(task.JobTask{}))
}

func (h *JobTaskHandler) QueryJobTask(r *restful.Request, w *restful.Response) {
	req := task.NewQueryTaskRequestFromHttp(r)
	set, err := h.service.QueryJobTask(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *JobTaskHandler) DescribeJobTask(r *restful.Request, w *restful.Response) {
	req := task.NewDescribeJobTaskRequest(r.PathParameter("id"))
	ins, err := h.service.DescribeJobTask(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

func (h *JobTaskHandler) RunJob(r *restful.Request, w *restful.Response) {
	req := pipeline.NewTask("")
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}

	req.UpdateFromToken(token.GetTokenFromRequest(r))
	set, err := h.service.RunJob(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *JobTaskHandler) UpdateJobTaskOutput(r *restful.Request, w *restful.Response) {
	req := task.NewUpdateJobTaskOutputRequest("")
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}
	req.Id = r.PathParameter("id")

	set, err := h.service.UpdateJobTaskOutput(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *JobTaskHandler) UpdateJobTaskStatus(r *restful.Request, w *restful.Response) {
	req := task.NewUpdateJobTaskStatusRequest("")
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}
	req.Id = r.PathParameter("id")

	set, err := h.service.UpdateJobTaskStatus(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *JobTaskHandler) AuditJobTask(r *restful.Request, w *restful.Response) {
	req := task.NewAuditJobTaskRequest(r.PathParameter("id"))
	if err := r.ReadEntity(req.Status); err != nil {
		response.Failed(w, err)
		return
	}

	tk := token.GetTokenFromRequest(r)
	req.Status.AuditBy = tk.UserId
	req.Status.AuditAt = time.Now().Unix()

	set, err := h.service.AuditJobTask(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

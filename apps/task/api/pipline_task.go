package api

import (
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
	PIPELINE_TASK_RESOURCE_NAME = "PipelineTask"
)

// 用户自己手动管理任务状态相关接口
func (h *PipelineTaskHandler) RegistryUserHandler() {
	tags := []string{"Pipeline任务管理"}

	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.GET("/").To(h.QueryPipelineTask).
		Doc("查询Pipeline运行任务列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, PIPELINE_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(task.QueryPipelineTaskRequest{}).
		Writes(task.PipelineTaskSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribePipelineTask).
		Doc("查询Pipeline运行任务详情").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, PIPELINE_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(task.DescribePipelineTaskRequest{}).
		Writes(task.PipelineTask{}))

	ws.Route(ws.POST("/").To(h.RunPipeline).
		Doc("运行Pipeline").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, PIPELINE_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Create.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(pipeline.RunPipelineRequest{}).
		Writes(task.PipelineTask{}))

	ws.Route(ws.DELETE("/{id}").To(h.DeletePipelineTask).
		Doc("删除Pipeline运行任务").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, PIPELINE_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Writes(task.PipelineTask{}))
}

func (h *PipelineTaskHandler) QueryPipelineTask(r *restful.Request, w *restful.Response) {
	req := task.NewQueryPipelineTaskRequestFromHttp(r)
	set, err := h.service.QueryPipelineTask(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *PipelineTaskHandler) DescribePipelineTask(r *restful.Request, w *restful.Response) {
	req := task.NewDescribePipelineTaskRequest(r.PathParameter("id"))
	set, err := h.service.DescribePipelineTask(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *PipelineTaskHandler) RunPipeline(r *restful.Request, w *restful.Response) {
	req := pipeline.NewRunPipelineRequest("")
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}
	req.UpdateFromToken(token.GetTokenFromRequest(r))
	set, err := h.service.RunPipeline(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *PipelineTaskHandler) DeletePipelineTask(r *restful.Request, w *restful.Response) {
	req := task.NewDeletePipelineTaskRequest(r.PathParameter("id"))
	ins, err := h.service.DeletePipelineTask(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

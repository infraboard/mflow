package api

import (
	"fmt"

	"github.com/infraboard/mcenter/apps/token"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mflow/apps/job"
)

func (h *handler) Registry() {
	tags := []string{"Job管理"}

	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.POST("/").To(h.CreateJob).
		Doc("创建Job").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Create.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(job.CreateJobRequest{}).
		Writes(job.Job{}))

	ws.Route(ws.GET("/").To(h.QueryJob).
		Doc("查询Job列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(job.QueryJobRequest{}).
		Writes(job.JobSet{}).
		Returns(200, "OK", job.JobSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribeJob).
		Doc("Job详情").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Disable).
		Metadata(label.Permission, label.Disable).
		Writes(job.Job{}).
		Returns(200, "OK", job.Job{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{id}").To(h.PutJob).
		Doc("修改Job").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(job.CreateJobRequest{}).
		Writes(job.Job{}).
		Returns(200, "OK", job.Job{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PATCH("/{id}").To(h.PatchJob).
		Doc("修改Job").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(job.CreateJobRequest{}).
		Writes(job.Job{}).
		Returns(200, "OK", job.Job{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{id}/status").To(h.UpdateJobStatus).
		Doc("修改Job状态").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(job.JobStatus{}).
		Writes(job.Job{}).
		Returns(200, "OK", job.Job{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.DELETE("/{id}").To(h.DeleteJob).
		Doc("删除Job").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Delete.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable))
}

func (h *handler) CreateJob(r *restful.Request, w *restful.Response) {
	req := job.NewCreateJobRequest()

	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}

	// 补充Scope信息
	tk := token.GetTokenFromRequest(r)
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace
	req.CreateBy = fmt.Sprintf("%s@%s", tk.Username, tk.Domain)

	ins, err := h.service.CreateJob(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

func (h *handler) QueryJob(r *restful.Request, w *restful.Response) {
	req := job.NewQueryJobRequestFromHTTP(r)

	set, err := h.service.QueryJob(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}

func (h *handler) DescribeJob(r *restful.Request, w *restful.Response) {
	req := job.NewDescribeJobRequestById(r.PathParameter("id"))
	req.Scope = token.GetTokenFromRequest(r).GenScope()

	// 解析describe by
	dbstr := r.QueryParameter("describe_by")
	if dbstr != "" {
		by, err := job.ParseDESCRIBE_BYFromString(dbstr)
		if err != nil {
			response.Failed(w, err)
			return
		}
		req.DescribeBy = by
	}

	ins, err := h.service.DescribeJob(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) PutJob(r *restful.Request, w *restful.Response) {
	tk := r.Attribute("token").(*token.Token)

	req := job.NewPutJobRequest(r.PathParameter("id"))
	if err := r.ReadEntity(req.Spec); err != nil {
		response.Failed(w, err)
		return
	}
	req.UpdateBy = tk.Username

	set, err := h.service.UpdateJob(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) PatchJob(r *restful.Request, w *restful.Response) {
	tk := r.Attribute("token").(*token.Token)

	req := job.NewPatchJobRequest(r.PathParameter("id"))
	if err := r.ReadEntity(req.Spec); err != nil {
		response.Failed(w, err)
		return
	}
	req.UpdateBy = tk.Username

	set, err := h.service.UpdateJob(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) UpdateJobStatus(r *restful.Request, w *restful.Response) {
	tk := r.Attribute("token").(*token.Token)

	req := job.NewUpdateJobStatusRequest(fmt.Sprintf("#%s", r.PathParameter("id")))
	if err := r.ReadEntity(req.Status); err != nil {
		response.Failed(w, err)
		return
	}
	req.Status.ChangeBy = tk.Username

	set, err := h.service.UpdateJobStatus(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) DeleteJob(r *restful.Request, w *restful.Response) {
	req := job.NewDeleteJobRequest(r.PathParameter("id"))
	set, err := h.service.DeleteJob(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

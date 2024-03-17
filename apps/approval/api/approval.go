package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mflow/apps/approval"
)

func init() {
	ioc.Api().Registry(&handler{})
}

type handler struct {
	service approval.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *handler) Init() error {
	h.log = log.Sub(approval.AppName)
	h.service = ioc.Controller().Get(approval.AppName).(approval.Service)
	h.Registry()
	return nil
}

// /prifix/cluster/
func (h *handler) Name() string {
	return approval.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry() {
	tags := []string{"审核管理"}

	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.POST("/").To(h.CreateApproval).
		Doc("创建审核").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Create.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(approval.CreateApprovalRequest{}).
		Writes(approval.Approval{}))

	ws.Route(ws.GET("/").To(h.QueryApproval).
		Doc("查询审核").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Disable).
		Reads(approval.QueryApprovalRequest{}).
		Writes(approval.ApprovalSet{}).
		Returns(200, "OK", approval.ApprovalSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribeApproval).
		Doc("审核详情").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Disable).
		Writes(approval.Approval{}).
		Returns(200, "OK", approval.Approval{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{id}").To(h.EditApproval).
		Doc("编辑审核").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Writes(approval.Approval{}).
		Returns(200, "OK", approval.Approval{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/{id}").To(h.UpdateApprovalStatus).
		Doc("提交审核").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(approval.UpdateApprovalStatusRequest{}).
		Writes(approval.Approval{}).
		Returns(200, "OK", approval.Approval{}).
		Returns(404, "Not Found", nil))
}

func (h *handler) CreateApproval(r *restful.Request, w *restful.Response) {
	req := approval.NewCreateApprovalRequest()
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}

	req.UpdateFromToken(token.GetTokenFromRequest(r))
	ins, err := h.service.CreateApproval(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

func (h *handler) QueryApproval(r *restful.Request, w *restful.Response) {
	req := approval.NewQueryApprovalRequestFromHTTP(r)
	set, err := h.service.QueryApproval(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) DescribeApproval(r *restful.Request, w *restful.Response) {
	req := approval.NewDescribeApprovalRequest(r.PathParameter("id"))
	ins, err := h.service.DescribeApproval(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) EditApproval(r *restful.Request, w *restful.Response) {
	req := approval.NewEditApprovalRequest(r.PathParameter("id"))

	if err := r.ReadEntity(req.Spec); err != nil {
		response.Failed(w, err)
		return
	}

	ins, err := h.service.EditApproval(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) UpdateApprovalStatus(r *restful.Request, w *restful.Response) {
	req := approval.NewUpdateApprovalStatusRequest(r.PathParameter("id"))

	if err := r.ReadEntity(req.Status); err != nil {
		response.Failed(w, err)
		return
	}

	ins, err := h.service.UpdateApprovalStatus(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcenter/apps/user"
	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/approval"
	"github.com/infraboard/mflow/apps/approval/api/callback"
)

func init() {
	ioc.Api().Registry(&callbackHandler{})
}

type callbackHandler struct {
	service approval.Service
	log     *zerolog.Logger
	mcenter *rpc.ClientSet
	ioc.ObjectImpl
}

func (h *callbackHandler) Init() error {
	h.log = logger.Sub(approval.AppName)
	h.service = ioc.Controller().Get(approval.AppName).(approval.Service)
	h.mcenter = rpc.C()
	return nil
}

func (h *callbackHandler) Name() string {
	return "callbacks"
}

func (h *callbackHandler) Version() string {
	return "v1"
}

func (h *callbackHandler) Registry(ws *restful.WebService) {
	tags := []string{"审核对接第三方"}
	ws.Route(ws.POST("/feishu").To(h.FeishuCard).
		Doc("飞书卡片处理审核").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Create.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(approval.CreateApprovalRequest{}).
		Writes(approval.Approval{}))
}

// 飞书卡片回调, 参考文档: https://open.feishu.cn/document/ukTMukTMukTM/uYjNwUjL2YDM14iN2ATN
func (h *callbackHandler) FeishuCard(r *restful.Request, w *restful.Response) {
	req := callback.NewFeishuCardCallback()
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}

	// 根据飞书UserId查询用户信息
	descReq := user.NewDescriptUserRequestByFeishuUserId(req.UserId)
	u, err := h.mcenter.User().DescribeUser(r.Request.Context(), descReq)
	if err != nil {
		response.Failed(w, err)
		return
	}

	in, err := req.BuildUpdateApprovalStatusRequest(u.Meta.Id)
	if err != nil {
		response.Failed(w, err)
		return
	}
	ins, err := h.service.UpdateApprovalStatus(r.Request.Context(), in)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 返回更新后的卡片信息
	msg, _ := ins.FeishuAuditNotifyMessage()
	sendReq, err := msg.BuildNotifyRequest()
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, sendReq.Content)
}

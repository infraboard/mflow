package api

import (
	"fmt"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	http "github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mflow/apps/trigger"
)

func init() {
	ioc.Api().Registry(&Handler{})
}

type Handler struct {
	log *zerolog.Logger
	svc trigger.Service
	ioc.ObjectImpl

	mcenter *rpc.ClientSet
}

func (h *Handler) Init() error {
	h.svc = ioc.Controller().Get(trigger.AppName).(trigger.Service)
	h.log = log.Sub(trigger.AppName)
	h.mcenter = rpc.C()
	h.Registry()
	return nil
}

func (h *Handler) Name() string {
	return trigger.AppName
}

func (h *Handler) Version() string {
	return "v1"
}

func (h *Handler) APIPrefix() string {
	return fmt.Sprintf("%s/%s/%s",
		http.Get().HTTPPrefix(),
		h.Version(),
		h.Name(),
	)
}

func (h *Handler) Registry() {
	tags := []string{"事件处理"}

	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.GET("records").To(h.QueryRecord).
		Doc("查询触发记录").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, "trigger_records").
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable))
	ws.Route(ws.POST("gitlab").To(h.HandleGitlabEvent).
		Doc("处理Gitlab Webhook事件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Create.Value()).
		Metadata(label.Auth, label.Disable).
		Metadata(label.Permission, label.Disable))
	ws.Route(ws.POST("mannul").To(h.MannulGitlabEvent).
		Doc("手动模拟Gitlab Webhook事件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Create.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable))
}

package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mflow/apps/trigger"
)

// 处理来自gitlab的事件
// Hook Header参考文档: https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#delivery-headers
// 参考文档: https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html
func (h *Handler) HandleGitlabEvent(r *restful.Request, w *restful.Response) {
	event, err := trigger.ParseGitLabEventFromRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 专门处理Gilab事件
	ge, err := event.GetGitlabEvent()
	if err != nil {
		response.Failed(w, err)
		return
	}
	event.SubName = ge.GetBranch()

	h.log.Debug().Msgf("accept event: %s", event.ToJson())
	ins, err := h.svc.HandleEvent(r.Request.Context(), event)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

// 查询repo 的gitlab地址, 手动获取信息, 触发手动事件
func (h *Handler) MannulGitlabEvent(r *restful.Request, w *restful.Response) {
	// 构造事件
	event := trigger.NewGitlabEvent("")

	// 读取模拟事件
	err := r.ReadEntity(event)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 事件关联信息填充
	err = h.BuildEvent(r.Request.Context(), event)
	if err != nil {
		response.Failed(w, err)
		return
	}

	h.log.Debug().Msgf("mannul event: %s", event.Raw)
	ins, err := h.svc.HandleEvent(r.Request.Context(), event)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

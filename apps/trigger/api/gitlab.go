package api

import (
	"context"
	"fmt"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcenter/apps/service"
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

	h.log.Debug().Msgf("mannul event: %s", event.ToJson())
	ins, err := h.svc.HandleEvent(r.Request.Context(), event)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

func (h *Handler) BuildEvent(ctx context.Context, in *trigger.Event) error {
	in.IsMock = true

	// 查询服务仓库信息
	descReq := service.NewDescribeServiceRequest(in.Token)
	svc, err := h.mcenter.Service().DescribeService(ctx, descReq)
	if err != nil {
		return err
	}
	h.log.Debug().Msgf("service: %s", svc)
	repo := svc.Spec.CodeRepository
	if repo == nil || repo.Token == "" {
		return fmt.Errorf("service %s[%s] no repo or private token info", svc.FullName(), svc.Meta.Id)
	}

	event, err := in.GetGitlabEvent()
	if err != nil {
		return err
	}

	// 补充Project相关信息
	p := event.Project
	p.Id = repo.ProjectIdToInt64()
	p.GitHttpUrl = repo.HttpUrl
	p.GitSshUrl = repo.SshUrl
	p.NamespacePath = repo.Namespace
	p.WebUrl = repo.WebUrl
	p.Name = svc.Spec.Name

	// 补充分支相关Commit信息
	// gc, err := repo.MakeGitlabConfig()
	// if err != nil {
	// 	return err
	// }
	// v4 := gitlab.NewGitlabV4(gc)
	// branchReq := gitlab.NewGetProjectBranchRequest()
	// branchReq.ProjectId = repo.ProjectId
	// branchReq.Branch = event.GetBranch()
	// if branchReq.Branch == "" {
	// 	branchReq.Branch = in.SubName
	// }
	// b, err := v4.Project().GetProjectBranch(ctx, branchReq)
	// if err != nil {
	// 	return fmt.Errorf("查询分支: %s 异常, %s", branchReq.Branch, err)
	// }
	// event.Commits = append(event.Commits, &trigger.Commit{
	// 	Id:        b.Commit.Id,
	// 	Message:   b.Commit.Message,
	// 	Title:     b.Commit.Title,
	// 	Timestamp: b.Commit.CommittedDate,
	// })
	return nil
}

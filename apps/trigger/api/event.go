package api

import (
	"context"
	"fmt"

	"github.com/infraboard/mcenter/apps/service"
	"github.com/infraboard/mflow/apps/trigger"
)

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
	return nil
}

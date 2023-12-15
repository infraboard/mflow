package api

import (
	"fmt"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/task"
)

func init() {
	ioc.Api().Registry(&JobTaskHandler{})
}

type JobTaskHandler struct {
	service task.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *JobTaskHandler) Init() error {
	h.log = logger.Sub(task.AppName)
	h.service = ioc.Controller().Get(task.AppName).(task.Service)
	return nil
}

func (h *JobTaskHandler) Name() string {
	return "job_tasks"
}

func (h *JobTaskHandler) Version() string {
	return "v1"
}

func (h *JobTaskHandler) APIPrefix() string {
	return fmt.Sprintf("%s/%s/%s",
		application.App().HTTPPrefix(),
		h.Version(),
		h.Name(),
	)
}

func (h *JobTaskHandler) Registry(ws *restful.WebService) {
	h.RegistryUserHandler(ws)
}

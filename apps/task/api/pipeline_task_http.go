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
	ioc.RegistryApi(&PipelineTaskHandler{})
}

type PipelineTaskHandler struct {
	service task.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *PipelineTaskHandler) Init() error {
	h.log = logger.Sub(task.AppName)
	h.service = ioc.GetController(task.AppName).(task.Service)
	return nil
}

func (h *PipelineTaskHandler) Name() string {
	return "pipeline_tasks"
}

func (h *PipelineTaskHandler) Version() string {
	return "v1"
}

func (h *PipelineTaskHandler) APIPrefix() string {
	return fmt.Sprintf("%s/%s/%s",
		application.App().HTTPPrefix(),
		h.Version(),
		h.Name(),
	)
}

func (h *PipelineTaskHandler) Registry(ws *restful.WebService) {
	h.RegistryUserHandler(ws)
}

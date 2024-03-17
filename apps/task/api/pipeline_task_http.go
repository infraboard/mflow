package api

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/task"
)

func init() {
	ioc.Api().Registry(&PipelineTaskHandler{})
}

type PipelineTaskHandler struct {
	service task.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *PipelineTaskHandler) Init() error {
	h.log = log.Sub(task.AppName)
	h.service = ioc.Controller().Get(task.AppName).(task.Service)
	h.Registry()
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
		http.Get().HTTPPrefix(),
		h.Version(),
		h.Name(),
	)
}

func (h *PipelineTaskHandler) Registry() {
	h.RegistryUserHandler()
}

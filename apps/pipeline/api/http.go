package api

import (
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/pipeline"
)

func init() {
	ioc.RegistryApi(&handler{})
}

type handler struct {
	service pipeline.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *handler) Init() error {
	h.log = logger.Sub(pipeline.AppName)
	h.service = ioc.GetController(pipeline.AppName).(pipeline.Service)
	return nil
}

func (h *handler) Name() string {
	return pipeline.AppName
}

func (h *handler) Version() string {
	return "v1"
}

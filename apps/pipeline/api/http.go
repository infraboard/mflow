package api

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/pipeline"
)

func init() {
	ioc.Api().Registry(&handler{})
}

type handler struct {
	service pipeline.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *handler) Init() error {
	h.log = logger.Sub(pipeline.AppName)
	h.service = ioc.Controller().Get(pipeline.AppName).(pipeline.Service)
	return nil
}

func (h *handler) Name() string {
	return pipeline.AppName
}

func (h *handler) Version() string {
	return "v1"
}

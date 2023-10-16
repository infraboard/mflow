package api

import (
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/build"
)

func init() {
	ioc.RegistryApi(&handler{})
}

type handler struct {
	service build.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *handler) Init() error {
	h.log = logger.Sub(build.AppName)
	h.service = ioc.GetController(build.AppName).(build.Service)
	return nil
}

func (h *handler) Name() string {
	return build.AppName
}

func (h *handler) Version() string {
	return "v1"
}

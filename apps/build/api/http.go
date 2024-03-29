package api

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/build"
)

func init() {
	ioc.Api().Registry(&handler{})
}

type handler struct {
	service build.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *handler) Init() error {
	h.log = log.Sub(build.AppName)
	h.service = ioc.Controller().Get(build.AppName).(build.Service)
	h.Registry()
	return nil
}

func (h *handler) Name() string {
	return build.AppName
}

func (h *handler) Version() string {
	return "v1"
}

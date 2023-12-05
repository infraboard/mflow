package api

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/job"
)

func init() {
	ioc.RegistryApi(&handler{})
}

type handler struct {
	service job.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *handler) Init() error {
	h.log = logger.Sub(job.AppName)
	h.service = ioc.GetController(job.AppName).(job.Service)
	return nil
}

func (h *handler) Name() string {
	return job.AppName
}

func (h *handler) Version() string {
	return "v1"
}

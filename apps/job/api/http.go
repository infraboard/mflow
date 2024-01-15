package api

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/job"
)

func init() {
	ioc.Api().Registry(&handler{})
}

type handler struct {
	service job.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *handler) Init() error {
	h.log = log.Sub(job.AppName)
	h.service = ioc.Controller().Get(job.AppName).(job.Service)
	return nil
}

func (h *handler) Name() string {
	return job.AppName
}

func (h *handler) Version() string {
	return "v1"
}

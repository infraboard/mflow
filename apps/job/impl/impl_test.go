package impl_test

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl job.Service
	ctx  = context.Background()
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.GetController(job.AppName).(job.Service)
}

package impl_test

import (
	"context"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl pipeline.Service
	ctx  = context.Background()
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.GetController(pipeline.AppName).(pipeline.Service)
}

package impl_test

import (
	"context"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl build.Service
	ctx  = context.Background()
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.GetController(build.AppName).(build.Service)
}

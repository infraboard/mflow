package impl_test

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl build.Service
	ctx  = context.Background()
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.Controller().Get(build.AppName).(build.Service)
}

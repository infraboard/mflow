package impl_test

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mflow/apps/trigger"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl trigger.Service
	ctx  = context.Background()
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.Controller().Get(trigger.AppName).(trigger.Service)
}

package impl_test

import (
	"context"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mflow/apps/approval"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl approval.Service
	ctx  = context.Background()
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.GetController(approval.AppName).(approval.Service)
}

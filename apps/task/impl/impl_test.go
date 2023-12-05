package impl_test

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl task.Service
	ctx  = context.Background()
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.GetController(task.AppName).(task.Service)
}

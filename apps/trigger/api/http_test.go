package api_test

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mflow/apps/trigger"
	"github.com/infraboard/mflow/apps/trigger/api"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl *api.Handler
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.Api().Get(trigger.AppName).(*api.Handler)
}

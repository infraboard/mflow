package api_test

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mflow/apps/task/api"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl *api.WebsocketHandler
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.Api().Get("ws/task").(*api.WebsocketHandler)
}

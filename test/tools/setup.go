package tools

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/ioc"
	test "github.com/infraboard/mflow/test/conf"

	// 注册所有服务
	_ "github.com/infraboard/mflow/apps"
)

func DevelopmentSetup() {
	// 针对http handler的测试需要提前设置默认数据格式
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// 加载单元测试的变量
	test.LoadConfigFromEnv()

	// 初始化全局app
	ioc.DevelopmentSetup()
}

package task

import (
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/application"
)

var (
	ERR_PIPELINE_IS_RUNNING = exception.NewAPIException(7001, "流水线当前处于运行中, 请取消或者等待之前运行结束").WithNamespace(application.Get().AppName)
)

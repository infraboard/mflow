package apps

import (
	// 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载
	_ "github.com/infraboard/mflow/apps/approval/api"
	_ "github.com/infraboard/mflow/apps/build/api"
	_ "github.com/infraboard/mflow/apps/job/api"
	_ "github.com/infraboard/mflow/apps/pipeline/api"
	_ "github.com/infraboard/mflow/apps/task/api"
	_ "github.com/infraboard/mflow/apps/trigger/api"
)

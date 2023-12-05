package apps

import (
	// 注册所有GRPC服务模块, 暴露给框架GRPC服务器加载, 注意 导入有先后顺序
	_ "github.com/infraboard/mflow/apps/approval/impl"
	_ "github.com/infraboard/mflow/apps/build/impl"
	_ "github.com/infraboard/mflow/apps/job/impl"
	_ "github.com/infraboard/mflow/apps/pipeline/impl"
	_ "github.com/infraboard/mflow/apps/task/impl"
	_ "github.com/infraboard/mflow/apps/trigger/impl"

	// 引入第三方存储模块
	_ "github.com/infraboard/mcube/v2/ioc/apps/oss/mongo"
)

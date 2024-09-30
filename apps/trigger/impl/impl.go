package impl

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/apps/trigger"
)

func init() {
	ioc.Controller().Registry(&impl{})
}

type impl struct {
	col *mongo.Collection
	log *zerolog.Logger
	trigger.UnimplementedRPCServer
	ioc.ObjectImpl

	// 构建配置
	build build.Service
	// 执行流水线
	task task.PipelineTaskService
	// mcenter客户端
	mcenter *rpc.ClientSet
}

func (i *impl) Init() error {
	i.col = ioc_mongo.DB().Collection(i.Name())
	i.log = log.Sub(i.Name())

	i.build = ioc.Controller().Get(build.AppName).(build.Service)
	i.task = ioc.Controller().Get(task.AppName).(task.Service)
	i.mcenter = rpc.C()

	trigger.RegisterRPCServer(grpc.Get().Server(), i)
	return nil
}

func (i *impl) Name() string {
	return trigger.AppName
}

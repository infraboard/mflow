package impl

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/infraboard/mcenter/clients/rpc"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/apps/trigger"
)

func init() {
	ioc.RegistryController(&impl{})
}

type impl struct {
	col *mongo.Collection
	log *zerolog.Logger
	trigger.UnimplementedRPCServer
	ioc.ObjectImpl

	// 构建配置
	build build.Service
	// 执行流水线
	task task.PipelineService
	// mcenter客户端
	mcenter *rpc.ClientSet
}

func (i *impl) Init() error {
	i.col = ioc_mongo.DB().Collection(i.Name())
	i.log = logger.Sub(i.Name())

	i.build = ioc.GetController(build.AppName).(build.Service)
	i.task = ioc.GetController(task.AppName).(task.Service)
	i.mcenter = rpc.C()
	return nil
}

func (i *impl) Name() string {
	return trigger.AppName
}

func (i *impl) Registry(server *grpc.Server) {
	trigger.RegisterRPCServer(server, i)
}

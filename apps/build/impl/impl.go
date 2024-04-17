package impl

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/pipeline"
)

func init() {
	ioc.Controller().Registry(&impl{})
}

type impl struct {
	build.UnimplementedRPCServer
	ioc.ObjectImpl

	mcenter  *rpc.ClientSet
	col      *mongo.Collection
	log      *zerolog.Logger
	pipeline pipeline.Service
}

func (i *impl) Init() error {
	i.col = ioc_mongo.DB().Collection(i.Name())
	i.log = log.Sub(i.Name())
	i.mcenter = rpc.C()
	i.pipeline = ioc.Controller().Get(pipeline.AppName).(pipeline.Service)

	build.RegisterRPCServer(grpc.Get().Server(), i)
	return nil
}

func (i *impl) Name() string {
	return build.AppName
}

package impl

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/build"
)

func init() {
	ioc.RegistryController(&impl{})
}

type impl struct {
	col *mongo.Collection
	log *zerolog.Logger
	build.UnimplementedRPCServer
	ioc.ObjectImpl

	mcenter *rpc.ClientSet
}

func (i *impl) Init() error {
	i.col = ioc_mongo.DB().Collection(i.Name())
	i.log = logger.Sub(i.Name())
	i.mcenter = rpc.C()
	return nil
}

func (i *impl) Name() string {
	return build.AppName
}

func (i *impl) Registry(server *grpc.Server) {
	build.RegisterRPCServer(server, i)
}

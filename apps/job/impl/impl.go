package impl

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/conf"
)

func init() {
	ioc.RegistryController(&impl{})
}

type impl struct {
	col *mongo.Collection
	log *zerolog.Logger
	job.UnimplementedRPCServer
	ioc.ObjectImpl
}

func (i *impl) Init() error {
	db, err := conf.C().Mongo.GetDB()
	if err != nil {
		return err
	}
	i.col = db.Collection(i.Name())
	i.log = logger.Sub(i.Name())
	return nil
}

func (i *impl) Name() string {
	return job.AppName
}

func (i *impl) Registry(server *grpc.Server) {
	job.RegisterRPCServer(server, i)
}

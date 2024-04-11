package impl

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/pipeline"
)

func init() {
	ioc.Controller().Registry(&impl{})
}

type impl struct {
	col *mongo.Collection
	log *zerolog.Logger
	pipeline.UnimplementedRPCServer
	ioc.ObjectImpl

	job job.Service
}

func (i *impl) Init() error {
	i.col = ioc_mongo.DB().Collection(i.Name())
	i.log = log.Sub(i.Name())
	i.job = ioc.Controller().Get(job.AppName).(job.Service)

	pipeline.RegisterRPCServer(grpc.Get().Server(), i)
	return nil
}

func (i *impl) Name() string {
	return pipeline.AppName
}

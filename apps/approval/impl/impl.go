package impl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/approval"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
)

func init() {
	ioc.Controller().Registry(&impl{})
}

type impl struct {
	col *mongo.Collection
	log *zerolog.Logger
	approval.UnimplementedRPCServer
	ioc.ObjectImpl

	pipeline pipeline.Service
	task     task.PipelineService
	mcenter  *rpc.ClientSet
}

func (s *impl) Init() error {
	s.col = ioc_mongo.DB().Collection(s.Name())
	indexs := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "create_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "domain", Value: -1},
				{Key: "namespace", Value: -1},
				{Key: "version", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := s.col.Indexes().CreateMany(context.Background(), indexs)
	if err != nil {
		return err
	}

	s.log = log.Sub(s.Name())
	s.pipeline = ioc.Controller().Get(pipeline.AppName).(pipeline.Service)
	s.task = ioc.Controller().Get(task.AppName).(task.Service)
	s.mcenter = rpc.C()

	approval.RegisterRPCServer(grpc.Get().Server(), s)
	return nil
}

func (s *impl) Name() string {
	return approval.AppName
}

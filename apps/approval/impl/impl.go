package impl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	ioc_mongo "github.com/infraboard/mcube/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/approval"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
)

func init() {
	ioc.RegistryController(&impl{})
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

	s.log = logger.Sub(s.Name())
	s.pipeline = ioc.GetController(pipeline.AppName).(pipeline.Service)
	s.task = ioc.GetController(task.AppName).(task.Service)
	s.mcenter = rpc.C()
	return nil
}

func (s *impl) Name() string {
	return approval.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	approval.RegisterRPCServer(server, s)
}

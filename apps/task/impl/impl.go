package impl

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/approval"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mflow/apps/trigger"

	// 加载并初始化Runner
	"github.com/infraboard/mflow/apps/task/runner"
	_ "github.com/infraboard/mflow/apps/task/runner/k8s"
	"github.com/infraboard/mflow/apps/task/webhook"
)

func init() {
	ioc.Controller().Registry(&impl{})
}

type impl struct {
	jcol *mongo.Collection
	pcol *mongo.Collection
	log  *zerolog.Logger
	task.UnimplementedJobRPCServer
	task.UnimplementedPipelineRPCServer
	ioc.ObjectImpl

	job      job.Service
	pipeline pipeline.Service
	approval approval.Service
	hook     *webhook.WebHook
	trigger  trigger.Service

	mcenter *rpc.ClientSet
}

func (i *impl) Init() error {
	db := ioc_mongo.DB()
	i.jcol = db.Collection("job_tasks")
	i.pcol = db.Collection("pipeline_tasks")
	i.log = log.Sub(i.Name())
	i.job = ioc.Controller().Get(job.AppName).(job.Service)
	i.pipeline = ioc.Controller().Get(pipeline.AppName).(pipeline.Service)
	i.approval = ioc.Controller().Get(approval.AppName).(approval.Service)
	i.trigger = ioc.Controller().Get(trigger.AppName).(trigger.Service)
	i.mcenter = rpc.C()
	if err := runner.Init(); err != nil {
		return err
	}

	i.hook = webhook.NewWebHook()
	return nil
}

func (i *impl) Name() string {
	return task.AppName
}

func (i *impl) Registry(server *grpc.Server) {
	task.RegisterJobRPCServer(server, i)
	task.RegisterPipelineRPCServer(server, i)
}

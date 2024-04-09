package impl

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mpaas/apps/k8s"
	mpaas "github.com/infraboard/mpaas/clients/rpc"
)

func init() {
	ioc.Controller().Registry(&impl{})
}

type impl struct {
	job.UnimplementedRPCServer
	ioc.ObjectImpl

	col   *mongo.Collection
	log   *zerolog.Logger
	mpaas *mpaas.ClientSet
}

func (i *impl) Init() error {
	i.col = ioc_mongo.DB().Collection(i.Name())
	i.log = log.Sub(i.Name())
	i.mpaas = mpaas.C()

	job.RegisterRPCServer(grpc.Get().Server(), i)
	return nil
}

func (i *impl) Name() string {
	return job.AppName
}

func (i *impl) InjectK8sCluster(
	ctx context.Context,
	domain, namespace string,
	jobs ...*job.Job) error {
	req := k8s.NewQueryClusterRequest()
	req.Scope.Domain = domain
	req.Scope.Namespace = namespace
	cs, err := i.mpaas.K8s().QueryCluster(ctx, req)
	if err != nil {
		return err
	}

	for i := range jobs {
		j := jobs[i]
		if j.Spec.RunnerType != job.RUNNER_TYPE_K8S_JOB {
			continue
		}

		kf := j.Spec.RunParams.GetParam("_kube_config_from")
		if kf != nil && kf.Value !=
			job.KUBE_CONF_FROM_MPAAS_K8S_CLUSTER_REF.String() {
			continue
		}

		parm := j.Spec.RunParams.GetParam("_kube_config")
		if parm != nil {
			parm.EnumOptions = []*job.EnumOption{}
			for _, c := range cs.Items {
				parm.EnumOptions = append(parm.EnumOptions,
					&job.EnumOption{
						Label:      c.Spec.Name,
						Value:      c.Meta.Id,
						Extensions: map[string]string{},
					})
			}
		}
	}

	return nil
}

package impl_test

import (
	"testing"

	"github.com/infraboard/mflow/apps/build"
	"github.com/infraboard/mflow/apps/job"
	"github.com/infraboard/mflow/test/conf"
	"github.com/infraboard/mflow/test/tools"
)

func TestQueryJob(t *testing.T) {
	req := job.NewQueryJobRequest()
	set, err := impl.QueryJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(set))
}

func TestCreateTestJob(t *testing.T) {
	req := job.NewCreateJobRequest()
	req.Name = "test"
	req.CreateBy = "test"
	req.RunnerSpec = tools.MustReadContentFile("test/test.yml")
	param := job.NewRunParamSet()
	param.Add(&job.RunParam{
		Required: true,
		Name:     "cluster_id",
		NameDesc: "job运行时的k8s集群",
		Value:    "k8s-test",
	})
	param.Add(&job.RunParam{
		Required: true,
		Name:     "namespace",
		NameDesc: "job运行时的namespace",
		Value:    "default",
	})
	req.RunParams = param

	ins, err := impl.CreateJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

func TestCreateBuildJob(t *testing.T) {
	req := job.NewCreateJobRequest()
	req.Name = "docker_build"
	req.DisplayName = "容器构建"
	req.CreateBy = "admin@default"
	req.RunnerSpec = tools.MustReadContentFile("test/container_build.yml")
	req.Description = "使用git拉取源代码, 然后使用kaniko完成应用容器镜像的构建与推送"

	ins, err := impl.CreateJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

func TestCreateDeployJob(t *testing.T) {
	req := job.NewCreateJobRequest()
	req.Name = "docker_deploy"
	req.DisplayName = "容器部署"
	req.CreateBy = "admin@default"
	req.RunnerSpec = tools.MustReadContentFile("test/container_deploy.yml")
	req.Description = "通过kubectl 的set image命令来更新部署镜像的版本"

	ins, err := impl.CreateJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

func TestUpdateDeployJob(t *testing.T) {
	req := job.NewPatchJobRequest(conf.C.DEPLOY_JOB_ID)
	req.Spec.RunnerSpec = tools.MustReadContentFile("test/container_deploy.yml")
	req.Spec.Labels["Language"] = "*"
	req.Spec.Icon = tools.MustReadContentFile("test/container_deploy.svg")
	param := job.NewRunParamSet()
	param.Add(&job.RunParam{
		Required:    true,
		Name:        "_kube_config",
		NameDesc:    "用于运行k8s job的访问配置",
		Value:       tools.MustReadContentFile("test/kube_config.yml"),
		IsSensitive: true,
		SearchLabel: false,
		Extensions: map[string]string{
			"format": "yaml",
		},
	})
	param.Add(&job.RunParam{
		Required:    true,
		Name:        "_namespace",
		NameDesc:    "k8s job运行时的namespace",
		Value:       "default",
		SearchLabel: true,
	})

	// 部署运行时变量
	param.Add(&job.RunParam{
		Required: true,
		Name:     "service_account",
		NameDesc: `部署时关联的账号,  默认default并没有修改pod权限, 创建方式: 
		kubectl create serviceaccount mflow-deploy &&
		kubectl create role deployments-manager --verb=get,update,list,watch --resource=deployments &&
		kubectl create rolebinding deployments-manager-binding --role=deployments-manager --serviceaccount=default:mflow-deploy`,
		Example:     "mflow-deploy",
		SearchLabel: true,
	})
	param.Add(&job.RunParam{
		Required:    false,
		Name:        build.SYSTEM_VARIABLE_IMAGE_REPOSITORY,
		NameDesc:    "应用部署的镜像仓库地址, 默认沿用原有的镜像仓库地址",
		Example:     "registry.cn-hangzhou.aliyuncs.com/infraboard/mcenter",
		ValueDesc:   "如果 前面有构建操作, 该变量由构建操作提供",
		SearchLabel: true,
	})
	param.Add(&job.RunParam{
		Required:    true,
		Name:        build.SYSTEM_VARIABLE_APP_VERSION,
		NameDesc:    "应用部署时的版本",
		ValueDesc:   "如果 前面有构建操作, 该变量由构建操作提供",
		Example:     "v0.0.1",
		SearchLabel: true,
	})
	req.Spec.RunParams = param

	ins, err := impl.UpdateJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

func TestUpdateBuildJob(t *testing.T) {
	req := job.NewPatchJobRequest(conf.C.BUILD_JOB_ID)
	req.Spec.RunnerSpec = tools.MustReadContentFile("test/container_build.yml")
	req.Spec.Labels["Language"] = "*"
	req.Spec.Icon = tools.MustReadContentFile("test/container_build.svg")
	param := job.NewRunParamSet()
	param.Add(&job.RunParam{
		Required:    true,
		Name:        "_kube_config",
		NameDesc:    "用于运行k8s job的访问配置, 默认位于: ~/.kube/config",
		Value:       tools.MustReadContentFile("test/kube_config.yml"),
		IsSensitive: true,
		SearchLabel: false,
	})
	param.Add(&job.RunParam{
		Required:    true,
		Name:        "_namespace",
		NameDesc:    "job运行时的namespace",
		Value:       "default",
		SearchLabel: true,
	})

	// 需要构建的代码信息
	param.Add(&job.RunParam{
		Required:    true,
		Name:        "GIT_SSH_URL",
		NameDesc:    "应用git代码仓库地址",
		ValueDesc:   "CI时, 系统根据Webhook自动注入",
		Example:     "git@github.com:infraboard/mpaas.git",
		SearchLabel: true,
	})
	param.Add(&job.RunParam{
		Required:  true,
		Name:      "GIT_BRANCH",
		NameDesc:  "需要拉去的代码分支",
		ValueDesc: "CI时, 系统根据Webhook自动注入",
		Example:   "master",
	})
	param.Add(&job.RunParam{
		Required:  true,
		Name:      "GIT_COMMIT_ID",
		NameDesc:  "应用git代码仓库地址",
		ValueDesc: "CI时, 系统根据Webhook自动注入",
		Example:   "32d63566098f7e0b0ac3a3d8ddffe71cc6cad7b0",
	})

	// 构建缓存
	param.Add(&job.RunParam{
		Required:  false,
		Name:      "CACHE_ENABLE",
		NameDesc:  "是否启动构建缓存, 构建缓存能有效加快构建的速度, 避免每次构建都通过远程网络下载构建依赖",
		Example:   "true",
		ValueType: job.PARAM_VALUE_TYPE_BOOLEAN,
		Value:     "true",
	})
	param.Add(&job.RunParam{
		Required: false,
		Name:     "CACHE_REPO",
		NameDesc: "构建缓存的镜像仓库地址, 默认为: [镜像推送地址/cache], 需要使用独立的缓存仓库时设置",
		Example:  "registry.cn-hangzhou.aliyuncs.com/build_cache/mpaas",
	})
	param.Add(&job.RunParam{
		Required:  false,
		Name:      "CACHE_COMPRESS",
		NameDesc:  "镜像缓存层压缩, 这可以降低缓存镜像的大小, 但是压缩时需要花费额外的内存, 对于大型构建特别有用, 默认为true, 但是如果出现内存不足错误, 请关闭",
		Example:   "true",
		ValueType: job.PARAM_VALUE_TYPE_BOOLEAN,
		Value:     "true",
	})
	param.Add(&job.RunParam{
		Required:  false,
		Name:      "CUSTOM_PLATFORM",
		NameDesc:  "目标镜像的平台架构信息, 详情参考: https://github.com/GoogleContainerTools/kaniko#flag---custom-platform",
		Example:   "linux/amd64",
		ValueType: job.PARAM_VALUE_TYPE_ENUM,
		EnumOptions: job.NewEnumOptionWithKVPaire(
			"linux/386", "X86",
			"linux/amd64", "AMD64",
			"linux/arm", "ARM",
			"linux/arm64", "ARM64",
		),
		Value: "linux/amd64",
	})

	// docker push registry.cn-hangzhou.aliyuncs.com/infraboard/mpaas:[镜像版本号]
	param.Add(&job.RunParam{
		Required:  true,
		Name:      build.SYSTEM_VARIABLE_IMAGE_REPOSITORY,
		NameDesc:  "镜像推送地址",
		ValueDesc: "CI时, 通过读取构建配置自动获取",
		Example:   "registry.cn-hangzhou.aliyuncs.com/infraboard/mpaas",
	})
	param.Add(&job.RunParam{
		Required:  true,
		Name:      build.SYSTEM_VARIABLE_APP_VERSION,
		NameDesc:  "镜像版本",
		ValueDesc: "CI时, 通过构建配置自动生成",
		Example:   "v0.0.2",
	})
	param.Add(&job.RunParam{
		Required:  false,
		Name:      "APP_DOCKERFILE",
		NameDesc:  "应用git代码仓库中用于构建镜像的Dockerfile路径",
		ValueDesc: "CI时, 通过读取构建配置自动获取",
		Value:     "Dockerfile",
		Example:   "Dockerfile",
	})
	param.Add(&job.RunParam{
		UsageType: job.PARAM_USAGE_TYPE_TEMPLATE,
		Required:  true,
		Name:      "git_ssh_secret",
		NameDesc:  "用于拉取git仓库代码的secret名称, kubectl create secret generic git-ssh-key --from-file=id_rsa=${HOME}/.ssh/id_rsa",
		ValueType: job.PARAM_VALUE_TYPE_K8S_SECRET,
		ValueDesc: "如果使用默认值, 请提前确认名为git-ssh-key的secret已经创建，namespace请参考前面参数",
		Example:   "git-ssh-key",
		Value:     "git-ssh-key",
	})
	param.Add(&job.RunParam{
		Required:  true,
		UsageType: job.PARAM_USAGE_TYPE_TEMPLATE,
		Name:      "image_push_secret",
		NameDesc:  "用于推送镜像的secret名称, kubectl create secret generic kaniko-secret --from-file=apps/job/impl/test/config.json 具体文档参考: https://github.com/GoogleContainerTools/kaniko#pushing-to-docker-hub",
		ValueType: job.PARAM_VALUE_TYPE_K8S_SECRET,
		ValueDesc: "如果使用默认值, 请提前确认名为kaniko-secret的secret已经创建, namespace请参考前面参数, 注意挂载文件名字config.json",
		Example:   "kaniko-secret",
		Value:     "kaniko-secret",
	})
	req.Spec.RunParams = param

	ins, err := impl.UpdateJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

func TestDeleteJob(t *testing.T) {
	req := job.NewDeleteJobRequest(conf.C.BUILD_JOB_ID)
	ins, err := impl.DeleteJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

func TestDescribeJob(t *testing.T) {
	req := job.NewDescribeJobRequestByName("#d98eb61c5d3cc142")
	ins, err := impl.DescribeJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

// 发布Job
func TestUpdateBuildJobStatus(t *testing.T) {
	req := job.NewUpdateJobStatusRequest("docker_build@default.default")
	req.Status.Stage = job.JOB_STAGE_PUBLISHED
	req.Status.Version = "v1"
	ins, err := impl.UpdateJobStatus(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

func TestUpdateDeployJobStatus(t *testing.T) {
	req := job.NewUpdateJobStatusRequest("docker_deploy@default.default")
	req.Status.Stage = job.JOB_STAGE_PUBLISHED
	req.Status.Version = "v1"
	ins, err := impl.UpdateJobStatus(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tools.MustToJson(ins))
}

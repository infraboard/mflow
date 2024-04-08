package job

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"dario.cat/mergo"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/pb/resource"
	"github.com/infraboard/mcube/v2/tools/sense"
	"github.com/infraboard/mflow/common/hash"
	k8sApp "github.com/infraboard/mpaas/apps/k8s"
	mpaas "github.com/infraboard/mpaas/clients/rpc"
	"github.com/infraboard/mpaas/provider/k8s"
	"github.com/infraboard/mpaas/provider/k8s/workload"
	v1 "k8s.io/api/core/v1"
)

// New 新建一个部署配置
func New(req *CreateJobRequest) (*Job, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	req.BuildSearchLabels()
	job := &Job{
		Meta:   resource.NewMeta(),
		Spec:   req,
		Status: NewJobStatus(),
	}

	job.Meta.Id = hash.FnvHash(job.UniqName())
	return job, nil
}

func (r *CreateJobRequest) BuildSearchLabels() {
	if r.Labels == nil {
		r.Labels = map[string]string{}
	}

	sl := r.RunParams.SearchLabels()
	for k, v := range sl {
		r.Labels[k] = v
	}
}

func NewJobSet() *JobSet {
	return &JobSet{
		Items: []*Job{},
	}
}

func (s *JobSet) Add(item *Job) {
	s.Items = append(s.Items, item)
}

func (s *JobSet) Desense() {
	for i := range s.Items {
		item := s.Items[i]
		item.Desense()
	}
}

func NewDefaultJob() *Job {
	return &Job{
		Spec: NewCreateJobRequest(),
	}
}

func (i *Job) Update(req *UpdateJobRequest) {
	i.Meta.UpdateAt = time.Now().Unix()
	i.Meta.UpdateBy = req.UpdateBy
	i.Spec = req.Spec
}

func (i *Job) Patch(req *UpdateJobRequest) error {
	i.Meta.UpdateAt = time.Now().Unix()
	i.Meta.UpdateBy = req.UpdateBy
	return mergo.MergeWithOverwrite(i.Spec, req.Spec)
}

func (i *Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*resource.Meta
		*CreateJobRequest
		Status *JobStatus `json:"status"`
	}{i.Meta, i.Spec, i.Status})
}

// 只可以删除未发布的Job
func (i *Job) CheckAllowDelete() error {
	if !i.Status.Stage.Equal(JOB_STAGE_DRAFT) {
		return fmt.Errorf("只有处于草稿状态的Job才允许删除")
	}
	return nil
}

func (i *Job) Desense() {
	if i.Spec == nil {
		return
	}

	if i.Spec.RunParams != nil {
		i.Spec.RunParams.Densense()
	}

	if i.Spec.RollbackParam != nil {
		i.Spec.RollbackParam.Densense()
	}
}

// 补充扩展属性
func (i *Job) AddExtension() {
	if i.Spec == nil {
		return
	}

	// 填充选项
	for idx := range i.Spec.RunParams.Params {
		param := i.Spec.RunParams.Params[idx]
		switch param.ValueType {
		case PARAM_VALUE_TYPE_BOOLEAN:
			param.AddEnumOptions(&EnumOption{
				Value: "true",
				Label: "是",
			}, &EnumOption{
				Value: "false",
				Label: "否",
			})
		}
	}
	i.Spec.Extension["uniq_name"] = i.UniqName()
}

// docker_build@default.default:v1
func (i *Job) UniqName() string {
	return fmt.Sprintf("%s@%s.%s:%s",
		i.Spec.Name,
		i.Spec.Namespace,
		i.Spec.Domain, i.GetVersion())
}

func (i *Job) GetVersion() string {
	if i.Status == nil {
		return ""
	}

	return i.Status.Version
}

func NewRunParamSet() *RunParamSet {
	return &RunParamSet{
		Params: []*RunParam{},
	}
}

// 参数脱敏, 注意 不能在运行过程中脱敏, 仅仅在需要显示时，调用该方法
func (r *RunParamSet) Densense() {
	for i := range r.Params {
		r.Params[i].Desense()
	}
}

// 绕开Merge, 直接注入, 因为Merge只允许注入Job声明的变量
// 非job声明的变量只能通过Add添加, 比如系统变量
// 注意系统变量是系统注入, 用户无法注入系统变量
func (r *RunParamSet) Add(items ...*RunParam) {
	for i := range items {
		item := items[i]
		if !item.IsSystemVariable() {
			item.Init()
			r.Params = append(r.Params, item)
		}
	}
}

// 检查是否有重复的参数
func (r *RunParamSet) CheckDuplicate() error {
	kc := map[string]int{}
	for i := range r.Params {
		p := r.Params[i]
		kc[p.Name]++
	}

	duplicates := []string{}
	for k, v := range kc {
		if v > 1 {
			duplicates = append(duplicates, fmt.Sprintf("%s duplicate count %d", k, v))
		}
	}

	if len(duplicates) > 0 {
		return fmt.Errorf("params %s duplicate", duplicates)
	}
	return nil
}

func (r *RunParamSet) Validate() error {
	err := r.CheckDuplicate()
	if err != nil {
		return err
	}

	for i := range r.Params {
		p := r.Params[i]
		if p.Required && p.Value == "" {
			return fmt.Errorf("参数: %s 不能为空", p.Name)
		}
	}

	return nil
}

// 从参数中提取k8s job执行器(runner)需要的参数
// 这里采用反射来获取Struc Tag, 然后根据Struct Tag 获取参数的具体值
// 关于反射 可以参考: https://blog.csdn.net/bocai_xiaodaidai/article/details/123668047
func (r *RunParamSet) K8SJobRunnerParams() *K8SJobRunnerParams {
	params := NewK8SJobRunnerParams()

	// params是一个Pointer Value, 如果需要获取值的类型需要这样处理:
	//	reflect.Indirect(reflect.ValueOf(params)).Type()
	// 因此这里直接采用K8SJobRunnerParams{}获取类型
	pt := reflect.TypeOf(K8SJobRunnerParams{})

	// go语言所有函数传的都是值，所以要想修改原来的值就需要传指
	// 通过Elem()返回指针指向的对象
	v := reflect.ValueOf(params).Elem()

	for i := 0; i < pt.NumField(); i++ {
		field := pt.Field(i)
		if field.IsExported() {
			tagValue := field.Tag.Get("param")
			param := r.GetParam(tagValue)
			if param != nil {
				switch param.ValueType {
				case PARAM_VALUE_TYPE_BOOLEAN:
					boolV, _ := strconv.ParseBool(param.Value)
					v.Field(i).SetBool(boolV)
				case PARAM_VALUE_TYPE_ENUM:
					from, err := ParseKUBE_CONF_FROMFromString(param.Value)
					if err != nil {
						log.L().Error().Msgf("parse enum error: %s", err)
					} else {
						v.Field(i).SetInt(int64(from))
					}
				default:
					v.Field(i).SetString(param.Value)
				}

			}

		}
	}

	return params
}

func (r *RunParamSet) GetDeploymentId() string {
	return r.GetParamValue(SYSTEM_VARIABLE_DEPLOY_ID)
}

func (r *RunParamSet) GetJobTaskId() string {
	return r.GetParamValue(SYSTEM_VARIABLE_JOB_TASK_ID)
}

func (r *RunParamSet) GetJobId() string {
	return r.GetParamValue(SYSTEM_VARIABLE_JOB_ID)
}

func (r *RunParamSet) GetPipelineTaskId() string {
	return r.GetParamValue(SYSTEM_VARIABLE_PIPELINE_TASK_ID)
}

// 获取需要注入容器的环境变量参数:
//
//	用户变量: 大写开头的变量, 因为一般环境变量都是大写的比如 DB_PASS,
//	系统变量: _开头为系统变量, 由Runner处理并注入, 比如 _DEPLOY_ID
//	Runner变量: 小写的变量, 用于系统内部使用, 不会注入, 比如 K8SJobRunnerParams 中的cluster_id
func (r *RunParamSet) EnvVars() (envs []v1.EnvVar) {
	for i := range r.Params {
		item := r.Params[i]
		// 只导出环境变量
		if !item.UsageType.Equal(PARAM_USAGE_TYPE_ENV) {
			continue
		}
		if item.Name != "" && (unicode.IsUpper(rune(item.Name[0])) || strings.HasPrefix(item.Name, "_")) {
			envs = append(envs, v1.EnvVar{
				Name:  item.Name,
				Value: item.Value,
			})
		}
	}
	return
}

func (r *RunParamSet) TemplateVars() (vars []*RunParam) {
	for i := range r.Params {
		item := r.Params[i]
		// 只导出模版变量
		if item.UsageType.Equal(PARAM_USAGE_TYPE_TEMPLATE) {
			vars = append(vars, item)
		}
	}
	return
}

func ParamsToEnvVar(params []*RunParam) (envs []v1.EnvVar) {
	for i := range params {
		item := params[i]
		envs = append(envs, v1.EnvVar{
			Name:  item.Name,
			Value: item.Value,
		})
	}
	return
}

// 获取参数的值
func (r *RunParamSet) GetParamValue(key string) string {
	for i := range r.Params {
		item := r.Params[i]
		if item.Name == key {
			return item.Value
		}
	}
	return ""
}

// 获取参数的值
func (r *RunParamSet) GetParam(key string) *RunParam {
	for i := range r.Params {
		item := r.Params[i]
		if item.Name == key {
			return item
		}
	}
	return nil
}

// 设置参数的值, 注意如果value为""则不会修改
func (r *RunParamSet) SetParamValue(key, value string, readOnly bool) {
	for i := range r.Params {
		param := r.Params[i]
		if param.IsEdit() && param.Name == key {
			param.Value = value
			param.ReadOnly = readOnly
			return
		}
	}
	log.L().Warn().Msgf("set param %s value failed, job no param or readonly", key)
}

func (r *RunParamSet) Merge(targets ...*RunParam) {
	for i := range targets {
		t := targets[i]
		t.TrimSpace()
		if t.Value != "" {
			r.SetParamValue(t.Name, t.Value, t.ReadOnly)
		}
	}
}

// 脱敏后的数据还原
func (r *RunParamSet) RestoreSensitive(set *RunParamSet) {
	for i := range set.Params {
		param := set.Params[i]
		if param.IsSensitive && r.GetParamValue(param.Name) != param.Value {
			r.Merge(param)
		}
	}
}

func (r *RunParamSet) SearchLabels() map[string]string {
	labels := map[string]string{}
	for i := range r.Params {
		p := r.Params[i]
		if p.SearchLabel && p.Value != "" {
			labels[p.Name] = p.Value
		}
	}

	return labels
}

func NewK8SJobRunnerParams() *K8SJobRunnerParams {
	return &K8SJobRunnerParams{}
}

func (p *K8SJobRunnerParams) KubeConfSecret(name string, mountPath string) *v1.Secret {
	secret := new(v1.Secret)
	secret.Name = name
	secret.StringData = map[string]string{
		"config": p.KubeConfig,
	}
	secret.Annotations = map[string]string{
		workload.ANNOTATION_SECRET_MOUNT: mountPath,
	}
	return secret
}

func (p *K8SJobRunnerParams) Client(ctx context.Context) (*k8s.Client, error) {
	// 如果集群配置托管在mpaas k8s cluster
	switch p.KubeConfigFrom {
	case KUBE_CONF_FROM_MPAAS_K8S_CLUSTER_REF:
		descReq := k8sApp.NewDescribeClusterRequest(p.KubeConfig)
		k8sCluster, err := mpaas.C().K8s().DescribeCluster(ctx, descReq)
		if err != nil {
			return nil, err
		}
		p.KubeConfig = k8sCluster.Spec.KubeConfig
		log.L().Debug().Msgf("load kube config from mpaas k8s cluster: %s", descReq.Id)
	}

	if p.KubeConfig == "" {
		return nil, fmt.Errorf("kube config not config")
	}
	return k8s.NewClient(p.KubeConfig)
}

func NewRunParam(name, value string) *RunParam {
	return &RunParam{
		Name:        name,
		Value:       value,
		EnumOptions: []*EnumOption{},
		ParamScope:  NewParamScope(),
		Extensions:  map[string]string{},
	}
}

func NewEnumOption(label, value string) *EnumOption {
	return &EnumOption{
		Label:      label,
		Value:      value,
		Extensions: map[string]string{},
	}
}

// label, value
func NewEnumOptionWithKVPaire(kvs ...string) (options []*EnumOption) {
	m := NewMapWithKVPaire(kvs...)
	for k, v := range m {
		options = append(options, NewEnumOption(v, k))
	}
	return nil
}

func NewParamScope() *ParamScope {
	return &ParamScope{
		Label: map[string]string{},
	}
}

func NewRunParamWithKVPaire(kvs ...string) (params []*RunParam) {
	m := NewMapWithKVPaire(kvs...)
	for k, v := range m {
		params = append(params, NewRunParam(k, v))
	}

	return
}

// 引用名称
func (p *RunParam) RefName() string {
	return fmt.Sprintf("${%s}", p.Name)
}

// 添加选项
func (p *RunParam) AddEnumOptions(options ...*EnumOption) {
	// 初始化
	if p.EnumOptions == nil {
		p.EnumOptions = []*EnumOption{}
	}

	// 如果不存在再注入
	for _, option := range options {
		v := p.GetEnumOptionByValue(option.Value)
		if v == nil {
			p.EnumOptions = append(p.EnumOptions, option)
		}
	}
}

// 根据value获取选项
func (p *RunParam) GetEnumOptionByValue(v string) *EnumOption {
	for _, option := range p.EnumOptions {
		if option.Value == v {
			return option
		}
	}
	return nil
}

// 是否允许修改
func (p *RunParam) IsEdit() bool {
	// 只读且有值时不允许修改
	if p.ReadOnly && p.Value != "" {
		return false
	}
	return true
}

// 设置ReadOnly
func (p *RunParam) SetReadOnly(v bool) *RunParam {
	p.ReadOnly = v
	return p
}

// 值脱敏
func (p *RunParam) Desense() *RunParam {
	if p.IsSensitive {
		p.Value = sense.DeSense(p.Value)
	}
	return p
}

// 剔除值里面的空白字符
func (p *RunParam) TrimSpace() {
	p.Value = strings.TrimSpace(p.Value)
}

// 值脱敏
func (p *RunParam) Init() {
	if p.Extensions == nil {
		p.Extensions = map[string]string{}
	}
}

// 是否是系统变量
func (p *RunParam) IsSystemVariable() bool {
	return p.UsageType == PARAM_USAGE_TYPE_SYSTEM
}

// 设置SearchLabel
func (p *RunParam) SetSearchLabel(v bool) *RunParam {
	p.SearchLabel = v
	return p
}

// 设置Required
func (p *RunParam) SetRequired(v bool) *RunParam {
	p.Required = v
	return p
}

// Markdown格式的简要说明
func (p *RunParam) MarkdownShortShow() string {
	return fmt.Sprintf("%s: %s", p.Name, p.Value)
}

func NewJobStatus() *JobStatus {
	return &JobStatus{}
}

func ParseRunParamFromBytes(content []byte) ([]*RunParam, error) {
	envs := []*RunParam{}
	lines := []string{}
	line := []byte{}
	for _, c := range content {
		if c == '\n' {
			lines = append(lines, string(line))
			line = []byte{}
		} else {
			line = append(line, c)
		}
	}

	for _, l := range lines {
		l := strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		kvs := strings.Split(l, "=")
		if len(kvs) != 2 {
			return nil, fmt.Errorf("环境变量格式错误: %s", kvs)
		}
		k, v := kvs[0], kvs[1]

		env := NewRunParam(k, strings.Trim(v, `"`))
		envs = append(envs, env)
	}

	return envs, nil
}

func (r *RunParam) FileLine() (line []byte) {
	return []byte(fmt.Sprintf("%s=%s\n", r.Name, r.Value))
}

func (r *RunParam) IsExport() bool {
	if r.Name == "" && unicode.IsUpper(rune(r.Name[0])) {
		return true
	}
	return false
}

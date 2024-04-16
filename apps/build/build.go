package build

import (
	"encoding/json"
	"path/filepath"
	"regexp"

	"github.com/infraboard/mcube/v2/pb/resource"
	"github.com/infraboard/mflow/apps/job"
)

// New 新建一个domain
func New(req *CreateBuildConfigRequest) (*BuildConfig, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	d := &BuildConfig{
		Meta: resource.NewMeta(),
		Spec: req,
	}

	return d, nil
}

func NewBuildConfigSet() *BuildConfigSet {
	return &BuildConfigSet{
		Items: []*BuildConfig{},
	}
}

func (s *BuildConfigSet) Add(item *BuildConfig) {
	s.Items = append(s.Items, item)
}

func (s *BuildConfigSet) Len() int {
	return len(s.Items)
}

// 比如: "foo.*"
func (s *BuildConfigSet) MatchSubEvent(branchRegExp string) *BuildConfigSet {
	set := NewBuildConfigSet()

	for i := range s.Items {
		item := s.Items[i]
		if item.Spec.Condition.MatchSubEvent(branchRegExp) {
			set.Add(item)
		}
	}

	return set
}

func NewDefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		Spec: NewCreateBuildConfigRequest(),
	}
}

func (b *BuildConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*resource.Meta
		*CreateBuildConfigRequest
	}{b.Meta, b.Spec})
}

func NewCreateBuildConfigRequest() *CreateBuildConfigRequest {
	return &CreateBuildConfigRequest{
		Enabled:       true,
		VersionPrefix: "v",
		Condition:     NewTrigger(),
		Labels:        make(map[string]string),
		CustomParams:  []*job.RunParam{},
		Extra:         make(map[string]string),
	}
}

func NewImageBuild() *ImageBuildConfig {
	return &ImageBuildConfig{
		DockerFile: "Dockerfile",
		Extra:      make(map[string]string),
	}
}

// 如果没有配置，则使用默认配置
func (c *ImageBuildConfig) GetDockerFileWithDefault(defaulDockerFile string) string {
	if c.DockerFile == "" {
		return defaulDockerFile
	}
	return c.DockerFile
}

// 如果没有配置，则使用默认配置
func (c *ImageBuildConfig) GetImageRepositoryWithDefault(defaultImageRepository string) string {
	if c.ImageRepository == "" {
		return defaultImageRepository
	}
	return c.ImageRepository
}

func NewPkgBuildConfig() *PkgBuildConfig {
	return &PkgBuildConfig{
		Extra: make(map[string]string),
	}
}

func NewTrigger() *Trigger {
	return &Trigger{
		Events:             []string{},
		SubEvents:          []string{},
		SubEventsMatchType: MATCH_TYPE_GLOB,
	}
}

func (t *Trigger) AddEvent(event string) {
	t.Events = append(t.Events, event)
}

func (t *Trigger) AddSubEvents(sub string) {
	t.SubEvents = append(t.SubEvents, sub)
}

// 关于Go语言正则表达式: http://c.biancheng.net/view/5124.html
func (t *Trigger) MatchSubEvent(pattern string) bool {
	if len(t.SubEvents) == 0 {
		return true
	}

	for _, b := range t.SubEvents {
		switch t.SubEventsMatchType {
		case MATCH_TYPE_GLOB:
			// 如果是* 直接匹配所有, 避免通配符/的问题
			if b == "*" {
				return true
			}
			ok, _ := filepath.Match(pattern, b)
			if ok {
				return true
			}
		case MATCH_TYPE_REGEXP:
			ok, _ := regexp.MatchString(pattern, b)
			if ok {
				return true
			}
		default:
			if pattern == b {
				return true
			}
		}
	}

	return false
}

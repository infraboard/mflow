package build

import (
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcenter/apps/policy"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc/config/validator"
	pb_request "github.com/infraboard/mcube/v2/pb/request"
)

const (
	AppName = "builds"
)

type Service interface {
	RPCServer
}

func (req *CreateBuildConfigRequest) Validate() error {
	if req.VersionPrefix == "" {
		req.VersionPrefix = "v"
	}
	return validator.Validate(req)
}

func NewQueryBuildConfigRequestFromHTTP(r *restful.Request) *QueryBuildConfigRequest {
	req := NewQueryBuildConfigRequest()
	req.Page = request.NewPageRequestFromHTTP(r.Request)
	req.Scope = token.GetTokenFromRequest(r).GenScope()
	req.Filters = policy.GetScopeFilterFromRequest(r)
	req.WithPipeline = r.QueryParameter("with_pipeline") == "true"

	sid := r.QueryParameter("service_id")
	if sid != "" {
		req.ServiceIds = strings.Split(sid, ",")
	}

	return req
}

func NewQueryBuildConfigRequest() *QueryBuildConfigRequest {
	return &QueryBuildConfigRequest{
		Page: request.NewDefaultPageRequest(),
	}
}

func (req *QueryBuildConfigRequest) AddService(serviceId string) {
	req.ServiceIds = append(req.ServiceIds, serviceId)
}

func (req *QueryBuildConfigRequest) SetEnabled(v bool) {
	req.Enabled = &v
}

func NewDescribeBuildConfigRequst(id string) *DescribeBuildConfigRequst {
	return &DescribeBuildConfigRequst{
		Id: id,
	}
}

func (req *DescribeBuildConfigRequst) Validate() error {
	return validator.Validate(req)
}

func NewDeleteBuildConfigRequest(id string) *DeleteBuildConfigRequest {
	return &DeleteBuildConfigRequest{
		Id: id,
	}
}

func NewPutBuildConfigRequest(id string) *UpdateBuildConfigRequest {
	return &UpdateBuildConfigRequest{
		Id:         id,
		UpdateMode: pb_request.UpdateMode_PUT,
		Spec:       NewCreateBuildConfigRequest(),
	}
}

func NewPatchBuildConfigRequest(id string) *UpdateBuildConfigRequest {
	return &UpdateBuildConfigRequest{
		Id:         id,
		UpdateMode: pb_request.UpdateMode_PATCH,
		Spec:       NewCreateBuildConfigRequest(),
	}
}

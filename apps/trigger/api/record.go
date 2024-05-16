package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mflow/apps/trigger"
)

func (h *Handler) QueryRecord(r *restful.Request, w *restful.Response) {
	req, err := trigger.NewQueryRecordRequestFromHTTP(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	set, err := h.svc.QueryRecord(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}

func (h *Handler) DescribeRecord(r *restful.Request, w *restful.Response) {
	req := trigger.NewDescribeRecordRequest(r.PathParameter("id"))
	ins, err := h.svc.DescribeRecord(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

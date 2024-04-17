package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mflow/apps/trigger"
)

func (h *Handler) QueryRecord(r *restful.Request, w *restful.Response) {
	req := trigger.NewQueryRecordRequestFromHTTP(r)

	set, err := h.svc.QueryRecord(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}

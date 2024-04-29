package api

import (
	"net/http"
	"time"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/websocket"
	"github.com/infraboard/mcenter/apps/endpoint"
	middleware "github.com/infraboard/mcenter/clients/rpc/middleware/auth/gorestful"
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mpaas/common/terminal"
	"github.com/rs/zerolog"

	"github.com/infraboard/mflow/apps/task"
)

func init() {
	ioc.Api().Registry(&WebsocketHandler{})
}

type WebsocketHandler struct {
	service task.Service
	log     *zerolog.Logger
	ioc.ObjectImpl
}

func (h *WebsocketHandler) Init() error {
	h.log = log.Sub(task.AppName)
	h.service = ioc.Controller().Get(task.AppName).(task.Service)
	h.Registry()
	return nil
}

func (h *WebsocketHandler) Name() string {
	return "ws/job_tasks"
}

func (h *WebsocketHandler) Version() string {
	return "v1"
}

func (h *WebsocketHandler) Registry() {
	h.RegistryUserHandler()
}

// 用户自己手动管理任务状态相关接口
func (h *WebsocketHandler) RegistryUserHandler() {
	tags := []string{"Job任务管理"}

	ws := gorestful.ObjectRouter(h)
	// Socket内鉴权
	ws.Route(ws.GET("/{id}/log").
		To(h.JobTaskLog).
		Doc("查询任务日志[WebSocket]").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, JOB_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Get.Value()).
		Reads(task.WatchJobTaskLogRequest{}).
		Writes(task.JobTaskStreamReponse{}))

	// Socket内鉴权
	ws.Route(ws.GET("/{id}/debug").
		To(h.DebugJobTask).
		Doc("任务在线Debug[WebSocket]").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, JOB_TASK_RESOURCE_NAME).
		Metadata(label.Action, label.Get.Value()).
		Reads(task.DebugJobTaskRequest{}).
		Writes(task.JobTaskStreamReponse{}))
}

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 60 * time.Second,
		ReadBufferSize:   8192,
		WriteBufferSize:  8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (h *WebsocketHandler) JobTaskLog(r *restful.Request, w *restful.Response) {
	ws, err := upgrader.Upgrade(w, r.Request, nil)
	if err != nil {
		response.Failed(w, err)
		return
	}

	term := task.NewTaskLogWebsocketTerminal(ws)

	// 开启认证与鉴权
	entry := endpoint.NewEntryFromRestRequest(r)
	if entry != nil {
		entry.SetAuthEnable(true)
		entry.SetPermissionEnable(true)
		err = middleware.Get().PermissionCheck(r, w, entry)
		if err != nil {
			term.Failed(err)
			return
		}
	}

	// 读取请求
	in := task.NewWatchJobTaskLogRequest(r.PathParameter("id"))
	if err = term.ReadReq(in); err != nil {
		term.Failed(err)
		return
	}

	// 输出日志到Term中
	if err = h.service.WatchJobTaskLog(in, term); err != nil {
		term.Failed(err)
		return
	}

	term.Success("ok")
}

func (h *WebsocketHandler) DebugJobTask(r *restful.Request, w *restful.Response) {
	ws, err := upgrader.Upgrade(w, r.Request, nil)
	if err != nil {
		response.Failed(w, err)
		return
	}

	term := terminal.NewWebSocketTerminal(ws)

	// 开启认证与鉴权
	entry := endpoint.NewEntryFromRestRequest(r).
		SetAuthEnable(true).
		SetPermissionEnable(true)
	err = middleware.Get().PermissionCheck(r, w, entry)
	if err != nil {
		term.Failed(err)
		return
	}

	// 读取请求
	in := task.NewDebugJobTaskRequest(r.PathParameter("id"))
	if err = term.ReadReq(in); err != nil {
		term.Failed(err)
		return
	}

	// 进入容器
	in.SetWebTerminal(term)
	h.service.DebugJobTask(r.Request.Context(), in)

	term.Success("ok")
}

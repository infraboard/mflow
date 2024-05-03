package webhook

import (
	"context"
	"sync"

	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mflow/apps/pipeline"
	"github.com/infraboard/mflow/apps/task"
	"github.com/rs/zerolog"
)

func NewWebHook() *WebHook {
	return &WebHook{
		log: log.Sub("webhook"),
	}
}

type WebHook struct {
	log *zerolog.Logger
}

func (h *WebHook) SendTaskStatus(ctx context.Context, hooks []*pipeline.WebHook, t task.TaskMessage) {
	if t == nil {
		return
	}

	if len(hooks) == 0 {
		return
	}

	if len(hooks) > MAX_WEBHOOKS_PER_SEND {
		t.AddErrorEvent("too many webhooks configs current: %d, max: %d", len(hooks), MAX_WEBHOOKS_PER_SEND)
		return
	}

	h.log.Debug().Msgf("start send job task[%s] webhook, total %d", t.ShowTitle(), len(hooks))
	wg := &sync.WaitGroup{}
	for i := range hooks {
		req := newJobTaskRequest(hooks[i], t, wg, h.log)
		go req.Push(ctx)
	}
	wg.Wait()
}

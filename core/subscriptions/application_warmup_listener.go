package subscriptions

import (
	"context"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/application"
)

type ApplicationWarmupListener struct {
	server application.Warmupper
}

func NewApplicationWarmupListener(server application.Warmupper) *ApplicationWarmupListener {
	return &ApplicationWarmupListener{
		server: server,
	}
}

func (r *ApplicationWarmupListener) OnNotify(ctx context.Context) {
	go func() {
		loaded := r.server.Loaded(ctx)
		r.server.SetIsReady(loaded)
	}()
}

func (r *ApplicationWarmupListener) MustSubscribe() bool {
	return r.server.HasWarmup()
}

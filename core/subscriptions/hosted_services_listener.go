package subscriptions

import (
	"context"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/application"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/services"
)

type HostedServiceListener struct {
	server application.SelfHosting
}

func NewHostedServiceListener(server application.SelfHosting) *HostedServiceListener {
	return &HostedServiceListener{
		server: server,
	}
}

func (r *HostedServiceListener) OnNotify(ctx context.Context) {
	hostedServices := r.server.GetHostedServices()
	for i := 0; i < len(hostedServices); i++ {
		hostedService := hostedServices[i]
		go func(ctx context.Context, hostedService services.IHostedService) {
			hostedService.Execute(ctx)
		}(ctx, hostedService)
	}
}

func (r *HostedServiceListener) MustSubscribe() bool {
	return len(r.server.GetHostedServices()) > 0
}

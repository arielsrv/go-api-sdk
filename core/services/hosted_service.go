package services

import (
	"context"
)

type IHostedService interface {
	Execute(ctx context.Context)
}

type HostedService struct {
	hostedService IHostedService
}

func NewHostedService(hostedService IHostedService) *HostedService {
	return &HostedService{
		hostedService: hostedService,
	}
}

func (h HostedService) Execute(ctx context.Context) {
	h.hostedService.Execute(ctx)
}

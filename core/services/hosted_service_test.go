package services_test

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/services"
)

type MockHostedService struct {
	mock.Mock
}

func (m *MockHostedService) Execute(_ context.Context) {
	time.Sleep(100 * time.Millisecond)
}

func TestHostedService_Execute(t *testing.T) {
	mockHostedService := new(MockHostedService)
	hostedService := services.NewHostedService(mockHostedService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		hostedService.Execute(ctx)
	}()

	time.Sleep(200 * time.Millisecond)

	wg.Wait()

	mockHostedService.AssertExpectations(t)
}

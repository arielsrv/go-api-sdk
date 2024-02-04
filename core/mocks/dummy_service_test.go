package mocks_test

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/mocks"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/services"
)

func TestDummyHostedService_Execute(t *testing.T) {
	hostedService := services.NewHostedService(mocks.NewDummySQSConsumerService())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		hostedService.Execute(ctx)
	}()

	go func() {
		time.Sleep(2000 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		require.NoError(t, err)
	}()

	t.Log()
}

func TestDummyHostedService_ExecuteStop(t *testing.T) {
	hostedService := services.NewHostedService(mocks.NewDummySQSConsumerService())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		hostedService.Execute(ctx)
	}()

	go func() {
		time.Sleep(2000 * time.Millisecond)
		stop()
	}()

	t.Log()
}

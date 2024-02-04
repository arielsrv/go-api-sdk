package mocks

import (
	"context"
	"time"

	log "gitlab.com/iskaypetcom/digital/sre/tools/dev/go-logger"
)

type DummySQSConsumerService struct {
}

func NewDummySQSConsumerService() *DummySQSConsumerService {
	return &DummySQSConsumerService{}
}

func (p DummySQSConsumerService) Execute(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Infof("consumer service stopped")
			return
		default:
			for {
				log.Infof("getting messages from queue ...")
				time.Sleep(1 * time.Minute)
			}
		}
	}
}

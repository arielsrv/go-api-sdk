package workers

import (
	"context"
	"time"

	log "gitlab.com/iskaypetcom/digital/sre/tools/dev/go-logger"
)

type HelloWorldService struct {
}

func (r HelloWorldService) Execute(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info("stopping ...")
			return
		default:
			for {
				log.Infof("hello world from Service ...")
				time.Sleep(1 * time.Minute)
			}
		}
	}
}

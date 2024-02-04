package main

import (
	"context"
	"time"

	log "gitlab.com/iskaypetcom/digital/sre/tools/dev/go-logger"
)

type MyBackgroundWorker struct {
}

func (r MyBackgroundWorker) Execute(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info("stopping ...")
			return
		default:
			for {
				log.Infof("hello world from Service ...")
				time.Sleep(2000 * time.Millisecond)
			}
		}
	}
}

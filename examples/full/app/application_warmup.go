package app

import (
	"context"
	log "gitlab.com/iskaypetcom/digital/sre/tools/dev/go-logger"
	"time"
)

type ApplicationWarmup struct {
}

// Loaded is called when the application is loaded
func (r ApplicationWarmup) Loaded(_ context.Context) bool {
	time.Sleep(1000 * time.Millisecond)
	log.Println("Warmup loaded")
	return true
}

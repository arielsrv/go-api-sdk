package application

import "context"

type Warmup interface {
	Loaded(ctx context.Context) bool
}

package subscriptions

import (
	"context"
)

type Listener interface {
	MustSubscribe() bool
	OnNotify(ctx context.Context)
}

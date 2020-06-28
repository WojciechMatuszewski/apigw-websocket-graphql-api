package pubsub

import (
	"context"
)

type CtxKey string

const ctxKeyMetadata = CtxKey("METADATA")

type CtxMetada struct {
	GqlID  string
	ID     string
	Params string
}

func WithContextMetadata(ctx context.Context, metadata CtxMetada) context.Context {
	return context.WithValue(ctx, ctxKeyMetadata, metadata)
}

func metadataFromContext(ctx context.Context) *CtxMetada {
	v, found := ctx.Value(ctxKeyMetadata).(CtxMetada)
	if !found {
		return nil
	}

	return &v
}

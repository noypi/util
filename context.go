package util

import (
	"context"
)

func WithCtxKVMap(ctx context.Context, m map[interface{}]interface{}) context.Context {
	for k, v := range m {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx

}

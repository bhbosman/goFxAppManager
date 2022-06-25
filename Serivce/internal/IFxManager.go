package internal

import "context"

type IFxManager interface {
	StopAll(ctx context.Context) error
	StartAll(ctx context.Context) error
	Stop(ctx context.Context, name ...string) error
	Start(ctx context.Context, name ...string) error
}

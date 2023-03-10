package app

import "context"

type App interface {
	Start() error
	Stop(ctx context.Context) error
}

package service

import "context"

//go:generate mockgen -source=file.go -destination=./mocks/file.go -package=mocks
type FileRepository interface {
	Save(ctx context.Context, data []byte, fileName string) error
	Read(ctx context.Context, fileName string) ([]byte, error)
	Delete(ctx context.Context, fileName string) error
}

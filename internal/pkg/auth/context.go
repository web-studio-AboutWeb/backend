package auth

import (
	"context"

	"web-studio-backend/internal/app/domain"
)

const userKey = "user_id"

func NewContext(ctx context.Context, ac *domain.AuthContext) context.Context {
	return context.WithValue(ctx, userKey, ac)
}

func FromContext(ctx context.Context) *domain.AuthContext {
	ac, ok := ctx.Value(userKey).(*domain.AuthContext)
	if !ok {
		// Panic here because this will break application.
		// Use this function only in appropriate places.
		panic("cannot find user_id in context")
	}

	return ac
}

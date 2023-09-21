package auth

import (
	"context"

	"web-studio-backend/internal/app/domain"
)

const userKey = "user_id"

// NewContext puts auth context in given context.
func NewContext(ctx context.Context, ac *domain.AuthContext) context.Context {
	return context.WithValue(ctx, userKey, ac)
}

// MustFromContext returns auth context from given context.
//
//	Panics if auth context not found.
func MustFromContext(ctx context.Context) *domain.AuthContext {
	ac, ok := ctx.Value(userKey).(*domain.AuthContext)
	if !ok {
		// Panic here because this will break the application.
		// Use this function only in appropriate places.
		panic("cannot find user_id in context")
	}

	return ac
}

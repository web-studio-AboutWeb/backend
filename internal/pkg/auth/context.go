package auth

import "context"

const userKey = "user_id"

func NewContext(ctx context.Context, id int32) context.Context {
	return context.WithValue(ctx, userKey, id)
}

func FromContext(ctx context.Context) int32 {
	id, ok := ctx.Value(userKey).(int32)
	if !ok {
		// Panic here because this will break application.
		// Use this function only in appropriate places.
		panic("cannot find user_id in context")
	}

	return id
}

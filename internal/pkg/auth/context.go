package auth

import "context"

const userKey = "user_id"

func NewContext(ctx context.Context, id int16) context.Context {
	return context.WithValue(ctx, userKey, id)
}

func FromContext(ctx context.Context) int16 {
	id, ok := ctx.Value(userKey).(int16)
	if !ok {
		// Panic here because this will break application.
		// Use this function only in appropriate places.
		panic("cannot find user_id in context")
	}

	return id
}

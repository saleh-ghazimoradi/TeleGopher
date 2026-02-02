package utils

import "context"

type ContextKey string

const (
	UserIdKey   ContextKey = "user_id"
	UserNameKey ContextKey = "user_name"
	PlatformKey ContextKey = "X-Platform"
)

func WithUserId(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, UserIdKey, id)
}

func WithIdFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIdKey).(int64)
	return id, ok
}

func WithUserName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, UserNameKey, name)
}

func UserNameFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(UserNameKey).(string)
	return name, ok
}

func WithPlatform(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, PlatformKey, platform)
}

func PlatformFromContext(ctx context.Context) (string, bool) {
	platform, ok := ctx.Value(PlatformKey).(string)
	return platform, ok
}

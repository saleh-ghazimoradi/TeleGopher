package utils

import "context"

type ContextKey string

const (
	UserIdKey   ContextKey = "user_id"
	NameKey     ContextKey = "name"
	PlatformKey ContextKey = "X-Platform"
)

func WithUserId(ctx context.Context, id uint) context.Context {
	return context.WithValue(ctx, UserIdKey, id)
}

func WithName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, NameKey, name)
}

func WithPlatform(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, PlatformKey, platform)
}

func UserIdFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(UserIdKey).(uint)
	return id, ok
}

func NameFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(NameKey).(string)
	return name, ok
}

func PlatformFromContext(ctx context.Context) (string, bool) {
	platform, ok := ctx.Value(PlatformKey).(string)
	return platform, ok
}

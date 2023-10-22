package helper

import (
	"context"
)

const (
	chatIdKey = "chatID"
)

func CtxWithChatIdValue(ctx context.Context, chatID int64) context.Context {
	return context.WithValue(context.Background(), chatIdKey, chatID)
}

func ChatIdFromCtx(ctx context.Context) int64 {
	return ctx.Value(chatIdKey).(int64)
}

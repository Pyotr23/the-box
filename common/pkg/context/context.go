package context

import "context"

type macAddressKey struct{}

func ContextWithMacAddress(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, macAddressKey{}, addr)
}

func MacAddressFromContext(ctx context.Context) string {
	addr, ok := ctx.Value(macAddressKey{}).(string)
	if !ok {
		return ""
	}
	return addr
}

package cli

import "context"

type contextKey string

const commandContextKey contextKey = "command-context"

func withCommandContext(ctx context.Context, value *commandContext) context.Context {
	return context.WithValue(ctx, commandContextKey, value)
}

func getCommandContext(ctx context.Context) *commandContext {
	value, _ := ctx.Value(commandContextKey).(*commandContext)
	if value == nil {
		return &commandContext{}
	}
	return value
}

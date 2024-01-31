package ego

import (
	"context"
)

type Job struct {
	ctx  context.Context
	f    FuncArgs
	args []any
}

func NewJob(ctx context.Context, f FuncArgs, args ...any) Job {
	return Job{ctx: ctx, f: f, args: args}
}

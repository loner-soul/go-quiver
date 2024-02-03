package ego

import (
	"context"
)

type Job struct {
	ctx  context.Context
	f    FuncArgs
	args []any
}

// JobQueue 任务缓冲队列
type JobQueue interface {
	Push(job Job)
	Pop() (j Job, ack func(), ok bool) // j : 任务； ack:确认消费后长度-1; ok close后返回false
	Len() int64
}

func NewJob(ctx context.Context, f FuncArgs, args ...any) Job {
	return Job{ctx: ctx, f: f, args: args}
}

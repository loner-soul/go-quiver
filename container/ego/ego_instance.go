package ego

import (
	"context"

	"github.com/loner-soul/go-quiver/util/cache"
)

const (
	DefaultName = "ego-default-instance"
)

var (
	instances = cache.NewSafeCache(func(name string) (*Ego, error) {
		return New(WithSize(5000)), nil
	})
	defaultEgo = Instance()
)

func Instance(names ...string) *Ego {
	name := DefaultName
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return instances.MustGet(name)
}

func SetDefaultEgoSize(size int64) {
	defaultEgo.SetEgoSize(size)
}

func Runf(ctx context.Context, task FuncArgs, args ...any) {
	defaultEgo.Runf(ctx, task, args)
}

func Run(ctx context.Context, task Func) {
	defaultEgo.Run(ctx, task)
}

func Done() {
	defaultEgo.Done()
}

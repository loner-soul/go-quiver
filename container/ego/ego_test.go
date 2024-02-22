package ego

import (
	"context"
	"fmt"
	"testing"
)

func Test_EgoUse(t *testing.T) {
	pool := New()

	pool.Runf(context.Background(), func(ctx context.Context, args ...any) {
		fmt.Println("call in goroutine")
	})
	pool.Run(context.Background(), func(ctx context.Context) {
		fmt.Println("call in goroutine")
	})

	pool.Done()
	fmt.Println("done")
}

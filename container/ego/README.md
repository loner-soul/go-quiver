# ego
一个简单的goroutine pool

1. 设计一个简单的goroutine pool
- 安全的启动go; 内置recover
- 限制goroutine数量
- 惰性增加goroutine
- 使用简单


example

```go
package main

import (
	"context"
	"fmt"

	"github.com/loner-soul/go-quiver/container/ego"
)

func main() {
	pool := ego.New()
	pool.Run(context.Background(), func(ctx context.Context) {
		fmt.Println("run in goroutine")
	})
	pool.Close()
}
```


# TODO
- [] 使用sync.Pool 复用goroutine/或者task
- [] 支持withCtx => group (可以分组等待)

- [] 验证job使用指针传递是否更好
- [] 压测比较

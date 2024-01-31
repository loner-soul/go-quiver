# ego
一个简单的goroutine pool

1. 设计一个简单的goroutine pool
- 安全的启动go; 内置recover
- 限制goroutine数量
- 惰性增加goroutine
- 使用简单


example
```go
func main() {
	pool := ego.New()
	pool.SetSize(100)
	a := 1
	b := 2
	pool.Run(func(a, b int){
		// ....
	}, a,b)

	pool.Wait()
}
```


# TODO
- [] 使用sync.Pool 复用goroutine/或者task
- [] 支持withCtx => group (可以分组等待)

- [] 验证job使用指针传递是否更好

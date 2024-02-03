package queue

// JobQueue 任务缓冲队列
type JobQueue interface {
	Push(T)
	Pop() (T, bool)
	Len() int64
}

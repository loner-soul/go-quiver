package queue

type Queue[T any] interface {
	Push(T)
	Pop() (T, bool)
	Len() int64
}

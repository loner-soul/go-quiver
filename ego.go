package ego

import (
	"sync"
)

type Ego struct {
	mutex sync.Mutex
	opt   Option
}

func New() *Ego {
	return &Ego{}
}

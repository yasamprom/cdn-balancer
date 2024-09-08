package usecases

import (
	"sync/atomic"
)

type usecases struct {
	// counter is used for counting requests
	counter atomic.Uint32
}

func New() *usecases {
	return &usecases{}
}

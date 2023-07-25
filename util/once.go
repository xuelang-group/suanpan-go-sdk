package util

import (
	"sync"
	"sync/atomic"
)

type OnceWithoutErr struct {
	m    sync.Mutex
	done uint32
}

func (o *OnceWithoutErr) Do(f func() error) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}

	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		if err := f(); err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
}

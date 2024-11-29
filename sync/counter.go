package counter

import "sync"

type Counter struct {
	mutex sync.Mutex
	value uint
}

func (c *Counter) Inc() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value++
}

func (c *Counter) Value() uint {
	return c.value
}

func NewCounter() *Counter {
	return &Counter{}
}

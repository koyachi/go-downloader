package downloader

type IntCounter struct {
	val int
}

func NewIntCounter() *IntCounter {
	return &IntCounter{val: 0}
}

func (c *IntCounter) Incr(n int) {
	c.val += n
}

func (c *IntCounter) Value() int {
	return c.val
}

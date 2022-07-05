package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup // wait call return
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex // protect m
	m  map[string]*call
}

// Do only return the first one calling
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// someone has calling, just wait
		c.wg.Wait()
		return c.val, c.err
	}
	// no one called
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// do the real call
	c.val, c.err = fn()
	c.wg.Done()

	// release this call
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

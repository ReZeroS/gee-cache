package geecache

import (
	"fmt"
	"log"
	"sync"
)

// Getter load cache from datasource
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc function type, not function instance, can do force convert instead of being called
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	// f(key) not f.Get(key)
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	mu.Lock()
	defer mu.Unlock()
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLocker()
	group := groups[name]
	defer mu.RUnlock()
	return group
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if cache, ok := g.mainCache.get(key); ok {
		log.Printf("[Geecache hit]")
		return cache, nil
	}
	// load from datasource
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

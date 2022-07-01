package lru

import "container/list"

type Cache struct {
	maxBytes         int64 // max allowed memory
	usedBytes        int64 // has used memory
	doublyLinkedList *list.List
	cache            map[string]*list.Element

	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:         maxBytes,
		doublyLinkedList: list.New(),
		cache:            make(map[string]*list.Element),
		OnEvicted:        onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	// hit cache, move front
	if ele, ok := c.cache[key]; ok {
		c.doublyLinkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	// miss cache
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.doublyLinkedList.Back()
	if ele != nil {
		c.doublyLinkedList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.doublyLinkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.doublyLinkedList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	//defensive programming
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.doublyLinkedList.Len()
}

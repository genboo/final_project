package cache

import (
	"sync"
)

type Key interface{}

type Cache struct {
	sync.Mutex
	capacity  int
	queue     List
	items     map[Key]*ListItem
	OnEvicted func(interface{})
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *Cache) Set(key Key, value interface{}) bool {
	lc.Lock()
	defer lc.Unlock()
	val, ok := lc.items[key]
	if ok {
		// если элемент уже был, обновить значение и переместить в начало
		val.Value = value
		lc.queue.MoveToFront(val)
	} else {
		// если не было, добавить
		item := lc.queue.PushFront(value)
		lc.items[key] = item
		// вытолкнуть последний при превышении ёмкости
		if lc.queue.Len() > lc.capacity {
			for k, v := range lc.items {
				if v == lc.queue.Back() {
					if lc.OnEvicted != nil {
						lc.OnEvicted(v.Value)
					}
					lc.queue.Remove(v)
					delete(lc.items, k)
					break
				}
			}
		}
	}
	return ok
}

func (lc *Cache) Get(key Key) (interface{}, bool) {
	lc.Lock()
	defer lc.Unlock()
	val, ok := lc.items[key]
	// если нашелся, переместить в начало
	if ok {
		lc.queue.MoveToFront(val)
		return val.Value, true
	}
	return nil, false
}

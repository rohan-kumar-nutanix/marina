/*
 * Copyright (c) 2018 Nutanix Inc. All rights reserved.
 *
 * Author: isha.singhal@nutanix.com
 *
 * This file adds the utility for a LRU cache.
 */

package misc

import (
	"container/list"
	"sync"
)

// EvictCallback is called when cache entry is evicted.
// Parameters:
//   key: The key of the element of LRU cache which is evicted.
//   value: The value of the element which is evicted.
type EvictCallback func(key string, value interface{})

// LRU cache structure.
type Cache struct {
	// Max size of the cache. Not guaranteed to be maintained if the pinned
	// elements exceeds this capacity as we don't evict pinned elements.
	capacity int
	// Map from the key in the cache to the linked list element actually
	// stored in the cache. The list element is of type entry. If the entry
	// is pinned, the value is an element in the pinnedList. Otherwise, the
	// value is an element in the lruList.
	items map[string]*list.Element
	// Doubly linked list for LRU implementation.
	lruList *list.List
	// The callback function which is called when cache entry is evicted.
	evictCallback EvictCallback
	// Mutex for thread safety for the cache.
	mutex *sync.RWMutex
	// List of pinned items. The items in the list are not evicted
	// automatically. At any given point, an item in the cache is either in
	// the pinnedList or the lruList but not both.
	pinnedList *list.List
}

// Element that is actually stored in the cache.
type entry struct {
	key   string
	value interface{}
	// Number of times the entry has been pinned.
	numPins int
}

// Returns the new cache of a given size.
func NewCache(size int, evictCallback EvictCallback) *Cache {
	cache := &Cache{
		capacity:      size,
		items:         make(map[string]*list.Element),
		lruList:       list.New(),
		evictCallback: evictCallback,
		mutex:         &sync.RWMutex{},
		pinnedList:    list.New(),
	}
	return cache
}

// Adds a value to the cache.
func (cache *Cache) Add(key string, val interface{}) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	element, ok := cache.items[key]
	if ok {
		// If the element is not pinned, move it to the front of the lruList.
		if !cache.isPinned(element) {
			cache.lruList.MoveToFront(element)
		}
		element.Value.(*entry).value = val
		return
	}

	// Add in the cache since it is not already present.
	cacheEntry := &entry{key, val, 0 /*numPins*/}
	linkedListElem := cache.lruList.PushFront(cacheEntry)
	cache.items[key] = linkedListElem

	// Check if cache size exceeds.
	evict := cache.Len() > cache.capacity
	if evict {
		_, _, _ = cache.removeOldest()
	}
}

// Get looks up a key's value from the cache.
func (cache *Cache) Get(key string) (value interface{}, ok bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	element, ok := cache.items[key]
	if ok {
		// If the element is not pinned, move it to the front of the lruList.
		if !cache.isPinned(element) {
			cache.lruList.MoveToFront(element)
		}
		return element.Value.(*entry).value, true
	}
	return nil, false
}

// This method removes the key from the cache. It returns True if the key was
// found and successfully removed, else it returns false. A pinned element cannot
// be removed.
func (cache *Cache) Remove(key string) (ok bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	element, ok := cache.items[key]
	if ok {
		// Cannot remove a pinned element.
		if cache.isPinned(element) {
			return false
		}
		cache.removeElement(element)
		return true
	}
	return false
}

// RemoveOldest removes the oldest item from the cache.
func (cache *Cache) removeOldest() (key string, value interface{}, ok bool) {
	ent := cache.lruList.Back()
	if ent != nil {
		cache.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return "", nil, false
}

// removeElement is used to remove a given list element from the cache.
func (cache *Cache) removeElement(element *list.Element) {
	cache.lruList.Remove(element)
	keyValue := element.Value.(*entry)
	delete(cache.items, keyValue.key)
	if cache.evictCallback != nil {
		cache.evictCallback(keyValue.key, keyValue.value)
	}
}

// Len returns the number of items in the cache.
func (cache *Cache) Len() int {
	return cache.lruList.Len() + cache.pinnedList.Len()
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
func (cache *Cache) Contains(key string) bool {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	_, ok := cache.items[key]
	return ok
}

// Purge is used to completely clear the cache.
func (cache *Cache) Purge() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	for k, v := range cache.items {
		if cache.evictCallback != nil {
			cache.evictCallback(k, v.Value.(*entry).value)
		}
		delete(cache.items, k)
	}
	cache.lruList.Init()
	cache.pinnedList.Init()
}

// Keys returns a slice of the keys in the cache. The unpinned keys ordered are
// oldest to newest in the slice. The pinned keys are at the end of the slice
// in the order they were inserted into the cache.
func (cache *Cache) Keys() []string {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	keys := make([]string, len(cache.items))
	i := 0
	for ent := cache.lruList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	for ent := cache.pinnedList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}

	return keys
}

// RemoveOldest removes the oldest item from the cache.
func (cache *Cache) RemoveOldest() (key string, value interface{}, ok bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	key, val, ok := cache.removeOldest()
	return key, val, ok
}

// GetOldest returns the oldest entry
func (cache *Cache) GetOldest() (key string, value interface{}, ok bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	ent := cache.lruList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return "", nil, false
}

// Pin an element already existing in the cache. Returns true if the element is
// in the cache and was pinned successfully. Return false otherwise.
func (cache *Cache) Pin(key string) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	element, ok := cache.items[key]
	// If the element is not present in the cache there is no way to pin it.
	if !ok {
		return false
	}
	return cache.pin(key, element)
}

// Insert a pinned element in the cache. Returns a bool value.
func (cache *Cache) AddPinned(key string, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	element, ok := cache.items[key]
	// If the element is already present in the cache, just pin it.
	if ok {
		element.Value.(*entry).value = value
		return cache.pin(key, element)
	}

	// New element. Add to the cache's pinned list.
	cacheEntry := &entry{key, value, 1 /*numPins*/}
	pinnedListElem := cache.pinnedList.PushFront(cacheEntry)
	cache.items[key] = pinnedListElem

	// Check if cache size exceeds the max capacity and evict if necessary.
	if cache.Len() > cache.capacity {
		_, _, _ = cache.removeOldest()
	}
	return true
}

// Unpin an element in the cache. If the element is not present in the cache
// return false.
func (cache *Cache) Unpin(key string) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	element, ok := cache.items[key]

	// If the element is not present in the cache there is no way to unpin it.
	if !ok {
		return false
	}

	cache.unpin(key, element)
	// As we have unpinned an element, let us try to keep the cache length
	// within the capacity. This is useful when the the number of pinned
	// items exceeds the max cache capacity.
	if cache.Len() > cache.capacity {
		_, _, _ = cache.removeOldest()
	}
	return true
}

// Returns true if the element is pinned. False, otherwise.
func (cache *Cache) isPinned(element *list.Element) bool {
	return element.Value.(*entry).numPins > 0
}

// Internal method to perform the pinning. Must be called with a lock.
func (cache *Cache) pin(key string, element *list.Element) bool {
	// If we are pinning it the first time, move the element from lruList to
	// pinnedList and then update the map so that the key points to the element
	// in the pinnedList.
	if !cache.isPinned(element) {
		cache.lruList.Remove(element)
		pinnedListElem := cache.pinnedList.PushFront(element.Value)
		cache.items[key] = pinnedListElem
	}
	// Increment the number of pins.
	element.Value.(*entry).numPins++
	return true
}

// Internal method to unpin an element. Must be called with a lock.
func (cache *Cache) unpin(key string, element *list.Element) bool {
	element.Value.(*entry).numPins--
	// If the element is still pinned, return. Else, move the element to the
	// lruList and then update the map so that the key now points to the element
	// in the lruList.
	if cache.isPinned(element) {
		return true
	} else {
		cache.pinnedList.Remove(element)
		lruListElem := cache.lruList.PushFront(element.Value)
		cache.items[key] = lruListElem
	}
	return true
}

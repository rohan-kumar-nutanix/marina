// Copyright (c) 2017 Nutanix Inc. All rights reserved.
//
// Author: zhi.xie@nutanix.com
//
// Implements thread safe slice, wraps the builtin slice
// with read/write mutex.
package misc

import (
	"sync"
)

// Slice type that can be safely shared between goroutines
type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

// Concurrent slice item
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// Appends an item to the concurrent slice and returns the current size
func (cs *ConcurrentSlice) Append(item interface{}) int {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
	return len(cs.items)
}

// Retrieve an item given at the specified index
func (cs *ConcurrentSlice) Get(index int) interface{} {
	cs.Lock()
	defer cs.Unlock()

	return cs.items[index]
}

// Removes and returns an item at the specified index
func (cs *ConcurrentSlice) Remove(index int) interface{} {
	cs.Lock()
	defer cs.Unlock()

	if len(cs.items) <= 0 {
		return nil
	}

	item := cs.items[index]
	tmpSlice := cs.items
	tmpSlice = append(tmpSlice[:index], tmpSlice[index+1:]...)
	cs.items = tmpSlice
	return item
}

// Get the current size of the slice
func (cs *ConcurrentSlice) Size() int {
	cs.Lock()
	defer cs.Unlock()
	return len(cs.items)
}

// Iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)
	f := func() {
		cs.Lock()
		defer cs.Lock()
		for index, value := range cs.items {
			c <- ConcurrentSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

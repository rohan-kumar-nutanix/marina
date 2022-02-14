/*
 * Copyright (c) 2020 Nutanix Inc. All rights reserved.
 * Author: monica.jeyachandran@nutanix.com
 *
 * This module handles set based operations on Version 4 UUIDs.
 *
 */

package uuid4

/* Set of UUIDs. Empty struct occupies 0 memory, hence used as
   value in the map. */
type UuidSet map[Uuid]struct{}

func UuidSetFromList(uuidList []Uuid) UuidSet {
	set := make(UuidSet)
	for _, uuid := range uuidList {
		set[uuid] = struct{}{}
	}
	return set
}

func UuidListFromSet(uSet UuidSet) []Uuid {
	uList := make([]Uuid, 0, len(uSet))
	for uuid, _ := range uSet {
		uList = append(uList, uuid)
	}
	return uList
}

func UnionUuid(sets ...UuidSet) UuidSet {
	set := make(UuidSet)
	for _, arr := range sets {
		for uuid, _ := range arr {
			if _, ok := set[uuid]; !ok {
				set[uuid] = struct{}{}
			}
		}
	}
	return set
}

func IntersectUuid(sets ...UuidSet) UuidSet {
	// store map of values to number of times
	// the value was encountered.
	uuidFreq := make(map[Uuid]int)
	numSets := len(sets)
	for _, arr := range sets {
		for uuid, _ := range arr {
			uuidFreq[uuid]++
		}
	}
	set := make(UuidSet)
	for uuid, freq := range uuidFreq {
		if freq == numSets {
			set[uuid] = struct{}{}
		}
	}
	return set
}

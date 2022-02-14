/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Utilities for interacting with strings.
 */

package base

type StringArr []string

func (sa StringArr) Contains(s string) bool {
	for _, s2 := range sa {
		if s2 == s {
			return true
		}
	}
	return false
}

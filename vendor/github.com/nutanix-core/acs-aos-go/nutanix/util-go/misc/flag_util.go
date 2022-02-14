/*
 * Copyright (c) 2019 Nutanix Inc. All rights reserved.
 *
 * Author: arun.vasudev@nutanix.com
 *
 * This file defines custom flag types by satisfying the Value interface
 * (https://golang.org/pkg/flag/#Value).
 *
 */

package misc

import (
	"strconv"
)

//-----------------------------------------------------------------------------
// Port number (uint16).

type PortNumber struct {
	Port *uint16
}

func (v PortNumber) String() string {
	if v.Port != nil {
		return strconv.FormatUint(uint64(*v.Port), 10)
	}
	return ""
}

func (v PortNumber) Set(s string) error {
	if u, err := strconv.ParseUint(s, 10, 16); err != nil {
		return err
	} else {
		*v.Port = uint16(u)
	}
	return nil
}

//-----------------------------------------------------------------------------

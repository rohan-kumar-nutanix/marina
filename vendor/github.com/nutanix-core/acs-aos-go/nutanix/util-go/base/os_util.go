/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Utilities for interacting with the OS. Inspired by util/base/os_util.cc
 */

package base

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

func GetCommandExitCode(err error) (int, error) {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), nil
		}
	}

	return -1, fmt.Errorf("Cannot get the exit code from the error %q", err)
}

func GetCurrentRSSFromProc() (int64, error) {
	dat, err := ioutil.ReadFile("/proc/self/statm")
	if err != nil {
		return 0, err
	}

	var total_pages, resident_pages int64
	fmt.Sscanf(string(dat), "%d %d",
		&total_pages, &resident_pages)
	pagesize := os.Getpagesize()
	return resident_pages * int64(pagesize), nil
}

func CheckRssLimit(rssLimit int64) (int64, error) {
	rssbytes, err := GetCurrentRSSFromProc()
	if err != nil {
		return -1, err
	}
	rssKib := rssbytes / 1024
	if rssKib > rssLimit {
		os.Exit(2)
	}
	return rssKib, nil
}

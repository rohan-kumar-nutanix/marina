// Copyright (c) 2017 Nutanix Inc. All rights reserved.
//
// Author: mantri@nutanix.com (Hemanth Mantri)
//
// Utility methods for file operations.

package misc

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	// Flags
	fileReadRetryIntervalInSecs = flag.Int(
		"file_read_retry_intervall_in_secs",
		10,
		"Interval in seconds after which reading of config files should be"+
			" retried.")

	fileReadNumRetries = flag.Int(
		"file_read_num_retries",
		6,
		"Number of times the config file read should be retried in case of "+
			"some intermittent error.")
)

// ============================================================================

// The code inside this block has been taken as is from io/ioutil/tempfile.go
// (https://github.com/golang/go/blob/master/src/io/ioutil/tempfile.go)
// The default ioutil.TmpFile implementation only takes prefix and no suffix
// parameter.
//
// With file extension added to suffix (eg: .yaml), the file editors can
// syntax hilight the text being edited. The following modified TmpFile()
// accepts another argument for the suffix string.

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// TempFile creates a new temporary file in the directory dir
// with a name beginning with prefix, opens the file for reading
// and writing, and returns the resulting *os.File.
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file. The caller can use f.Name()
// to find the pathname of the file. It is the caller's responsibility
// to remove the file when no longer needed.
// Nutanix: Added 'suffix' argument to this method.
func TempFile(dir, prefix string, suffix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextSuffix()+suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}

//=============================================================================

func ReadFileWithBackoff(filePath string) ([]byte, error) {
	cb := NewConstantBackoff(time.Duration(
		*fileReadRetryIntervalInSecs)*time.Second,
		*fileReadNumRetries)
	var fileContent []byte
	var err error
	for {
		fileContent, err = ioutil.ReadFile(filePath)
		if err == nil {
			break
		}

		duration := cb.Backoff()
		if duration == Stop {
			return nil, err
		}
	}
	return fileContent, nil
}

func ReadDirWithBackoff(dirname string) ([]os.FileInfo, error) {
	cb := NewConstantBackoff(time.Duration(
		*fileReadRetryIntervalInSecs)*time.Second,
		*fileReadNumRetries)

	for {
		fileInfo, err := ioutil.ReadDir(dirname)
		if err == nil {
			return fileInfo, nil
		}

		duration := cb.Backoff()
		if duration == Stop {
			return nil, err
		}
	}
}

func OpenWithBackoff(filePath string) (*os.File, error) {
	cb := NewConstantBackoff(time.Duration(
		*fileReadRetryIntervalInSecs)*time.Second,
		*fileReadNumRetries)

	for {
		file, err := os.Open(filePath)
		if err == nil {
			return file, nil
		}

		duration := cb.Backoff()
		if duration == Stop {
			return nil, err
		}
	}
}

func OpenFileWithBackoff(filePath string, flag int, perm os.FileMode) (
	*os.File, error) {
	cb := NewConstantBackoff(time.Duration(
		*fileReadRetryIntervalInSecs)*time.Second,
		*fileReadNumRetries)

	for {
		file, err := os.OpenFile(filePath, flag, perm)
		if err == nil {
			return file, nil
		}

		duration := cb.Backoff()
		if duration == Stop {
			return nil, err
		}
	}
}

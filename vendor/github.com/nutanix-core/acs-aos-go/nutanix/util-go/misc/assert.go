package misc

import (
	"github.com/golang/glog"
)

func Assert(cond bool, logMsg string, args ...interface{}) {
	// TODO (akshay.muramatti) Add support for debug and release build types.
	if !cond {
		glog.Fatalf(logMsg, args...)
	}
}

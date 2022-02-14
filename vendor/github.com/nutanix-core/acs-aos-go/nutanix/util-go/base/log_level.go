package base

import (
	"github.com/golang/glog"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/version"
	"regexp"
)

// Log level verbosity.
const (
	VLOG_RARE glog.Level = iota
	VLOG_PERIODIC
	VLOG_PER_REQUEST_INFO
	VLOG_PER_REQUEST_DEBUG
)

type BuildGetter func() string

type Build struct {
	GetBuildVersion BuildGetter
}

func GetBuildVersion() string {
	return version.BuildVersion
}

func BuildInfo(bd BuildGetter) *Build {
	return &Build{GetBuildVersion: bd}
}

func (build *Build) IsDebugBuild() bool {
	debugPattern, err := regexp.Compile("^[A-Za-z0-9\\.]+-(opt|dbg)-")
	if err != nil {
		glog.Fatal("Failed to compile debug build regex")
	}
	ver := build.GetBuildVersion()
	return debugPattern.MatchString(ver)
}

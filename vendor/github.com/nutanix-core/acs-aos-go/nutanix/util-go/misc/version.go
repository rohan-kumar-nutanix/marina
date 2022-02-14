// Copyright (c) 2019 Nutanix Inc. All rights reserved.
//
// Author: raghu.rapole@nutanix.com
//
// This file defines Golang version of methods defined in:
//   util/infrastructure/version.py
//

package misc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Nutanix software version is stored in release_version file,
// It has specific format like el7.3-opt-euphrates-5.10.1-stable-2a9ae949a7a01b04629ae5616ba78fcd01aad922
var nutanixVersionFormatRe = regexp.MustCompile(`^(\S+)-(\S+)-(\S+)-(\S+)-(\S+)-(\S+)$`)

// Nutanix master branch will have format like below:
// el7.3-opt-master-07e99216f885306187bcab63693584d70bb8c91b
var nutanixVersionMasterFormatRe = regexp.MustCompile(`^(\S+)-(\S+)-(\S+)-(\S+)$`)

// Nutanix -E branch will have format like below:
// el6-release-danube-4.5.1-E-stable-415411da0d7d8fcee1deee6b546d65e75cb2593a
var nutanixVersionEFormatRe = regexp.MustCompile(`^(\S+)-(\S+)-(\S+)-(\S+-E)-(\S+)-(\S+)$`)

// Location of the release version file on CVM.
const relVerFile = "/etc/nutanix/release_version"

// Returns the full release string of the format:
// 'el7.3-opt-euphrates-5.10.1-stable-2a9ae949a7a01b04629ae5616ba78fcd01aad922'
func GetClusterVersion(numeric bool) (string, error) {
	relVer, err := ioutil.ReadFile(relVerFile)
	if err != nil {
		return "", err
	}
	if numeric {
		return ExtractNosVersionFromReleaseVersion(string(relVer))
	} else {
		return string(relVer), nil
	}
}

// Get NOS version name from release version name.
// e.g. if release version name = 1.6-release-congo-3.5.2-stable-d743d74ec then
//				return 3.5.2
// Different regex is used for Master branch, full name is treated as version.
// Returns version on success, None otherwise.
func ExtractNosVersionFromReleaseVersion(relVersion string) (string, error) {
	if relVersion == "" {
		return "", errors.New("Invalid input release version")
	}

	var match []string
	relVersion = strings.TrimSpace(relVersion)
	match = nutanixVersionEFormatRe.FindStringSubmatch(relVersion)
	if match != nil {
		return match[4], nil
	}
	match = nutanixVersionFormatRe.FindStringSubmatch(relVersion)
	if match != nil {
		return match[4], nil
	}
	match = nutanixVersionMasterFormatRe.FindStringSubmatch(relVersion)
	if match != nil {
		return relVersion, nil
	}

	return "", errors.New(fmt.Sprintf("Invalid version format %s", relVersion))
}

// Converts input '.' separated string to list of integers. If any part of the
// string is not convertable to integer or if input is empty string, it will
// return [9999].
func convertStringToIntList(str string) []int {
	strList := strings.Split(str, ".")
	var intList []int
	for _, x := range strList {
		if val, err := strconv.Atoi(x); err != nil {
			return []int{9999}
		} else {
			intList = append(intList, val)
		}
	}
	return intList
}

// Compares version numbers of the format '3.5.2.1'.
// Returns 0 : if v1 == v2
// Returns 1 : if v1 > v2
// Returns -1: if v1 < v2
//
// If either version is not in the correct format, they will be the greater,
// unless neither can be parsed in which case they are equal. The case where a
// branch can't be parsed is if the cluster is running master, or some other
// non-release branch, or empty string.  This is only expected in a test/debug
// situation and is the motivation for making an unparseable format greater.
func CompareVersions(v1 string, v2 string) int {
	if strings.Compare(v1, "master") == 0 {
		v1 = "9999"
	}
	if strings.Compare(v2, "master") == 0 {
		v2 = "9999"
	}

	v1IntList := convertStringToIntList(v1)
	v2IntList := convertStringToIntList(v2)

	var maxLen int
	if len(v1IntList) < len(v2IntList) {
		maxLen = len(v2IntList)
	} else {
		maxLen = len(v1IntList)
	}

	v1NormIntList := make([]int, maxLen)
	v2NormIntList := make([]int, maxLen)
	copy(v1NormIntList, v1IntList)
	copy(v2NormIntList, v2IntList)

	for i, e := range v1NormIntList {
		if e > v2NormIntList[i] {
			return 1
		} else if e < v2NormIntList[i] {
			return -1
		}
	}
	return 0
}

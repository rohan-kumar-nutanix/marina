// +build use_init_nutanix
// Copyright (c) 2017 Nutanix Inc. All rights reserved.
//
// InitNutanix

package misc

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/version"
)

var Version = flag.Bool("version", false, "If true, print version and exit.")

var FlagFile = flag.String("flagfile", "", "")

var UnitTestMode bool = false

func InitNutanix() {
	flag.Parse()
	if *Version {
		PrintVersion()
		os.Exit(0)
	}

	// Set 'log_dir' to default value if it has not been set, retain the
	// value otherwise.
	logDirPath := flag.Lookup("log_dir").Value.String()
	if logDirPath == "" {
		flag.Set("log_dir", "/home/nutanix/data/logs")
	}

	if *FlagFile != "" {
		ParseFlagFile(*FlagFile)
	}
}

//-----------------------------------------------------------------------------

func PrintVersion() {
	progname := os.Args[0]
	fmt.Printf("%s version %s\n", progname, version.BuildVersion)

	// Print padding.
	for range progname {
		fmt.Printf(" ")
	}

	fmt.Printf(" last-commit-date %s\n", version.BuildLastCommitDate)
}

//-----------------------------------------------------------------------------

// TODO: handle more cases. Currently can only handle something like:
// --perf_mode=true. Does not handle something like: -v=3
func ParseFlagFile(flagFile string) {
	f, err := os.Open(flagFile)
	if err != nil {
		fmt.Println("Warning: Couldn't open flag file", flagFile)
		return
	}
	defer f.Close()

	flags := bufio.NewScanner(f)
	flags.Split(bufio.ScanLines)
	for flags.Scan() {
		line := flags.Text()
		parts := strings.Split(line, "=")
		flagname := parts[0][2:]
		value := parts[1]
		flag.Set(flagname, value)
	}
}

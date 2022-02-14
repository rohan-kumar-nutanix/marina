package misc

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/golang/glog"
)

func FindProcess(processName string, args string) (string, bool) {
	grepCmd := "/bin/ps -ewwf" +
		" | /bin/grep '" + processName + "'" +
		" | /bin/grep '" + args + "'" +
		" | /bin/grep -v grep" +
		" | /usr/bin/awk '{print $2}'"

	glog.Infof("Checking for process running with cmd: %s", grepCmd)
	cmd := exec.Command("bash", "-c", grepCmd)
	out, err := cmd.Output()
	if err != nil {
		glog.Errorf("Error while checking for process: %s", err)
		return "", false
	}

	output := string(out)
	return output, !(output == "")
}

func KillProcesses(processName string, args string, signum int,
	skipAfterCheck bool) error {
	grepCmd := "/bin/ps -ewwf" +
		" | /bin/grep '" + processName + "'" +
		" | /bin/grep '" + args + "'" +
		" | /bin/grep -v grep" +
		" | /usr/bin/awk '{print $2}'"

	cmd := exec.Command("bash", "-c", grepCmd)
	out, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("Error while checking for process: %s", err)
	}
	if string(out) == "" {
		glog.Info("No target process to kill.")
		return nil
	}

	killCmd := fmt.Sprintf("%s | /usr/bin/xargs kill -%d", grepCmd, signum)

	return DoKill(killCmd, grepCmd, skipAfterCheck)
}

func DoKill(killCmd string, grepCmd string, skipAfterCheck bool) error {
	var err error
	numTries := 0
	for numTries < 3 {
		glog.Infof("Killing processes with cmd: %s", killCmd)
		cmd := exec.Command("bash", "-c", killCmd)
		_, err = cmd.Output()
		if err != nil {
			glog.Infof("Error while killing process: %s,", err)
			numTries++
		} else {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	if err != nil {
		return err
	}
	if skipAfterCheck {
		return nil
	}

	// Verify that there are no matching processes left. grep will exit with
	// status 1 if there are no matches.
	glog.Infof("Checking that no processes are left with cmd: %s", grepCmd)
	cmd := exec.Command("bash", "-c", grepCmd)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Failed to double check if killed processes are still "+
			"alive with command: %s", grepCmd)
	}
	if string(out) != "" {
		return fmt.Errorf("Unexpected process is still alive after executing kill "+
			"command: %s", killCmd)
	}
	return nil
}

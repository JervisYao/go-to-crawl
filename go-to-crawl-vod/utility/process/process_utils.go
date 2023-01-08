package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//根据进程名判断进程是否运行
func CheckRunning(port int) (string, error) {
	a := fmt.Sprintf("lsof -i:%d|sed -n '2p'|awk '{print $2}'", port)
	result, err := exec.Command("/bin/sh", "-c", a).Output()
	pid := ""
	if err != nil {
		return pid, err
	}
	pid = strings.TrimSpace(string(result))
	return pid, nil
}

func CheckProcessRunning(processName string) (string, error) {
	a := fmt.Sprintf(`ps -ef | grep -v grep | grep %s | sed -n '1p'| awk '{print $2}'`, processName)
	result, err := exec.Command("/bin/sh", "-c", a).Output()
	pid := ""
	if err != nil {
		return pid, err
	}
	pid = strings.TrimSpace(string(result))
	return pid, nil
}

func KillPid(pidStr string) {
	if pidStr == "" {
		return
	}

	pid, _ := strconv.Atoi(pidStr)
	runProcess, _ := os.FindProcess(pid)

	KillProcess(runProcess)

}

func KillProcess(process *os.Process) {
	if process == nil {
		return
	}
	_ = process.Kill()
	processStatus, _ := process.Wait()
	if processStatus != nil && processStatus.Exited() {
		return
	}
}

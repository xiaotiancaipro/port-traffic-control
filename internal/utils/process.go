package utils

import (
	"fmt"
	"os"
	"port-traffic-control/internal/logger"
	"strconv"
	"strings"
	"syscall"
)

func NewProcessUtil(log *logger.Log) *ProcessUtil {
	return &ProcessUtil{
		Log: log,
	}
}

func (pu *ProcessUtil) CheckRunning(file string) (pid int, err error) {

	data, err := os.ReadFile(file)
	if err != nil {
		err = fmt.Errorf("file read failed, Error=%v", err)
		pu.Log.Error(err)
		return
	}

	pid, err = strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		err = fmt.Errorf("type conversion failed, Error=%v", err)
		pu.Log.Error(err)
		return
	}

	// 检查进程是否存在
	process, err := os.FindProcess(pid)
	if err != nil {
		err = fmt.Errorf("the process does not exist, Error=%v", err)
		pu.Log.Error(err)
		return
	}

	// 发送 0 信号检查进程状态
	if err = process.Signal(syscall.Signal(0)); err != nil {
		err = fmt.Errorf("send 0 to signal an error, Error=%v", err)
		pu.Log.Error(err)
		return
	}

	return

}

func (pu *ProcessUtil) WritePIDFile(file string, pid int) error {
	err := os.WriteFile(file, []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		pu.Log.Errorf("Failed to write file, Error=%v", err)
		return err
	}
	return nil
}

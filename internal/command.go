package internal

import (
	"log"
	"os/exec"
)

// Command 封装 os 命令行参数
type Command []string

// NewCommand Command 的构造函数
func NewCommand(cmd string, args ...string) Command {
	return append([]string{cmd}, args...)
}

// Run 执行 os 命令行
func (cmd Command) Run(workDir string) (string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Dir = workDir
	logBuf, err := c.CombinedOutput()
	log.Print(append([]string{workDir}, cmd...), " output:\n", string(logBuf))
	return string(logBuf), err
}

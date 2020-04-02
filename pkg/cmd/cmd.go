package cmd

import (
	"bingo/pkg/log"
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var bash string

func init() {
	if runtime.GOOS == "windows" {
		bash = `C:\Program Files\Git\bin\bash.exe`
	} else {
		bash = "/bin/sh"
	}
}

func GetRemoteCmdStr(remote string, format string, args ...interface{}) string {
	_, str := RunRemoteCmd(remote, format, args...)
	return str
}
func GetCmdStr(format string, args ...interface{}) string {
	_, str := RunCmd(format, args...)
	return str
}
func GetCmdTrimStr(format string, args ...interface{}) string {
	return strings.TrimSpace(GetCmdStr(format, args...))
}

func GetCmd(remote string, format string, args ...interface{}) string {
	cmd := fmt.Sprintf(format, args...)
	if strings.TrimSpace(remote) == "" || strings.TrimSpace(GetHostname(remote)) == strings.TrimSpace(GetHostname()) {
		return cmd
	}
	return fmt.Sprintf("ssh -o \"StrictHostKeyChecking no\" %s \"%s\"", remote, cmd)
}
func RunRemoteCmd(remote string, format string, args ...interface{}) (bool, string) {
	return RunCmd(GetCmd(remote, format, args...))
}
func RunCmd(format string, args ...interface{}) (succ bool, str string) {
	command := fmt.Sprintf(format, args...)
	cmd := exec.Command(bash, "-c", command)
	f := func(c bool, s string) {
		s = strings.TrimSpace(s)
		succ = c
		str = s
		//if c != 0 {
		//	log.Warning("cmd[%s] return %d: %s", command, c, s)
		//} else {
		log.Debug("cmd[%s] return %v: %s", command, c, s)
		//}
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		f(false, fmt.Sprintf("StdoutPipe: %s", err))
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		f(false, fmt.Sprintf("StderrPipe: %s", err))
	}
	if err := cmd.Start(); err != nil {
		f(false, fmt.Sprintf("Start: %s", err))
	}

	bytesErr, err := read(stderr)
	if err != nil {
		f(false, fmt.Sprintf("ReadAll stderr: %s", err))
	}
	//if len(bytesErr) != 0 {
	//	f(false, fmt.Sprintf("stderr is not nil: %s", bytesErr))
	//}

	bytes, err := read(stdout)
	if err != nil {
		f(false, fmt.Sprintf("ReadAll stdout: %s", err))
	}
	if err := cmd.Wait(); err != nil {
		f(false, fmt.Sprintf("[%s]%s", err, bytesErr))
		return
	}

	f(true, fmt.Sprintf("%s", bytes))
	return
}
func RunCmdDetail(format string, args ...interface{}) (exitCode int, stdOut, stdErr string) {
	command := fmt.Sprintf(format, args...)
	cmd := exec.Command(bash, "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Debug("StdoutPipe: %s", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Debug("StderrPipe: %s", err)
	}
	if err := cmd.Start(); err != nil {
		log.Debug("Start: %s", err)
	}
	bytesErr, err := read(stderr)
	if err != nil {
		log.Debug("ReadAll stderr: %s", err)
	}
	//if len(bytesErr) != 0 {
	//	f(false, fmt.Sprintf("stderr is not nil: %s", bytesErr))
	//}
	bytes, err := read(stdout)
	if err != nil {
		log.Debug("ReadAll stderr: %s", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Debug("[%s]--%s", err, bytesErr)
		code, _ := strconv.Atoi(strings.Trim(err.Error(), "exit status "))
		return code, string(bytes), string(bytesErr)
	}
	return 0, string(bytes), string(bytesErr)
}
func RunCmdQuiet(format string, args ...interface{}) bool {
	command := fmt.Sprintf(format, args...)
	cmd := exec.Command(bash, "-c", command)
	if err := cmd.Start(); err != nil {
		log.Error("cmd start failed: %s", err)
		return false
	}
	if err := cmd.Wait(); err != nil {
		log.Error("cmd wait failed: %s", err)
		return false
	}
	return true
}
func RunAwkCmd(remote string, remoteCmd, localCmd string) (bool, string) {
	return RunCmd(fmt.Sprintf("ssh -o \"StrictHostKeyChecking no\" %s \"%s\" |%s", remote, remoteCmd, localCmd))
}

func RunRemote(remote string, format string, args ...interface{}) error {
	return Run(GetCmd(remote, format, args...))
}
func Run(format string, args ...interface{}) error {
	command := fmt.Sprintf(format, args...)
	// log.Debug("cmd: %s", command)
	cmd := exec.Command(bash, "-c", command)

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	out := bufio.NewReader(stdout)
	outerr := bufio.NewReader(stderr)

	cmd.Start()
	exit := false
	fn := func() {
		for {
			if exit {
				break
			}
			line, _, _ := out.ReadLine()
			s := strings.TrimSpace(string(line))
			if s != "" {
				fmt.Println(s)
				//fmt.Println(color.Blue(s))
			}
			linerr, _, _ := outerr.ReadLine()
			serr := strings.TrimSpace(string(linerr))
			if s != "" {
				fmt.Println(serr)
				//fmt.Println(color.Red(s))
			}
		}
	}
	go fn()

	read(stderr)
	// bytesErr, _ := read(stderr)

	//if len(bytesErr) != 0 {
	//	fmt.Printf("stderr: %s", string(bytesErr))
	//fmt.Printf("stderr: %s", color.Red(string(bytesErr)))
	//	return log.NewError("stderr: %s", bytesErr)
	//}
	if err = cmd.Wait(); err != nil {
		// return log.NewError("[%s]%s", err, bytesErr)
		return err
	}
	exit = true
	return nil
}

// Read is the simple call of BlockReader
func read(r io.Reader) ([]byte, error) {
	var data []byte
	// defer r.Close()
	bufSize := 1024 * 10
	buf := make([]byte, bufSize) //一次读取多少个字节
	bfRd := bufio.NewReader(r)
	for {
		n, err := bfRd.Read(buf)
		data = append(data, buf[:n]...)
		if err != nil { //遇到任何错误立即返回，并忽略EOF错误信息
			if err == io.EOF {
				return data, nil
			}
			return data, err
		}
	}
	return data, nil
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/er1c-zh/go-now/log"
	"os"
	"syscall"
)

var (
	cmd string
)

func init() {
	flag.StringVar(&cmd, "cmd", "", "empty as cli, daemon as server.")
}

func main() {
	defer log.Flush()
	flag.Parse()
	log.Info("cmd: %s", cmd)
	switch cmd {
	case "daemon":
		(&Daemon{}).Run()
	default:
		var result []ResultItem
		defer func() {
			j, _ := json.Marshal(result)
			fmt.Printf("%s\n", string(j))
			return
		}()
		daemon, err := readStatusFile()
		if err != nil {
			log.Error("readStatusFile fail: %s", err.Error())
			createDaemon()
			return
		}
		if !daemon.Running {
			log.Warn("daemon not running, try start daemon.")
			createDaemon()
			return
		}
		// TODO
	}
}

func createDaemon() {
	oldPid := os.Getpid()
	s, err := os.Executable()
	if err != nil {
		log.Warn("createDaemon fail: %s", err.Error())
		return
	}
	log.Info("cmd: %s, args: %v", s, os.Args)
	pid, err := syscall.ForkExec(s, []string{s, "-cmd", "daemon"}, &syscall.ProcAttr{})
	if err != nil {
		log.Warn("ForkExec fail: %s", err.Error())
		return
	}
	log.Info("create daemon process successfully, pid: %d.", oldPid, pid)
}

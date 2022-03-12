package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/er1c-zh/go-now/log"
	"os"
	"syscall"
)

const (
	logPath = "/var/tmp/devdocs_alfred/access.log"
)

var (
	cmd   string
	query string
)

func init() {
	flag.StringVar(&cmd, "cmd", "", "empty: run as cli.\nkill: kill daemon.\ndaemon: create daemon process.")
	flag.StringVar(&query, "query", "", "param for cli query.")
}

func main() {
	log.SetLogger(log.NewFileLogger(logPath))
	defer log.Flush()
	flag.Parse()
	log.Info("cmd: %s, query: %s", cmd, query)
	switch cmd {
	case "daemon":
		(&Daemon{}).Run()
	case "kill":
		daemon, err := readStatusFile()
		if err != nil {
			log.Error("readStatusFile fail: %s", err.Error())
			return
		}
		if !daemon.Running {
			return
		}
		proc, err := os.FindProcess(daemon.Pid)
		if err != nil {
			log.Error("FindProcess find daemon process fail: %s", err.Error())
			return
		}
		err = proc.Kill()
		if err != nil {
			log.Error("proc.Kill fail: %s", err.Error())
			return
		}
		return
	default:
		var result []ResultItem
		defer func() {
			j, _ := json.Marshal(AlfredResp{Items: result})
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
		cli, err := NewCli(daemon)
		if err != nil {
			result = []ResultItem{{Title: "NewCli fail"}}
			return
		}
		result = cli.Router(cmd, query)
		return
	}
}

func createDaemon() {
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
	log.Info("create daemon process successfully, pid: %d.", pid)
}

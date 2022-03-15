package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/er1c-zh/devdocs_alfred/fuzzy_match"
	"github.com/er1c-zh/go-now/log"
)

type Daemon struct{}

const (
	PidFilePath    = "/var/tmp/devdocs_alfred/pid"
	StatusFilePath = "/var/tmp/devdocs_alfred/status"
)

///////////////
// Daemon api
///////////////

var (
	cmdList = []Cmd{
		{"doc", "查询文档列表"},
		{"cache", "缓存文档"},
		{"update", "更新文档"},
		{"reset", "重置本地设置"},
	}
)

func (d *Daemon) CmdList(req *RpcReq, resp *RpcResp) error {
	log.Info("CmdList, req.Query: %s", req.Query)
	resp.Data = make([]ResultItem, 0, len(cmdList))
	for _, item := range cmdList {
		resp.Data = append(resp.Data, item.ToAlfred())
	}
	if req.Cmd == "" {
		// 没有输入时展示所有的
		return nil
	}
	resp.Data = fuzzy_match.FuzzyMatch[ResultItem](req.Cmd, resp.Data,
		fuzzy_match.NewOption().WithLimit(5).WithMinScore(1).WithDebug())
	return nil
}

func (d *Daemon) DocList(req *RpcReq, resp *RpcResp) error {
	log.Info("DocList, req.Query: %s", req.Query)
	data, err := GetDocsList()
	if err != nil {
		return err
	}
	resp.Data = make([]ResultItem, 0, len(data))
	for _, item := range data {
		resp.Data = append(resp.Data, item.ToAlfred())
	}
	resp.Data = fuzzy_match.FuzzyMatch[ResultItem](req.Query, resp.Data,
		fuzzy_match.NewOption().WithLimit(5))
	return nil
}

func (d *Daemon) SearchDoc(req *RpcReq, resp *RpcResp) error {
	log.Info("SearchDoc, req.Query: %s", req.Query)
	// req.Query -> {doc_slug}_{query}
	queryList := strings.SplitN(req.Query, "_", 2)
	if len(queryList) != 2 {
		log.Warn("SearchDoc invalid query: '%s'", req.Query)
		return nil
	}
	docSlug, _query := queryList[0], queryList[1]
	data, err := GetDocIndex(docSlug)
	if err != nil {
		log.Error("SearchDoc.GetDocIndex fail: %s", err.Error())
		return nil
	}

	data = fuzzy_match.FuzzyMatch(_query, data,
		fuzzy_match.NewOption().WithLimit(10).WithMinScore(1))
	resp.Data = make([]ResultItem, 0, len(data))
	for _, item := range data {
		resp.Data = append(resp.Data, item.ToAlfred(docSlug))
	}
	return nil
}

//////////////////////////////
// Daemon 非接口
//////////////////////////////

// Run 启动Daemon
func (d *Daemon) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("panic: %v", err)
			log.Flush()
			fmt.Printf("panic: %v\n", err)
		}
	}()
	shouldContinue, err, clearFunc := testOrCreatePidFile()
	if err != nil {
		log.Fatal("testOrCreatePidFile fail: %s", err.Error())
		return
	}
	if !shouldContinue {
		return
	}
	if clearFunc != nil {
		defer clearFunc()
	}
	// create a local rpc server
	err = rpc.Register(d)
	if err != nil {
		log.Fatal("rpc.Register fail: %s", err.Error())
		return
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("net.Listen fail: %s", err.Error())
		return
	}

	defer func() {
		_ = l.Close()
	}()

	err = writeStatusFile(statusFileStruct{
		Addr: l.Addr().String(),
		Pid:  os.Getpid(),
	})
	if err != nil {
		log.Fatal("writeStatusFile fail: %s", err.Error())
		return
	}

	go func() {
		// TODO load cache file
	}()
	if err := http.Serve(l, nil); err != nil {
		log.Fatal("http.Server fail: %s", err.Error())
		return
	}
	return
}

type statusFileStruct struct {
	Addr    string
	Pid     int
	Running bool `json:"-"`
}

func writeStatusFile(status statusFileStruct) error {
	if err := os.MkdirAll(filepath.Dir(StatusFilePath), os.FileMode(0755)); err != nil {
		return err
	}
	j, _ := json.Marshal(status)
	if err := ioutil.WriteFile(StatusFilePath, j, 0644); err != nil {
		return err
	}
	return nil
}

func readStatusFile() (statusFileStruct, error) {
	var status statusFileStruct
	statusBytes, err := ioutil.ReadFile(StatusFilePath)
	if err != nil {
		log.Error("readStatusFile ReadFile fail: %s", err.Error())
		return status, err
	}
	err = json.Unmarshal(statusBytes, &status)
	if err != nil {
		log.Error("readStatusFile Unmarshal fail: %s", err.Error())
		return status, err
	}
	status.Running = IsProcessExist(strconv.FormatInt(int64(status.Pid), 10))
	return status, nil
}

func IsProcessExist(pidStr string) bool {
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		return false
	}
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	}
	if process == nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		switch err.Error() {
		case "os: process already finished":
			return false
		case "err: operation not permitted":
			return true
		default:
			log.Warn("unknown signal err: %s", err.Error())
			return false
		}
	}
	return true
}

func testOrCreatePidFile() (created bool, err error, cleanFunc func()) {
	pidContent, err := ioutil.ReadFile(PidFilePath)
	if err == nil {
		if IsProcessExist(strings.TrimSpace(string(pidContent))) {
			log.Info("daemon already running!")
			return false, nil, nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(PidFilePath), os.FileMode(0755)); err != nil {
		return false, err, nil
	}
	if err := ioutil.WriteFile(PidFilePath, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return false, err, nil
	}
	return true, nil, func() {
		_ = os.Remove(PidFilePath)
	}
}

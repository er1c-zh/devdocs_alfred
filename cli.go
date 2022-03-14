package main

import (
	"fmt"
	"github.com/er1c-zh/go-now/log"
	"net/rpc"
)

type Cli struct {
	DaemonStatus statusFileStruct
	client       *rpc.Client
}

func NewCli(daemonStatus statusFileStruct) (*Cli, error) {
	cli := &Cli{
		DaemonStatus: daemonStatus,
	}
	_cli, err := rpc.DialHTTP("tcp", daemonStatus.Addr)
	if err != nil {
		log.Fatal("dial daemon server(%s) fail: %s", daemonStatus.Addr, err.Error())
		return nil, fmt.Errorf("dial server fail: %s", err.Error())
	}
	cli.client = _cli
	return cli, nil
}

func (c *Cli) Router(cmd string, query string) []ResultItem {
	api := ""
	switch cmd {
	case "doc":
		api = "Daemon.DocList"
	case "search":
		api = "Daemon.SearchDoc"
	default:
		api = "Daemon.CmdList"
	}
	if api == "" {
		return []ResultItem{{Title: "unsupported cmd"}}
	}

	var data RpcResp
	log.Info("cli.Call(%s, %s, %s)", api, cmd, query)
	err := c.client.Call(api, &RpcReq{Cmd: cmd, Query: query}, &data)
	if err != nil {
		log.Error("Call %s fail: %s", api, err.Error())
		return GenMsgResultItemList("Call daemon fail.")
	}
	return data.Data
}

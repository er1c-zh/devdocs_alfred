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
	switch cmd {
	case "doc":
		var data RpcResp
		err := c.client.Call("Daemon.DocList", &RpcReq{Query: query}, &data)
		if err != nil {
			log.Error("Call DocList fail: %s", err.Error())
			return []ResultItem{{Title: "call daemon err"}}
		}
		return data.Data
	}
	return []ResultItem{{Title: "unsupported cmd"}}
}

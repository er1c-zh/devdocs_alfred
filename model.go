package main

import "fmt"

type AlfredResp struct {
	Items []ResultItem `json:"items"`
}

type ResultItem struct {
	Type         string `json:"type"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Arg          string `json:"arg"`
	Autocomplete string `json:"autocomplete"`
	Icon         struct {
		Type string `json:"type"`
		Path string `json:"path"`
	} `json:"icon"`
}

func GenMsgResultItemList(msg string) []ResultItem {
	return []ResultItem{
		{Title: msg},
	}
}

func (r ResultItem) GetString() string {
	return r.Title
}

type RpcReq struct {
	Cmd   string
	Query string
}
type RpcResp struct {
	Data []ResultItem
}

////////////////////////////////
// base model
////////////////////////////////

type Cmd struct {
	Title string
	Desc  string
}

func (c Cmd) ToAlfred() ResultItem {
	return ResultItem{
		Title:        c.Title,
		Subtitle:     c.Desc,
		Autocomplete: c.Title,
	}
}

////////////////////////////////
// devdocs model
////////////////////////////////

// RespDocMeta 文档元数据
type RespDocMeta struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Type  string `json:"type"`
	Links struct {
		Home string `json:"home"`
		Code string `json:"code"`
	} `json:"links"`
	Version string `json:"version"`
	Release string `json:"release"`
	Mtime   int    `json:"mtime"`
	DbSize  int    `json:"db_size"`
}

func (m RespDocMeta) ToAlfred() ResultItem {
	var dest ResultItem
	dest.Title = m.Name
	dest.Subtitle = fmt.Sprintf("%s %s %s", m.Name, m.Version, m.Release)
	dest.Arg = m.Slug
	dest.Autocomplete = dest.Title
	return dest
}

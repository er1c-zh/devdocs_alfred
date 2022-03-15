package main

import (
	"fmt"
	"strings"
)

type AlfredResp struct {
	Items []ResultItem `json:"items"`
}

type ResultItem struct {
	Type         string     `json:"type"`
	Title        string     `json:"title"`
	Subtitle     string     `json:"subtitle"`
	Arg          string     `json:"arg"`
	Autocomplete string     `json:"autocomplete"`
	Icon         ResultIcon `json:"icon"`
}

type ResultIcon struct {
	Type string `json:"-"`
	Path string `json:"path"`
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
	dest.Title = fmt.Sprintf("%s_%s_%s", m.Name, m.Version, m.Release)
	dest.Subtitle = m.Name
	dest.Arg = m.Slug
	dest.Autocomplete = dest.Title
	dest.Icon = ResultIcon{
		Path: fmt.Sprintf("./icons/docs/%s_16@2x.png", strings.ToLower(m.Name)),
	}
	return dest
}

type DocIndex struct {
	Entries []DocIndexEntry `json:"entries"`
	Types   []DocIndexType  `json:"types"`
}

type DocIndexEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

func (e DocIndexEntry) ToAlfred(docSlug string) ResultItem {
	var dest ResultItem
	dest.Title = e.Name
	dest.Subtitle = fmt.Sprintf("分类: %s", e.Type)
	dest.Arg = fmt.Sprintf("%s/%s/%s", UrlBase, docSlug, e.Path)
	return dest
}

func (e DocIndexEntry) GetString() string {
	return e.Name
}

type DocIndexType struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
}

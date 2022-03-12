package main

import (
	"encoding/json"
	"fmt"
	"github.com/er1c-zh/go-now/log"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	UrlBase = "https://devdocs.io"
)

var (
	cache sync.Map
)

func cacheWrapper[T any](key string, loader func() (T, error), reload bool) (T, error) {
	if v, ok := cache.Load(key); ok && !reload {
		if typed, ok := v.(T); ok {
			return typed, nil
		}
	}
	//var v T
	v, err := loader()
	if err != nil {
		return v, err
	}
	cache.Store(key, v)
	return v, nil
}

// GetDocsList 获取支持的文档的列表
func GetDocsList() ([]RespDocMeta, error) {
	const (
		api = "/docs/docs.json"
	)
	return cacheWrapper[[]RespDocMeta](api, func() ([]RespDocMeta, error) {
		return loader[[]RespDocMeta](api)
	}, false)
}

func loader[T any](uri string) (T, error) {
	t0 := time.Now()
	defer func() {
		log.Trace("%s, cost: %dms", uri, time.Now().Sub(t0).Milliseconds())
	}()
	var t T
	url := buildUrl(uri)
	resp, err := http.Get(url)
	if err != nil {
		log.Error("%s http.Get fail: %s", url, err.Error())
		return t, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	byteList, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("%s ReadAll fail: %s", url, err.Error())
		return t, err
	}
	err = json.Unmarshal(byteList, &t)
	if err != nil {
		log.Error("%s Unmarshal fail: %s", url, err.Error())
		return t, err
	}
	return t, nil
}

func buildUrl(path string) string {
	return fmt.Sprintf("%s%s", UrlBase, path)
}

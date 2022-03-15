package main

import (
	"encoding/json"
	"github.com/er1c-zh/go-now/log"
	"testing"
)

func TestGetDocsList(t *testing.T) {
	defer log.Flush()
	v, err := GetDocsList()
	if err != nil {
		t.Error(err)
		return
	}
	j, _ := json.Marshal(v)
	t.Logf("%s", string(j))
}

func TestGetDocIndex(t *testing.T) {
	defer log.Flush()
	v, err := GetDocIndex("ruby~3")
	if err != nil {
		t.Error(err)
		return
	}
	j, _ := json.Marshal(v)
	t.Logf("%s", string(j))
}

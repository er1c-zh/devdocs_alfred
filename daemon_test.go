package main

import (
	"github.com/er1c-zh/go-now/log"
	"os"
	"strconv"
	"testing"
)

func TestIsProcessExist(t *testing.T) {
	defer log.Flush()
	type args struct {
		pidStr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"default",
			args{pidStr: "1220"},
			false,
		},
		{
			"default",
			args{pidStr: strconv.FormatInt(int64(os.Getpid()), 10)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsProcessExist(tt.args.pidStr); got != tt.want {
				t.Errorf("IsProcessExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

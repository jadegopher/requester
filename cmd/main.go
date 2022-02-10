package main

import (
	"context"
	"fmt"
	"time"

	"requester/config"
	"requester/getter"
	"requester/getter/md5"
	"requester/master"
)

const (
	defaultTimeout = time.Second * 10
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hashGetter := getter.NewGetter(defaultTimeout, md5.NewCalculator())
	mstr := master.NewMaster(cfg.GoroutinesCount, hashGetter)

	result := mstr.ProcessTasks(context.Background(), cfg.URLs)

	for k, v := range result {
		fmt.Println(k, v)
	}
}

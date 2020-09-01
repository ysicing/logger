// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package main

import (
	"github.com/ysicing/go-utils/extime"
	"github.com/ysicing/logger"
)

func init() {
	cfg := logger.LogConfig{Simple: false}
	logger.InitLogger(&cfg)
}

func main() {
	logger.Debug("debug")
	logger.Debugf("", "1", 2, 3, extime.GetToday())
	logger.Log.Sugar().Debug("1", 2, 3, extime.GetToday())
	logger.Info("info")
	logger.Error("error")
	logger.Exit("error exit")
}

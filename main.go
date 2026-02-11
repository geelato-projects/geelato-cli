package main

import (
	"os"

	"github.com/geelato/cli/cmd"
	"github.com/geelato/cli/pkg/logger"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	logger.Infof("Geelato CLI version %s", version)
	if err := cmd.Execute(); err != nil {
		logger.Errorf("执行失败: %v", err)
		os.Exit(1)
	}
}

package main

import (
    "github.com/qjpcpu/logger"
)

func main() {
    logger.NewLogBuilder("my-logger").File("/tmp/logfile").Level("info").Build()
    logger.LoggerOf("my-logger").Debug("hi debug")
    logger.LoggerOf("my-logger").Info("information")
    logger.NewLogBuilder("another-logger").File("/tmp/logfile2").Level("debug").Build()
    logger.LoggerOf("another-logger").Debug("debug info")
    logger.LoggerOf("another-logger").Info("information")
}

package main

import (
    "github.com/qjpcpu/logger"
)

func main() {
    //default logger
    logger.Debug("default debug")
    logger.Warning("default Warning")
    logger.NewLogBuilder("my-logger").Level("notice").Build()
    logger.LoggerOf("my-logger").Debug("hi debug")
    logger.LoggerOf("my-logger").Error("hi error")
}

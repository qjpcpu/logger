package main

import (
    "github.com/qjpcpu/logger"
    "time"
)

func main() {
    logger.NewLogBuilder("my-logger").File("/tmp/logfile").Rotate("YYYYMMDDHHmm", "* * * * *", "2m").Level("info").Build()
    go func() {
        for {
            time.Sleep(time.Second)
            logger.LoggerOf("my-logger").Debug("hi debug")
            logger.LoggerOf("my-logger").Info("information")
        }
    }()
    time.Sleep(time.Hour)
}

logger
======

### Simple usage

```
package main

import (
    "github.com/qjpcpu/logger"
)

func main() {
    //default logger with log level DEBUG
    logger.Debug("default debug")
    logger.Warning("default Warning")
    logger.NewLogBuilder("my-logger").Level("notice").Build()
    logger.LoggerOf("my-logger").Debug("hi debug")
    logger.LoggerOf("my-logger").Error("hi error")
}
```

### Log to file

If the file was rotated, logger would auto recreate log file.
```
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
```

### Log rotate

```
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
```

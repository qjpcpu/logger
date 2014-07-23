package logger

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
)

const (
    default_id = "default"
)
const (
    DEBUG   = 0
    INFO    = 1
    NOTICE  = 2
    WARNING = 3
    ERROR   = 4
    FATAL   = 5
)

var LEVELS = []string{"[DEBUG] ", "[INFO] ", "[NOTICE] ", "[WARNING] ", "[ERROR] ", "[FATAL] "}

type Logger struct {
    logger  *log.Logger
    level   int
    id      string
    rotator *Rotator
}

var loggers map[string]*Logger = make(map[string]*Logger)

type LogBuilder struct {
    id           string
    filename     string
    log_level    int
    flags        int
    time_format  string
    rotate_cycle string
    keep_time    string
}

func NewLogBuilder(nid string) *LogBuilder {
    if nid == "" || nid == default_id {
        LoggerOf(default_id).Fatal("Invalnid logger id:", nid)
    }
    if _, ok := loggers[nid]; ok {
        LoggerOf(default_id).Fatalf("%s has alreday exist.", nid)
    }
    return &LogBuilder{
        id:        nid,
        filename:  "",
        log_level: DEBUG,
        flags:     log.LstdFlags | log.Lshortfile,
    }
}
func (lb *LogBuilder) Level(l string) *LogBuilder {
    switch strings.ToUpper(l) {
    case "DEBUG":
        lb.log_level = DEBUG
    case "INFO":
        lb.log_level = INFO
    case "ERROR":
        lb.log_level = ERROR
    case "WARNING":
        lb.log_level = WARNING
    case "NOTICE":
        lb.log_level = NOTICE
    case "FATAL":
        lb.log_level = FATAL
    }
    return lb
}
func (lb *LogBuilder) File(fn string) *LogBuilder {
    lb.filename, _ = filepath.Abs(fn)
    return lb
}
func (lb *LogBuilder) Rotate(time_format, schedule, keep string) *LogBuilder {
    if lb.filename == "" {
        Fatal("Pls specify log file first")
    }
    lb.time_format = time_format
    lb.rotate_cycle = schedule
    lb.keep_time = keep
    return lb
}
func (lb *LogBuilder) Flags(flg int) *LogBuilder {
    lb.flags = flg
    return lb
}
func (lb *LogBuilder) Build() *Logger {
    if lb.filename == "" {
        lg := initLogger(lb.id, lb.log_level, lb.flags)
        loggers[lb.id] = lg
        return lg
    } else {
        lg := initFileLogger(lb.id, lb.filename, lb.log_level, lb.flags)
        loggers[lb.id] = lg
        if lb.time_format != "" && lb.rotate_cycle != "" {
            lg.rotator = NewRotator(lb.filename, lb.time_format, lb.rotate_cycle, lb.keep_time)
            lg.rotator.Start()
        }
        return lg
    }
}

// Get logger by id
func LoggerOf(id string) *Logger {
    if id == "" {
        id = default_id
    }
    if _, ok := loggers[default_id]; !ok {
        nlogger := &Logger{
            logger: log.New(os.Stdout, LEVELS[DEBUG], log.Lshortfile|log.LstdFlags),
            level:  DEBUG,
            id:     default_id,
        }
        loggers[default_id] = nlogger
    }
    if _, ok := loggers[id]; !ok {
        loggers[default_id].Fatalf("logger %s not exist.", id)
    }
    return loggers[id]
}

func initLogger(id string, level int, flags int) *Logger {
    nlogger := &Logger{
        logger: log.New(os.Stdout, LEVELS[level], flags),
        level:  level,
        id:     id,
    }
    loggers[id] = nlogger
    return nlogger
}
func initFileLogger(id, filename string, level int, flags int) *Logger {
    writer, _ := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
    nlogger := &Logger{
        logger: log.New(writer, LEVELS[level], flags),
        level:  level,
        id:     id,
    }
    loggers[id] = nlogger
    go func() {
        watcher, err := NewWatcher()
        if err != nil {
            return
        }
        defer watcher.Close()
        for {
            err = watcher.AddWatch(filepath.Dir(filename), IN_DELETE|IN_MOVE)
            if err != nil {
                fmt.Println(err)
                return
            }
            for {
                <-watcher.Event
                if _, err = os.Stat(filename); os.IsNotExist(err) {
                    break
                }
            }
            writer.Close()
            writer, _ = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
            loggers[id].logger = log.New(writer, LEVELS[level], flags)
        }
    }()
    return loggers[id]
}

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func (lg *Logger) Fatal(args ...interface{}) {
    if lg.level <= FATAL {
        lg.logger.SetPrefix(LEVELS[FATAL])
        lg.logger.Fatalln(args...)
    }
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func (lg *Logger) Fatalf(format string, args ...interface{}) {
    if lg.level <= FATAL {
        lg.logger.SetPrefix(LEVELS[FATAL])
        lg.logger.Fatalf(format, args...)
    }
}

// Warning logs a message using WARNING as log level.
func (lg *Logger) Warning(args ...interface{}) {
    if lg.level <= WARNING {
        lg.logger.SetPrefix(LEVELS[WARNING])
        lg.logger.Println(args...)
    }
}
func (lg *Logger) Warningf(format string, args ...interface{}) {
    if lg.level <= WARNING {
        lg.logger.SetPrefix(LEVELS[WARNING])
        lg.logger.Printf(format, args...)
    }
}

func (lg *Logger) Error(args ...interface{}) {
    if lg.level <= ERROR {
        lg.logger.SetPrefix(LEVELS[ERROR])
        lg.logger.Println(args...)
    }
}
func (lg *Logger) Errorf(format string, args ...interface{}) {
    if lg.level <= ERROR {
        lg.logger.SetPrefix(LEVELS[ERROR])
        lg.logger.Printf(format, args...)
    }
}

// Notice logs a message using INFO as log level.
func (lg *Logger) Notice(args ...interface{}) {
    if lg.level <= NOTICE {
        lg.logger.SetPrefix(LEVELS[NOTICE])
        lg.logger.Println(args...)
    }
}
func (lg *Logger) Noticef(format string, args ...interface{}) {
    if lg.level <= NOTICE {
        lg.logger.SetPrefix(LEVELS[NOTICE])
        lg.logger.Printf(format, args...)
    }
}

// Info logs a message using INFO as log level.
func (lg *Logger) Info(args ...interface{}) {
    if lg.level <= INFO {
        lg.logger.SetPrefix(LEVELS[INFO])
        lg.logger.Println(args...)
    }
}
func (lg *Logger) Infof(format string, args ...interface{}) {
    if lg.level >= INFO {
        lg.logger.SetPrefix(LEVELS[INFO])
        lg.logger.Printf(format, args...)
    }
}

// Debug logs a message using DEBUG as log level.
func (lg *Logger) Debug(args ...interface{}) {
    if lg.level <= DEBUG {
        lg.logger.SetPrefix(LEVELS[DEBUG])
        lg.logger.Println(args...)
    }
}
func (lg *Logger) Debugf(format string, args ...interface{}) {
    if lg.level <= DEBUG {
        lg.logger.SetPrefix(LEVELS[DEBUG])
        lg.logger.Printf(format, args...)
    }
}

// default logger
// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
    LoggerOf(default_id).Fatal(args...)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
    LoggerOf(default_id).Fatalf(format, args...)
}

// Error logs a message using INFO as log level.
func Error(args ...interface{}) {
    LoggerOf(default_id).Error(args...)
}
func Errorf(format string, args ...interface{}) {
    LoggerOf(default_id).Errorf(format, args...)
}

// Warning logs a message using WARNING as log level.
func Warning(args ...interface{}) {
    LoggerOf(default_id).Warning(args...)
}
func Warningf(format string, args ...interface{}) {
    LoggerOf(default_id).Warningf(format, args...)
}

// Notice logs a message using INFO as log level.
func Notice(args ...interface{}) {
    LoggerOf(default_id).Notice(args...)
}
func Noticef(format string, args ...interface{}) {
    LoggerOf(default_id).Noticef(format, args...)
}

// Info logs a message using INFO as log level.
func Info(args ...interface{}) {
    LoggerOf(default_id).Info(args...)
}
func Infof(format string, args ...interface{}) {
    LoggerOf(default_id).Infof(format, args...)
}

// Debug logs a message using DEBUG as log level.
func Debug(args ...interface{}) {
    LoggerOf(default_id).Debug(args...)
}
func Debugf(format string, args ...interface{}) {
    LoggerOf(default_id).Debugf(format, args...)
}

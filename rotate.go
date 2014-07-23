package logger

import (
    "github.com/Unknwon/com"
    "github.com/robfig/cron"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

type Rotator struct {
    filename string
    suffix   string
    cycle    string
    keeptime string
}

func NewRotator(fn, sf, c, k string) *Rotator {
    return &Rotator{
        filename: fn,
        suffix:   sf,
        cycle:    parseCycle(c),
        keeptime: k,
    }
}
func parseCycle(c string) string {
    if len(strings.Split(c, " ")) != 5 {
        Fatal("Cycle format error, the rotate cycle should look like crontab time schedule.")
    }
    return "0 " + c
}

//support d(day) h(hour) m(minute)
func parseKeeptime(raw string) string {
    raw = strings.ToLower(raw)
    crontab := []string{"0", "*", "*", "*", "*", "*"}
    if i := strings.Index(raw, "d"); i > 0 {
        num, err := strconv.Atoi(raw[0:i])
        if err != nil || num < 1 {
            Fatal("keeptime format err")
        }
        crontab[3] += "/" + raw[0:i]
        crontab[2], crontab[1] = "0", "0"
    } else if i := strings.Index(raw, "h"); i > 0 {
        num, err := strconv.Atoi(raw[0:i])
        if err != nil || num < 1 {
            Fatal("keeptime format err")
        }
        crontab[2] += "/" + raw[0:i]
        crontab[1] = "0"
    } else if i := strings.Index(raw, "m"); i > 0 {
        num, err := strconv.Atoi(raw[0:i])
        if err != nil || num < 1 {
            Fatal("keeptime format err")
        }
        crontab[1] += "/" + raw[0:i]
    }
    return strings.Join(crontab, " ")
}
func findOld(filename, format string, t time.Time) []string {
    limit := filename + "." + com.DateT(t, format)
    list, _ := filepath.Glob(filename + ".*")
    result := []string{}
    for _, fn := range list {
        if fn <= limit {
            result = append(result, fn)
        }
    }
    return result
}
func (r *Rotator) Start() {
    c := cron.New()
    // rotate
    c.AddFunc(r.cycle, func() {
        suffix := com.DateT(time.Now(), r.suffix)
        os.Rename(r.filename, r.filename+"."+suffix)
    })
    // remove old log
    c.AddFunc(r.cycle, func() {
        t := time.Now()
        if strings.Contains(r.keeptime, "d") {
            if num, err := strconv.Atoi(strings.TrimRight(r.keeptime, "d")); err == nil {
                t = t.AddDate(0, 0, -num)
            }
        } else if strings.Contains(r.keeptime, "h") {
            if num, err := strconv.Atoi(strings.TrimRight(r.keeptime, "h")); err == nil {
                t = t.Add(time.Duration(-num) * time.Hour)
            }
        } else if strings.Contains(r.keeptime, "m") {
            if num, err := strconv.Atoi(strings.TrimRight(r.keeptime, "m")); err == nil {
                t = t.Add(time.Duration(-num) * time.Minute)
            }
        } else {
            return
        }
        list := findOld(r.filename, r.suffix, t)
        for _, f := range list {
            os.Remove(f)
        }
    })
    c.Start()
}

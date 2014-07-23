package logger

import (
    "fmt"
    "testing"
    "time"
)

func TestParsekeeptime(t *testing.T) {
    str := parseKeeptime("1d")
    if str != "0 0 0 */1 * *" {
        t.Fatal("parse 1d err")
    }
    str = parseKeeptime("2h")
    if str != "0 0 */2 * * *" {
        t.Fatal("parse 2h err")
    }
    str = parseKeeptime("3d5m")
    if str != "0 0 0 */3 * *" {
        t.Fatal("parse 3d5m err")
    }
}
func TestRotate(t *testing.T) {
}

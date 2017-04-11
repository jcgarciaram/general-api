package apiutils

import (
    "fmt"
	"time"
    "strings"
)

type CustomTime struct {
	time.Time
}

const ctLayout = "January 2, 2006 3:04 PM"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
    s := strings.Trim(string(b), "\"")
    if s == "null" {
       ct.Time = time.Time{}
       return
    }
    ct.Time, err = time.Parse(ctLayout, s)
    return
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
  if ct.Time.UnixNano() == nilTime {
    return []byte("null"), nil
  }
  return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
}

var nilTime = (time.Time{}).UnixNano()
func (ct *CustomTime) IsSet() bool {
    return ct.UnixNano() != nilTime
}
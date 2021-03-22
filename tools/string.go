package tools

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}

func GetCurrentTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func StructToJsonStr(e interface{}) (string, error) {
	if b, err := json.Marshal(e); err == nil {
		return string(b), err
	} else {
		return "", err
	}
}

func Strip(s string) (r string) {
	if s != " " {
		s = strings.Trim(s, " ")
		s = strings.Trim(s, "\t")
		s = strings.Trim(s, "\n")
		r = strings.Trim(s, "\r")
	}
	return
}

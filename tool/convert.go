package tool

import (
	"fmt"
	"strconv"
)

func ToString(_v interface{}) string {
	return fmt.Sprintf("%v", _v)
}

func StringToInt64(_v string) int64 {
	n, _ := strconv.ParseInt(_v, 10, 64)
	return n
}

func StringToUint64(_v string) uint64 {
	n, _ := strconv.ParseUint(_v, 10, 64)
	return n
}

func StringToInt(_v string) int {
	n, _ := strconv.ParseInt(_v, 10, 64)
	return int(n)
}

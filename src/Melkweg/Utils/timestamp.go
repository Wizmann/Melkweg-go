package Util

import (
    "time"
)

func GetTimestamp() uint64 {
   return uint64(time.Now().UnixNano() / 1000000)
}

package utils

import (
	"fmt"
	"github.com/ghnexpress/traefik-ratelimit/log"
	"runtime"
)

func ShowErrorLogs(errData error) {
	if errData != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Log(fmt.Errorf("[%s][at line %d] %s", file, line, errData.Error()))
	}
}

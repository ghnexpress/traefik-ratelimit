package utils

import (
	"fmt"
	"log"
	"runtime"
)

func ShowErrorLogs(errData error) {
	if errData != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Println(fmt.Errorf("[%s][at line %d] %s", file, line, errData.Error()))
	}
}

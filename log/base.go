package log

import (
	"fmt"
	"os"
)

func Log(args ...any) {
	os.Stdout.WriteString(fmt.Sprintf("[rate-limit-middleware-plugin] %s\n", fmt.Sprint(args)))
}

type ErrorPublisher interface {
	SendError(err error) error
}

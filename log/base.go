package log

import (
	"fmt"
	"os"
)

func Log(text string) {
	os.Stdout.WriteString(fmt.Sprintf("rate-limit-middleware-plugin] %s\n", text))
}

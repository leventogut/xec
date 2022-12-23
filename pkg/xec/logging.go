package xec

import (
	"time"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()

}

// Info outputs informational messages.
func Info(message string, kv map[string]interface{}) {
	url := "test"
	sugar.Infow(message, "url", url, "attempt", 3, "backoff", time.Second)
	sugar.Infof("Failed to fetch URL: %s", url)
}

// // Warning outputs warning messages
// func Warning(message string) {
// 	severity := "warning"
// 	WriteLogLine(severity, message)
// }

// // WriteLogLine writes the given
// func WriteLogLine(severity string, message string, args ...interface{}) {
// 	// Date
// 	// Color
// 	dateTime := time.Date().UTC().Format("RFC3339")
// 	if severity == "warning" {
// 		color.Set(color.FgYellow)
// 		fmt.Printf("%s - %s", dateTime.String(), strings.ToUpper(severity))
// 		color.Unset()
// 	}

// }

// color.EnableColor()
// color.DisableColor()

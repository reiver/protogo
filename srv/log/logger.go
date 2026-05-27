package logsrv

import (
	"io"
	"os"

	"codeberg.org/reiver/go-log"

	"protogo/cfg"
)

var writer io.Writer = os.Stdout

var logger log.Logger = log.CreateLogger(writer, cfg.LogLevel())

// Logger return the logger.
//
// Example usage:
//
//	log := logsrv.Begin()
//	defer log.End()
//	
//	log.Inform(field.S("Hello world!"))
//
// See also:
//
//	• [Begin]
func Logger() log.Logger {
	return logger
}

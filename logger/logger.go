package logger

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

const timestampFormat = "2006-01-02 15:04:05.001 -0700 MST"

// ConfigureLogger init log settings
func ConfigureLogger(verbose bool) *logrus.Logger {
	var log = logrus.New()
	// default level is info
	if verbose {
		log.SetLevel(logrus.DebugLevel)
	}
	multiWriter := io.MultiWriter(os.Stderr, &lumberjack.Logger{
		Filename:   "logs/server.log", // Filename is the file to write logs to.  Backup log files will be retained in the same directory.
		MaxSize:    50,                // MaxSize is the maximum size in megabytes of the log file before it gets rotated
		MaxBackups: 5,                 // MaxBackups is the maximum number of old log files to retain.
		MaxAge:     30,                // MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
	})
	log.SetOutput(multiWriter)

	dateFormatter := &logrus.JSONFormatter{
		TimestampFormat: timestampFormat,
	}
	// output in JSON format
	log.SetFormatter(dateFormatter)

	return log
}

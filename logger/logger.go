package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func ConfigureLogger(verbose bool) *logrus.Logger {
	var log = logrus.New()
	if verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	formatter := new(logrus.JSONFormatter)
	formatter.TimestampFormat = "01-01-2001 13:00:00"

	log.SetFormatter(formatter)
	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetOutput(os.Stdout)

	return log
}

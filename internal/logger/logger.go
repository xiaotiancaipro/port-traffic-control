package logger

import (
	"fmt"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func New(path string, prefix string) (log *Log, err error) {
	_logger := logrus.New()
	_logger.SetFormatter(new(Formatter))
	_logger.SetReportCaller(true)
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s/%s-%%Y-%%m-%%d.log", path, prefix),
		rotatelogs.WithMaxAge(time.Hour*24*14),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		return
	}
	_logger.SetOutput(writer)
	_logger.SetLevel(logrus.InfoLevel)
	log = &Log{
		Logger: _logger,
	}
	return
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	format := fmt.Sprintf(
		"%s - %s - %s - %d - %s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		strings.ToUpper(entry.Level.String()),
		entry.Caller.Function,
		entry.Caller.Line,
		entry.Message,
	)
	return []byte(format), nil
}

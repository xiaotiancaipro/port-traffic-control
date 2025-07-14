package logger

import "github.com/sirupsen/logrus"

type Log struct {
	*logrus.Logger
}

type Formatter struct {
	logrus.TextFormatter
}

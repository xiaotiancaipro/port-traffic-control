package utils

import "port-traffic-control/internal/logger"

type Utils struct {
	ProcessUtil  *ProcessUtil
	StringUtil   *StringUtil
	ResponseUtil *ResponseUtil
}

type ProcessUtil struct {
	Log *logger.Log
}

type StringUtil struct{}

type ResponseUtil struct {
	Log *logger.Log
}

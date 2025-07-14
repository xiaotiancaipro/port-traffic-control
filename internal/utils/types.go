package utils

import "port-traffic-control/internal/logger"

type Utils struct {
	StringUtil   *StringUtil
	ResponseUtil *ResponseUtil
}

type StringUtil struct{}

type ResponseUtil struct {
	Log *logger.Log
}

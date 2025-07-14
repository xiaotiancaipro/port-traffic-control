package utils

import "port-traffic-control/internal/logger"

func New(log *logger.Log) *Utils {
	return &Utils{
		StringUtil: &StringUtil{},
		ResponseUtil: &ResponseUtil{
			Log: log,
		},
	}
}

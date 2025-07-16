package utils

import "port-traffic-control/internal/logger"

func New(log *logger.Log) *Utils {
	return &Utils{
		ProcessUtil: &ProcessUtil{
			Log: log,
		},
		StringUtil: &StringUtil{},
		ResponseUtil: &ResponseUtil{
			Log: log,
		},
	}
}

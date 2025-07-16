package srv

import (
	"net/http"
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/logger"
)

type Srv struct {
	Config    *configs.APIConfig
	Log       *logger.Log
	Server    *http.Server
	ServerErr chan error
}

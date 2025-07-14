package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (a *Auth) mount(c *gin.Context) {
	if err := a.allowedAuthorization(c.GetHeader("Authorization")); err != nil {
		a.ResponseUtil.Unauthorized(c, err.Error())
		c.Abort()
		return
	}
	c.Next()
	return
}

func (a *Auth) allowedAuthorization(requestHeader string) (err error) {
	if requestHeader == "" {
		err = fmt.Errorf("the request header does not have an Authorization field, RequestHeader=%s", requestHeader)
		a.Log.Error(err)
		return
	}
	var scheme, accessKey string
	_, err = fmt.Sscanf(requestHeader, "%s %s", &scheme, &accessKey)
	if (err != nil) || (scheme != "Bearer") {
		err = fmt.Errorf("the Authorization field in the request header is incorrectly formatted, RequestHeader=%s", requestHeader)
		a.Log.Error(err)
		return
	}
	if accessKey != a.Config.AccessKey {
		err = fmt.Errorf("access key error")
		a.Log.Error(err)
		return
	}
	return
}

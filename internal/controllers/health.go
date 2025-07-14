package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"port-traffic-control/internal/extensions"
	"port-traffic-control/internal/services"
	"strings"
)

func (hc *HealthController) Health(c *gin.Context) {

	var tables []string
	tx := hc.DB.
		Raw("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';").
		Pluck("tablename", &tables)
	if err := tx.Error; err != nil {
		hc.Log.Errorf("Unable to connect to database, Error=%v", err)
		hc.ResponseUtil.Error(c, services.InternalServerError)
		return
	}

	var notExistTables []string
	lookupMap := hc.StringUtil.SetupLookupMap(tables)
	for table := range extensions.Tables {
		if lookupMap[table] {
			continue
		}
		notExistTables = append(notExistTables, table)
	}
	if len(notExistTables) > 0 {
		hc.Log.Errorf("The data table does not exist, NotExistTables=\"%s\"", strings.Join(notExistTables, ", "))
		hc.ResponseUtil.Error(c, services.InternalServerError)
		return
	}

	info := fmt.Sprintf("Service is running")
	hc.Log.Info(info)
	hc.ResponseUtil.Success(c, info, nil)
	return

}

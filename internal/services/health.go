package services

import (
	"fmt"
	"port-traffic-control/internal/extensions"
	"strings"
)

func (hs *HealthService) Check() error {

	var tables []string
	tx := hs.DB.
		Raw("SELECT name FROM 'sqlite_master' WHERE type='table' ORDER BY name;").
		Pluck("tablename", &tables)
	if err := tx.Error; err != nil {
		err = fmt.Errorf("unable to connect to database, Error=%v", err)
		hs.Log.Error(err)
		return err
	}

	var notExistTables []string
	lookupMap := hs.StringUtil.SetupLookupMap(tables)
	for table := range extensions.Tables {
		if lookupMap[table] {
			continue
		}
		notExistTables = append(notExistTables, table)
	}
	if len(notExistTables) > 0 {
		err := fmt.Errorf("the data table does not exist, NotExistTables=\"%s\"", strings.Join(notExistTables, ", "))
		hs.Log.Error(err)
		return err
	}

	return nil

}

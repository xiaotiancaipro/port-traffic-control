package extensions

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"port-traffic-control/internal/configs"
	"time"
)

var Tables = map[string]any{}

func NewDB(config *configs.DatabaseConfig) (db *gorm.DB, err error) {

	db, err = gorm.Open(sqlite.Open(config.Path), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("database connection failed, Error=%v", err)
		return
	}

	for name, table := range Tables {
		if err_ := db.AutoMigrate(&table); err_ != nil {
			err = fmt.Errorf("table migration failed, TableName=%s, Error=%v", name, err_)
			return
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		err = fmt.Errorf("failed to get the underlying SQL.DB, Error=%v", err)
		return
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return

}

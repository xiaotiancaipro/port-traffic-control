package extensions

import (
	"fmt"
	"path/filepath"
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Tables = map[string]any{
	"groups": models.Groups{},
	"ports":  models.Ports{},
}

func NewDB(config *configs.DatabaseConfig) (db *gorm.DB, err error) {

	file := filepath.Join(config.Path, config.File)
	db, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
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

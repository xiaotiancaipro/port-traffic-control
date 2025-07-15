package services

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"port-traffic-control/internal/models"
)

func (ps *PortsService) Create(groupID uuid.UUID, port int32) error {
	ports := models.Ports{
		GroupID: groupID,
		Port:    port,
	}
	err := ps.DB.Transaction(func(tx *gorm.DB) error {
		tx_ := tx.
			Model(&models.Ports{}).
			Create(&ports)
		if err_ := tx_.Error; err_ != nil {
			return fmt.Errorf("failed to insert data, Errors=%v", err_)
		}
		return nil
	})
	if err != nil {
		ps.Log.Error(err)
		return err
	}
	return nil
}

func (ps *PortsService) Delete(ports models.Ports) error {
	err := ps.DB.Transaction(func(tx *gorm.DB) error {
		tx_ := tx.
			Model(&models.Ports{}).
			Where(&models.Ports{
				ID: ports.ID,
			}).
			Select("Status").
			Updates(&models.Ports{
				Status: 0,
			})
		if err_ := tx_.Error; err_ != nil {
			return fmt.Errorf("failed to delete data, Errors=%v", err_)
		}
		return nil
	})
	if err != nil {
		ps.Log.Error(err)
		return err
	}
	return nil
}

package services

import (
	"fmt"
	"gorm.io/gorm"
	"port-traffic-control/internal/models"
)

func (gs *GroupsService) Create(bandwidth int32, portMaxNum int32) error {
	newGroups := models.Groups{
		Bandwidth:  bandwidth,
		PortMaxNum: portMaxNum,
	}
	err := gs.DB.Transaction(func(tx *gorm.DB) error {
		tx_ := tx.
			Model(&models.Groups{}).
			Create(&newGroups)
		if err_ := tx_.Error; err_ != nil {
			return fmt.Errorf("failed to insert data, Errors=%v", err_)
		}
		return nil
	})
	if err != nil {
		gs.Log.Error(err)
		return err
	}
	return nil
}

func (gs *GroupsService) Delete(groups models.Groups) error {
	err := gs.DB.Transaction(func(tx *gorm.DB) error {
		tx_ := tx.
			Model(&models.Groups{}).
			Where(&models.Groups{
				ID: groups.ID,
			}).
			Select("Status").
			Updates(&models.Groups{
				Status: 0,
			})
		if err_ := tx_.Error; err_ != nil {
			return fmt.Errorf("failed to delete data, Errors=%v", err_)
		}
		return nil
	})
	if err != nil {
		gs.Log.Error(err)
		return err
	}
	return nil
}

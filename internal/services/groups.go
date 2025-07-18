package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"port-traffic-control/internal/models"
)

var GroupNotFoundError = fmt.Errorf("group not found")

func (gs *GroupsService) IsNotExists(err error) error {
	if errors.Is(err, GroupNotFoundError) {
		return nil
	}
	return err
}

func (gs *GroupsService) Create(bandwidth int32, portMaxNum int32) (groups models.Groups, err error) {

	gs.Lock.Lock()
	defer gs.Lock.Unlock()

	newHandle := -1
	var handleMaxGroups models.Groups
	tx := gs.DB.
		Model(&models.Groups{}).
		Order("created_at desc").
		First(&handleMaxGroups)
	if err = tx.Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("database query failed, Error=%v", err)
			gs.Log.Error(err)
			return
		}
		newHandle = 0
	}
	if newHandle < 0 {
		newHandle = int(handleMaxGroups.Handle + 1)
	}

	groups = models.Groups{
		Handle:     int32(newHandle),
		Bandwidth:  bandwidth,
		PortMaxNum: portMaxNum,
	}
	err = gs.DB.Transaction(func(tx *gorm.DB) error {
		tx_ := tx.
			Model(&models.Groups{}).
			Create(&groups)
		if err_ := tx_.Error; err_ != nil {
			return fmt.Errorf("failed to insert data, Errors=%v", err_)
		}
		return nil
	})
	if err != nil {
		gs.Log.Error(err)
		return
	}

	return

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

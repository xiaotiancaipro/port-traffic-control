package services

import (
	"errors"
	"fmt"
	"port-traffic-control/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	var handleMaxGroups *models.Groups
	tx := gs.DB.
		Model(&models.Groups{}).
		Where("status != ?", 0).
		Order("handle desc").
		First(&handleMaxGroups)
	if err = tx.Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("database query failed, Error=%v", err)
			gs.Log.Error(err)
			return
		}
	}

	newHandle := int32(1)
	if handleMaxGroups != nil {
		newHandle = handleMaxGroups.Handle + 1
	}

	groups = models.Groups{
		Handle:     newHandle,
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

func (gs *GroupsService) UpdateUsage(groups models.Groups, usage int32) error {
	err := gs.DB.Transaction(func(tx *gorm.DB) error {
		tx_ := tx.
			Model(&models.Groups{}).
			Where(&models.Groups{
				ID: groups.ID,
			}).
			Select("usage").
			Updates(&models.Groups{
				Usage: usage,
			})
		if err_ := tx_.Error; err_ != nil {
			return fmt.Errorf("failed to update data, Errors=%v", err_)
		}
		return nil
	})
	if err != nil {
		gs.Log.Error(err)
		return err
	}
	return nil
}

func (gs *GroupsService) UpdateStatus(groups models.Groups, status int8) error {
	err := gs.DB.Transaction(func(tx *gorm.DB) error {
		tx_ := tx.
			Model(&models.Groups{}).
			Where(&models.Groups{
				ID: groups.ID,
			}).
			Select("status").
			Updates(&models.Groups{
				Status: status,
			})
		if err_ := tx_.Error; err_ != nil {
			return fmt.Errorf("failed to update data, Errors=%v", err_)
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
	return gs.UpdateStatus(groups, 0)
}

func (gs *GroupsService) GetByID(id uuid.UUID) (groups models.Groups, err error) {
	tx := gs.DB.
		Model(&models.Groups{}).
		Where(&models.Groups{
			ID:     id,
			Status: 1,
		}).
		First(&groups)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			err = GroupNotFoundError
			return
		}
		err = fmt.Errorf("database query failed, Error=%v", tx.Error)
		gs.Log.Error(err)
		return
	}
	return
}

func (gs *GroupsService) ListAll() (groups []models.Groups, err error) {
	tx := gs.DB.
		Model(&models.Groups{}).
		Where(&models.Groups{
			Status: 1,
		}).
		Find(&groups)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			err = GroupNotFoundError
			return
		}
		err = fmt.Errorf("database query failed, Error=%v", tx.Error)
		gs.Log.Error(err)
		return
	}
	return
}

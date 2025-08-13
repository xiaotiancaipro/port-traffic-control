package services

import (
	"errors"
	"fmt"
	"port-traffic-control/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var PortNotFoundError = fmt.Errorf("port not found")

func (ps *PortsService) IsNotExists(err error) error {
	if errors.Is(err, PortNotFoundError) {
		return nil
	}
	return err
}

func (ps *PortsService) Create(groupID uuid.UUID, port int32) (ports models.Ports, err error) {
	ports = models.Ports{
		GroupID: groupID,
		Port:    port,
	}
	err = ps.DB.Transaction(func(tx *gorm.DB) error {
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
		return
	}
	return
}

func (ps *PortsService) UpdateStatus(ports models.Ports, status int8) error {
	err := ps.DB.Transaction(func(tx *gorm.DB) error {
		tx_ := tx.
			Model(&models.Ports{}).
			Where(&models.Ports{
				ID: ports.ID,
			}).
			Select("status").
			Updates(&models.Ports{
				Status: status,
			})
		if err_ := tx_.Error; err_ != nil {
			return fmt.Errorf("failed to update data, Errors=%v", err_)
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
	return ps.UpdateStatus(ports, 0)
}

func (ps *PortsService) DeleteByID(id uuid.UUID) error {
	return ps.Delete(models.Ports{ID: id})
}

func (ps *PortsService) ListActivePorts() (ports []models.Ports, err error) {
	tx := ps.DB.
		Model(&models.Ports{}).
		Where(&models.Ports{
			Status: 1,
		}).
		Order("created_at asc").
		Find(&ports)
	if tx.Error != nil {
		err = fmt.Errorf("database query failed, Errors=%v", tx.Error)
		ps.Log.Error(err)
		return
	}
	return
}

func (ps *PortsService) ListActivePortsByGroupID(groupID uuid.UUID) (ports []models.Ports, err error) {
	tx := ps.DB.
		Model(&models.Ports{}).
		Where(&models.Ports{
			GroupID: groupID,
			Status:  1,
		}).
		Order("created_at asc").
		Find(&ports)
	if err_ := tx.Error; err_ != nil {
		if errors.Is(gorm.ErrRecordNotFound, err_) {
			err = PortNotFoundError
			return
		}
		err = fmt.Errorf("database query failed, Errors=%v", tx.Error)
		ps.Log.Error(err)
		return
	}
	return
}

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

func (ps *PortsService) CreateInTx(tx *gorm.DB, groupID uuid.UUID, port int32) error {
	op := tx.
		Model(&models.Ports{}).
		Create(&models.Ports{
			GroupID: groupID,
			Port:    port,
		})
	if err_ := op.Error; err_ != nil {
		err_ = fmt.Errorf("failed to insert data, Errors=%v", err_)
		ps.Log.Error(err_)
		return err_
	}
	return nil
}

func (ps *PortsService) UpdateStatusInTx(tx *gorm.DB, ports models.Ports, status int8) error {
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
		err_ = fmt.Errorf("failed to update data, Errors=%v", err_)
		ps.Log.Error(err_)
		return err_
	}
	return nil
}

func (ps *PortsService) DeleteInTx(tx *gorm.DB, ports models.Ports) error {
	return ps.UpdateStatusInTx(tx, ports, 0)
}

func (ps *PortsService) DeleteByIDInTx(tx *gorm.DB, id uuid.UUID) error {
	return ps.DeleteInTx(tx, models.Ports{ID: id})
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

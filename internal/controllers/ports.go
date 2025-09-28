package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (pc *PortsController) Add(c *gin.Context) {

	request := RequestBodyPorts{}
	if !pc.ResponseUtil.ParsingRequest(c, &request) {
		return
	}

	groupUUID, err := uuid.Parse(request.GroupID)
	if err != nil {
		pc.ResponseUtil.BadRequest(c, "Invalid groupID")
		return
	}

	if len(request.PortList) == 0 {
		pc.ResponseUtil.BadRequest(c, "portList is empty")
		return
	}

	group, err := pc.GroupsService.GetByID(groupUUID)
	if err != nil {
		if err_ := pc.GroupsService.IsNotExists(err); err_ != nil {
			pc.Log.Errorf("Error getting group")
			pc.ResponseUtil.Error(c, "Error getting group")
			return
		}
		pc.Log.Errorf("Group not found")
		pc.ResponseUtil.Error(c, "Group not found")
		return
	}

	exists, err := pc.PortsService.ListActivePorts()
	if err != nil {
		if err_ := pc.PortsService.IsNotExists(err); err_ != nil {
			pc.Log.Errorf("Error getting ports")
			pc.ResponseUtil.Error(c, "Error getting ports")
			return
		}
	}

	existsMap := make(map[int32]struct{}, len(exists))
	for _, port := range exists {
		existsMap[port.Port] = struct{}{}
	}

	num := group.PortMaxNum - group.Usage
	usage := group.Usage
	successful := make([]int32, 0)
	failed := make([]ResponseBodyPortsFailedItem, 0)
	failedItem := func(port_ int32, err_ string) ResponseBodyPortsFailedItem {
		return ResponseBodyPortsFailedItem{
			Port:  port_,
			Error: err_,
		}
	}
	for _, port := range request.PortList {

		if num <= 0 {
			failed = append(failed, failedItem(port, "Available quantity has been exhausted"))
			continue
		}

		if _, ok := existsMap[port]; ok {
			failed = append(failed, failedItem(port, "The port already exists"))
			continue
		}

		if err = pc.DB.Transaction(func(tx *gorm.DB) error {

			if err_ := pc.PortsService.CreateInTx(tx, group.ID, port); err_ != nil {
				pc.Log.Errorf("Error adding port, Port=%d, Error=%v", port, err_)
				return err_
			}

			usage++
			if err_ := pc.GroupsService.UpdateUsageInTX(tx, group, usage); err_ != nil {
				usage--
				pc.Log.Errorf("Error update usage, Port=%d, Error=%v", port, err_)
				return err_
			}

			parentMinor := uint32(group.Handle)
			childMinor := uint32(port)
			rateCeil := uint32(group.Bandwidth)
			if err_ := pc.TCService.CreateChildClass(parentMinor, childMinor, 1, rateCeil); err_ != nil {
				usage--
				pc.Log.Errorf("Error create tc child class, Port=%d, Error=%v", port, err_)
				return err_
			}

			successful = append(successful, port)
			existsMap[port] = struct{}{}
			num--
			return nil

		}); err != nil {
			failed = append(failed, failedItem(port, err.Error()))
		}

	}

	pc.ResponseUtil.Success(c, "Successfully", ResponseBodyPorts{
		Successful: successful,
		Failed:     failed,
	})
	return

}

func (pc *PortsController) Remove(c *gin.Context) {

	request := RequestBodyPorts{}
	if !pc.ResponseUtil.ParsingRequest(c, &request) {
		return
	}

	groupUUID, err := uuid.Parse(request.GroupID)
	if err != nil {
		pc.ResponseUtil.BadRequest(c, "Invalid groupID")
		return
	}

	if len(request.PortList) == 0 {
		pc.ResponseUtil.BadRequest(c, "portList is empty")
		return
	}

	group, err := pc.GroupsService.GetByID(groupUUID)
	if err != nil {
		if err_ := pc.GroupsService.IsNotExists(err); err_ != nil {
			pc.Log.Errorf("Error getting group")
			pc.ResponseUtil.Error(c, "Error getting group")
			return
		}
		pc.Log.Errorf("Group not found")
		pc.ResponseUtil.Error(c, "Group not found")
		return
	}

	exists, err := pc.PortsService.ListActivePortsByGroupID(group.ID)
	if err != nil {
		if err_ := pc.PortsService.IsNotExists(err); err_ != nil {
			pc.Log.Errorf("Error getting ports")
			pc.ResponseUtil.Error(c, "Error getting ports")
			return
		}
	}

	existsMap := make(map[int32]uuid.UUID, len(exists))
	for _, port := range exists {
		existsMap[port.Port] = port.ID
	}

	usage := group.Usage
	successful := make([]int32, 0)
	failed := make([]ResponseBodyPortsFailedItem, 0)
	failedItem := func(port_ int32, err_ string) ResponseBodyPortsFailedItem {
		return ResponseBodyPortsFailedItem{
			Port:  port_,
			Error: err_,
		}
	}
	for _, port := range request.PortList {

		portID, ok := existsMap[port]
		if !(ok && portID != uuid.Nil) {
			failed = append(failed, failedItem(port, "The port not exists"))
			continue
		}

		if err = pc.DB.Transaction(func(tx *gorm.DB) error {

			if err_ := pc.PortsService.DeleteByIDInTx(tx, portID); err_ != nil {
				pc.Log.Errorf("Error remove port, Port=%d", port)
				return err_
			}

			usage--
			if err_ := pc.GroupsService.UpdateUsageInTX(tx, group, usage); err_ != nil {
				usage++
				pc.Log.Errorf("Error update usage, Port=%d", port)
				return err_
			}

			parentMinor := uint32(group.Handle)
			childMinor := uint32(port)
			if err_ := pc.TCService.DeleteChildClass(parentMinor, childMinor); err_ != nil {
				pc.Log.Errorf("Error delete tc child class, Port=%d", port)
				return err_
			}

			successful = append(successful, port)
			existsMap[port] = uuid.Nil
			return nil

		}); err != nil {
			failed = append(failed, failedItem(port, err.Error()))
		}

	}

	pc.ResponseUtil.Success(c, "Successfully", ResponseBodyPorts{
		Successful: successful,
		Failed:     failed,
	})

}

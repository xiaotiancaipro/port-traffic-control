package controllers

import (
	"port-traffic-control/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			pc.ResponseUtil.Error(c, "Error getting group")
			return
		}
		pc.ResponseUtil.Error(c, "Group not found")
		return
	}

	exists, err := pc.PortsService.ListActivePorts()
	if err != nil {
		if err_ := pc.PortsService.IsNotExists(err); err_ != nil {
			pc.ResponseUtil.Error(c, "Error getting ports")
			return
		}
	}

	existsMap := make(map[int32]struct{}, len(exists))
	for _, port := range exists {
		existsMap[port.Port] = struct{}{}
	}

	var (
		num        = group.PortMaxNum - group.Usage
		usage      = group.Usage
		successful []int32
		failed     []ResponseBodyPortsFailedItem
	)
	for _, port := range request.PortList {
		if num <= 0 {
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: "Available quantity has been exhausted",
			})
			continue
		}
		if _, ok := existsMap[port]; ok {
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: "The port already exists",
			})
			continue
		}
		ports, err_ := pc.PortsService.Create(group.ID, port)
		if err_ != nil {
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: err_.Error(),
			})
			pc.Log.Errorf("Error adding port, Port=%d", port)
			continue
		}
		// TODO tc add child class
		if err_ = pc.PortsService.UpdateStatus(ports, 1); err_ != nil {
			// TODO tc remove child class
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: err_.Error(),
			})
			pc.Log.Errorf("Error update port status, Port=%d", port)
			continue
		}
		usage++
		if err_ = pc.GroupsService.UpdateUsage(group, usage); err_ != nil {
			usage--
			// TODO tc remove child class
			if err__ := pc.PortsService.UpdateStatus(ports, 0); err__ != nil {
				pc.Log.Errorf("Error[Ser] update port status, Port=%d", port)
			}
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: err_.Error(),
			})
			pc.Log.Errorf("Error update usage, Port=%d", port)
			continue
		}
		successful = append(successful, port)
		existsMap[port] = struct{}{}
		num--
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
			pc.ResponseUtil.Error(c, "Error getting group")
			return
		}
		pc.ResponseUtil.Error(c, "Group not found")
		return
	}

	exists, err := pc.PortsService.ListActivePortsByGroupID(group.ID)
	if err != nil {
		if err_ := pc.PortsService.IsNotExists(err); err_ != nil {
			pc.ResponseUtil.Error(c, "Error getting ports")
			return
		}
	}

	existsMap := make(map[int32]uuid.UUID, len(exists))
	for _, port := range exists {
		existsMap[port.Port] = port.ID
	}

	var (
		usage      = group.Usage
		successful []int32
		failed     []ResponseBodyPortsFailedItem
	)
	for _, port := range request.PortList {
		portID, ok := existsMap[port]
		if !(ok && portID != uuid.Nil) {
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: "The port not exists",
			})
			continue
		}
		// TODO tc remove child class
		if err_ := pc.PortsService.DeleteByID(portID); err_ != nil {
			// TODO tc add child class
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: err_.Error(),
			})
			pc.Log.Errorf("Error remove port, Port=%d", port)
			continue
		}
		usage--
		if err_ := pc.GroupsService.UpdateUsage(group, usage); err_ != nil {
			usage++
			// TODO tc add child class
			if err__ := pc.PortsService.UpdateStatus(models.Ports{ID: portID}, 1); err__ != nil {
				pc.Log.Errorf("Error[Ser] activate port, Port=%d", port)
			}
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: err_.Error(),
			})
			pc.Log.Errorf("Error update usage, Port=%d", port)
			continue
		}
		successful = append(successful, port)
		existsMap[port] = uuid.Nil
	}

	pc.ResponseUtil.Success(c, "Successfully", ResponseBodyPorts{
		Successful: successful,
		Failed:     failed,
	})
	return

}

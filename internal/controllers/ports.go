package controllers

import (
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

	activePorts, err := pc.PortsService.ListActivePorts()
	if err != nil {
		if err_ := pc.PortsService.IsNotExists(err); err_ != nil {
			pc.ResponseUtil.Error(c, "Error getting ports")
			return
		}
	}

	activePortsSet := make(map[int32]struct{}, len(activePorts))
	for _, port := range activePorts {
		activePortsSet[port.Port] = struct{}{}
	}

	var addList []int32
	for _, port := range request.PortList {
		if _, ok := activePortsSet[port]; ok {
			continue
		}
		addList = append(addList, port)
	}

	var successfulList, failedList []int32
	for _, port := range addList {
		ports, err_ := pc.PortsService.Create(group.ID, port)
		if err_ != nil {
			failedList = append(failedList, port)
			pc.Log.Errorf("Error adding port %d", port)
			continue
		}
		// TODO tc add child class
		if err_ := pc.PortsService.UpdateStatus(ports, 1); err_ != nil {
			return
		}
		successfulList = append(successfulList, port)
	}

	if len(failedList) > 0 {
		pc.ResponseUtil.Success(c, "Partial success", ResponseBodyPorts{
			Flag:           0,
			SuccessfulList: successfulList,
			FailedList:     failedList,
		})
		return
	}

	pc.ResponseUtil.Success(c, "Add ports successfully", ResponseBodyPorts{
		Flag:           1,
		SuccessfulList: successfulList,
		FailedList:     nil,
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

	activePorts, err := pc.PortsService.ListActivePortsByGroupID(group.ID)
	if err != nil {
		if err_ := pc.PortsService.IsNotExists(err); err_ != nil {
			pc.ResponseUtil.Error(c, "Error getting ports")
			return
		}
	}

	portToRecord := make(map[int32]uuid.UUID, len(activePorts))
	for _, rec := range activePorts {
		portToRecord[rec.Port] = rec.ID
	}

	var successfulList, failedList []int32
	for _, port := range request.PortList {
		id, ok := portToRecord[port]
		if !ok {
			continue
		}
		if err_ := pc.PortsService.DeleteByID(id); err_ != nil {
			failedList = append(failedList, port)
			pc.Log.Errorf("Error remove port %d", port)
			continue
		}
		// TODO tc remove child class
		successfulList = append(successfulList, port)
	}

	if len(failedList) > 0 {
		pc.ResponseUtil.Success(c, "Partial success", ResponseBodyPorts{
			Flag:           0,
			SuccessfulList: successfulList,
			FailedList:     failedList,
		})
		return
	}

	pc.ResponseUtil.Success(c, "Add ports successfully", ResponseBodyPorts{
		Flag:           1,
		SuccessfulList: successfulList,
		FailedList:     nil,
	})
	return

}

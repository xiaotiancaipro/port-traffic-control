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
		num        int32 = group.PortMaxNum - group.Usage
		usage      int32 = group.Usage
		fixUsage   int32 = 0
		successful []int32
		failed     []ResponseBodyPortsFailedItem
	)
	for _, port := range request.PortList {
		if _, ok := existsMap[port]; ok {
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: "The port already exists",
			})
			continue
		}
		if num <= 0 {
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: "The maximum number of ports that can be accommodated by the group has been exceeded",
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
		usage++
		if err_ = pc.GroupsService.UpdateUsage(group, usage); err_ != nil {
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: err_.Error(),
			})
			pc.Log.Errorf("Error update usage, Port=%d", port)
			continue
		}
		if err_ = pc.PortsService.UpdateStatus(ports, 1); err_ != nil {
			fixUsage++
			failed = append(failed, ResponseBodyPortsFailedItem{
				Port:  port,
				Error: err_.Error(),
			})
			pc.Log.Errorf("Error update port status, Port=%d", port)
			continue
		}
		successful = append(successful, port)
		existsMap[port] = struct{}{}
		num--
	}

	if fixUsage > 0 {
		if err_ := pc.GroupsService.UpdateUsage(group, usage-fixUsage); err_ != nil {
			pc.Log.Errorf(
				"The usage correction failed and needs to be corrected manually, GroupID=%s, UsageFixd=%d",
				group.ID,
				group.Usage-fixUsage,
			)
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
		pc.ResponseUtil.Success(c, "Partial success", ResponseBodyPortsOld{
			Flag:           0,
			SuccessfulList: successfulList,
			FailedList:     failedList,
		})
		return
	}

	pc.ResponseUtil.Success(c, "Add ports successfully", ResponseBodyPortsOld{
		Flag:           1,
		SuccessfulList: successfulList,
		FailedList:     nil,
	})
	return

}

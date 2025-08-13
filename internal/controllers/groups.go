package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (gc *GroupsController) Create(c *gin.Context) {

	request := RequestBodyGroupsCreate{}
	if !gc.ResponseUtil.ParsingRequest(c, &request) {
		return
	}

	groups, err := gc.GroupsService.Create(request.Bandwidth, request.PortMaxNum)
	if err != nil {
		gc.Log.Error("Error creating groups")
		gc.ResponseUtil.Error(c, "Error creating groups")
		return
	}

	err = gc.TCService.CreateParentClass(uint32(groups.Handle), uint32(request.Bandwidth))
	if err != nil {
		gc.Log.Error("Error creating tc parent class")
		gc.ResponseUtil.Error(c, "Error creating tc parent class")
		return
	}

	err = gc.GroupsService.UpdateStatus(groups, 1)
	if err != nil {
		gc.Log.Error("Error updating status")
		gc.ResponseUtil.Error(c, "Error updating status")
		return
	}

	gc.Log.Infof("Create groups successfully, GroupsID=%s", groups.ID)
	gc.ResponseUtil.Success(c, "Create groups successfully", ResponseBodyGroupsCreate{GroupID: groups.ID})
	return

}

func (gc *GroupsController) Get(c *gin.Context) {

	request := RequestBodyGroupsGet{}
	if !gc.ResponseUtil.ParsingRequest(c, &request) {
		return
	}

	groupUUID, err := uuid.Parse(request.GroupID)
	if err != nil {
		gc.ResponseUtil.BadRequest(c, "Invalid groupID")
		return
	}

	groups, err := gc.GroupsService.GetByID(groupUUID)
	if err != nil {
		if err_ := gc.GroupsService.IsNotExists(err); err_ != nil {
			gc.ResponseUtil.Error(c, "Error getting group")
			return
		}
		gc.ResponseUtil.Error(c, "Group not found")
		return
	}

	ports, err := gc.PortsService.ListActivePortsByGroupID(groups.ID)
	if err != nil {
		if err_ := gc.PortsService.IsNotExists(err); err_ != nil {
			gc.ResponseUtil.Error(c, "Error getting ports")
			return
		}
	}

	portNums := make([]int32, 0, len(ports))
	for _, p := range ports {
		portNums = append(portNums, p.Port)
	}

	gc.ResponseUtil.Success(c, "Get group successfully", ResponseBodyGroupsGet{
		Bandwidth:  groups.Bandwidth,
		PortMaxNum: groups.PortMaxNum,
		PortList:   portNums,
	})
	return

}

func (gc *GroupsController) Delete(c *gin.Context) {

	request := RequestBodyGroupsDelete{}
	if !gc.ResponseUtil.ParsingRequest(c, &request) {
		return
	}

	groupUUID, err := uuid.Parse(request.GroupID)
	if err != nil {
		gc.ResponseUtil.BadRequest(c, "Invalid groupID")
		return
	}

	groups, err := gc.GroupsService.GetByID(groupUUID)
	if err != nil {
		if err_ := gc.GroupsService.IsNotExists(err); err_ != nil {
			gc.ResponseUtil.Error(c, "Error getting group")
			return
		}
		gc.ResponseUtil.Error(c, "Group not found")
		return
	}

	if err := gc.GroupsService.Delete(groups); err != nil {
		gc.Log.Errorf("Error deleting group, GroupID=%s", groups.ID)
		gc.ResponseUtil.Error(c, "Error deleting group")
		return
	}

	gc.Log.Infof("Delete groups successfully, GroupsID=%s", groups.ID)
	gc.ResponseUtil.Success(c, "Delete groups successfully", nil)
	return

}

func (gc *GroupsController) List(c *gin.Context) {

	groups, err := gc.GroupsService.ListAll()
	if err != nil {
		if err_ := gc.GroupsService.IsNotExists(err); err_ != nil {
			gc.ResponseUtil.Error(c, "Error getting group")
			return
		}
		gc.ResponseUtil.Success(c, "List groups successfully", ResponseBodyGroupsList{Groups: nil})
		return
	}

	groupsList := make([]uuid.UUID, 0, len(groups))
	for _, g := range groups {
		groupsList = append(groupsList, g.ID)
	}
	gc.ResponseUtil.Success(c, "List groups successfully", ResponseBodyGroupsList{Groups: groupsList})
	return

}

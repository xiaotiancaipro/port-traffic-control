package controllers

import "github.com/gin-gonic/gin"

func (gc *GroupsController) Create(c *gin.Context) {

	type requestBody struct {
		Bandwidth  int32 `json:"bandwidth"`
		PortMaxNum int32 `json:"portMaxNum"`
	}

	type responseBody struct {
		Flag int8 `json:"flag"`
	}

	request := requestBody{}
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

	gc.Log.Infof("Create groups successfully, GroupsID=%s", groups.ID)
	gc.ResponseUtil.Success(c, "Create groups successfully", responseBody{Flag: 1})
	return

}

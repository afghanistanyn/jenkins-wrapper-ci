package jenkins_manage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"jenkins-wrapper-ci/global"
	"jenkins-wrapper-ci/model/common/response"
	"jenkins-wrapper-ci/model/jenkins_manage"
	jenkins_manageReq "jenkins-wrapper-ci/model/jenkins_manage/request"
	"jenkins-wrapper-ci/utils"
)

type AppApi struct {
}




// CreateApp 创建App
//	@Tags		App
//	@Summary	创建App
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.App	true	"创建App"
//	@Success	200		{string}	string				"{"success":true,"data":{},"msg":"获取成功"}"
//	@Router		/app/createApp [post]
func (appApi *AppApi) CreateApp(c *gin.Context) {
	var app jenkins_manage.App
	err := c.ShouldBindJSON(&app)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    app.CreatedBy = utils.GetUserID(c)

	if (app.CustomConfig == "" && app.BuildParams == nil) {
		app.BuildParams = jenkins_manage.DefaultBuildParams
	}


	var verify utils.Rules
	if(app.CustomConfig != "") {
		verify = utils.Rules{
			"ProjectId": {utils.NotEmpty()},
			"AppName": {utils.NotEmpty()},
		}
	} else {
		verify = utils.Rules{
			"ProjectId":{utils.NotEmpty()},
			"AppName": {utils.NotEmpty()},
			"GitRepo":{utils.NotEmpty()},
			"BuildParams": {utils.NotEmpty()},
		}
	}
	if err := utils.Verify(app, verify); err != nil {
    		response.FailWithMessage(err.Error(), c)
    		return
    }
    // validate customConfig xml format
    // go xml: unsupported version "1.1"; only version 1.0 is supported
    //if (app.CustomConfig != "") {
	//	if ok := utils.IsValidXML([]byte(app.CustomConfig)); !ok {
	//		response.FailWithMessage("自定义jenkins配置xml格式错误", c)
	//		return
	//	}
	//}

	currentUserAuthorityId := c.GetUint("currentUserAuthorityId")
	// test priv
	if !(projectService.IsProjectManager(app.ProjectId, currentUserAuthorityId) || utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("创建应用失败!", zap.Error(err))
		response.FailWithMessage("创建应用失败", c)
	}


	if err := appService.CreateApp(&app); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败, "+ err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteApp 删除App
//	@Tags		App
//	@Summary	删除App
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.App	true	"删除App"
//	@Success	200		{string}	string				"{"success":true,"data":{},"msg":"删除成功"}"
//	@Router		/app/deleteApp [delete]
func (appApi *AppApi) DeleteApp(c *gin.Context) {
	var app jenkins_manage.App
	err := c.ShouldBindJSON(&app)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	currentUserAuthorityId := c.GetUint("currentUserAuthorityId")
	// test priv
	if !(projectService.IsProjectManager(app.ProjectId, currentUserAuthorityId) || utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("删除应用失败!", zap.Error(err))
		response.FailWithMessage("删除应用失败", c)
	}


	app.DeletedBy = utils.GetUserID(c)
	if err := appService.DeleteApp(app); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}


// UpdateApp 更新App
//	@Tags		App
//	@Summary	更新App
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.App	true	"更新App"
//	@Success	200		{string}	string				"{"success":true,"data":{},"msg":"更新成功"}"
//	@Router		/app/updateApp [put]
func (appApi *AppApi) UpdateApp(c *gin.Context) {
	var app jenkins_manage.App
	err := c.ShouldBindJSON(&app)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	currentUserAuthorityId := c.GetUint("currentUserAuthorityId")
	// test priv
	if !(projectService.IsProjectManager(app.ProjectId, currentUserAuthorityId) || utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("更新应用失败!", zap.Error(err))
		response.FailWithMessage("更新应用失败", c)
	}

	app.UpdatedBy = utils.GetUserID(c)

	if (app.CustomConfig == "" && app.BuildParams == nil) {
		app.BuildParams = jenkins_manage.DefaultBuildParams
	}

	var verify utils.Rules
	if(app.CustomConfig != "") {
		verify = utils.Rules{
			"GVA_MODEL.ID": {utils.NotEmpty()},
			"ProjectId": {utils.NotEmpty()},
			"AppName": {utils.NotEmpty()},
		}
	} else {
		verify = utils.Rules{
			"GVA_MODEL.ID": {utils.NotEmpty()},
			"ProjectId":{utils.NotEmpty()},
			"AppName": {utils.NotEmpty()},
			"GitRepo":{utils.NotEmpty()},
			"BuildParams": {utils.NotEmpty()},
		}
	}

    if err := utils.Verify(app, verify); err != nil {
      	response.FailWithMessage(err.Error(), c)
      	return
    }
	if err := appService.UpdateApp(app); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindApp 用id查询App
//	@Tags		App
//	@Summary	用id查询App
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	query		jenkins_manage.App	true	"用id查询App"
//	@Success	200		{string}	string				"{"success":true,"data":{},"msg":"查询成功"}"
//	@Router		/app/findApp [get]
func (appApi *AppApi) FindApp(c *gin.Context) {
	var app jenkins_manage.App
	err := c.ShouldBindQuery(&app)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if app, err := appService.GetApp(app.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"app": app}, c)
	}
}

// GetAppList 分页获取App列表
//	@Tags		App
//	@Summary	分页获取App列表
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	query		jenkins_manageReq.AppSearch	true	"分页获取App列表"
//	@Success	200		{string}	string						"{"success":true,"data":{},"msg":"获取成功"}"
//	@Router		/app/getAppList [get]
func (appApi *AppApi) GetAppList(c *gin.Context) {
	var searchInfo jenkins_manageReq.AppSearch
	//err := c.ShouldBindQuery(&searchInfo)
	err := c.BindQuery(&searchInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := appService.GetAppInfoList(searchInfo); err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败", c)
    } else {
        response.OkWithDetailed(response.PageResult{
            List:     list,
            Total:    total,
            Page:     searchInfo.Page,
            PageSize: searchInfo.PageSize,
        }, "获取成功", c)
    }
}

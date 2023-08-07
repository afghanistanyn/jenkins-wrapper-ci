package jenkins_manage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"jenkins-wrapper-ci/global"
	"jenkins-wrapper-ci/model/common/request"
	"jenkins-wrapper-ci/model/common/response"
	"jenkins-wrapper-ci/model/jenkins_manage"
	jenkins_manageReq "jenkins-wrapper-ci/model/jenkins_manage/request"
	"jenkins-wrapper-ci/utils"
)

type ProjectApi struct {
}


// CreateProject 创建Project
//	@Tags		Project
//	@Summary	创建Project
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.Project	true	"创建Project"
//	@Success	200		{string}	string					"{"success":true,"data":{},"msg":"获取成功"}"
//	@Router		/project/createProject [post]
func (projectApi *ProjectApi) CreateProject(c *gin.Context) {
	var project jenkins_manage.Project
	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    project.CreatedBy = utils.GetUserID(c)
    verify := utils.Rules{
        "ProjectName":{utils.NotEmpty()},
    }
	if err := utils.Verify(project, verify); err != nil {
    	response.FailWithMessage(err.Error(), c)
    	return
    }
    // 权限验证
	currentUserAuthorityId := c.GetUint("currentUserAuthorityId")
	// test priv
	if !(utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("创建项目失败!", zap.Error(err))
		response.FailWithMessage("创建项目失败", c)
	}

	if err := projectService.CreateProject(&project); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteProject 删除Project
//	@Tags		Project
//	@Summary	删除Project
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.Project	true	"删除Project"
//	@Success	200		{string}	string					"{"success":true,"data":{},"msg":"删除成功"}"
//	@Router		/project/deleteProject [delete]
func (projectApi *ProjectApi) DeleteProject(c *gin.Context) {
	var project jenkins_manage.Project
	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	currentUserAuthorityId := c.GetUint("currentUserAuthorityId")
	// test priv
	if !(projectService.IsProjectManager(project.ID, currentUserAuthorityId) || utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("删除项目失败!", zap.Error(err))
		response.FailWithMessage("删除项目失败", c)
	}

    project.DeletedBy = utils.GetUserID(c)
	if err := projectService.DeleteProject(project); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteProjectByIds 批量删除Project
//	@Tags		Project
//	@Summary	批量删除Project
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		request.IdsReq	true	"批量删除Project"
//	@Success	200		{string}	string			"{"success":true,"data":{},"msg":"批量删除成功"}"
//	@Router		/project/deleteProjectByIds [delete]
func (projectApi *ProjectApi) DeleteProjectByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	currentUserAuthorityId := c.GetUint("currentUserAuthorityId")
	// test priv
	if !(projectService.IsProjectManagers(IDS, currentUserAuthorityId) || utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("修改项目失败!", zap.Error(err))
		response.FailWithMessage("修改项目失败", c)
	}

    deletedBy := utils.GetUserID(c)
	if err := projectService.DeleteProjectByIds(IDS,deletedBy); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateProject 更新Project
//	@Tags		Project
//	@Summary	更新Project
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.Project	true	"更新Project"
//	@Success	200		{string}	string					"{"success":true,"data":{},"msg":"更新成功"}"
//	@Router		/project/updateProject [put]
func (projectApi *ProjectApi) UpdateProject(c *gin.Context) {
	var project jenkins_manage.Project
	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    project.UpdatedBy = utils.GetUserID(c)
    verify := utils.Rules{
          "ProjectName":{utils.NotEmpty()},
    }
    if err := utils.Verify(project, verify); err != nil {
      	response.FailWithMessage(err.Error(), c)
      	return
    }

	currentUserAuthorityId := c.GetUint("currentUserAuthorityId")
	// test priv
	if !(projectService.IsProjectManager(project.ID, currentUserAuthorityId) || utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("修改项目失败!", zap.Error(err))
		response.FailWithMessage("修改项目失败", c)
	}

	if err := projectService.UpdateProject(project); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindProject 用id查询Project
//	@Tags		Project
//	@Summary	用id查询Project
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	query		jenkins_manage.Project	true	"用id查询Project"
//	@Success	200		{string}	string					"{"success":true,"data":{},"msg":"查询成功"}"
//	@Router		/project/findProject [get]
func (projectApi *ProjectApi) FindProject(c *gin.Context) {
	var project jenkins_manage.Project
	err := c.ShouldBindQuery(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if project, err := projectService.GetProject(project.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"project": project}, c)
	}
}

// GetProjectList 分页获取Project列表
//	@Tags		Project
//	@Summary	分页获取Project列表
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	query		jenkins_manageReq.ProjectSearch	true	"分页获取Project列表"
//	@Success	200		{string}	string							"{"success":true,"data":{},"msg":"获取成功"}"
//	@Router		/project/getProjectList [get]
func (projectApi *ProjectApi) GetProjectList(c *gin.Context) {
	var pageInfo jenkins_manageReq.ProjectSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if pageInfo.Mine {
		pageInfo.CurrentUserId = c.GetUint("currentUserId")
		pageInfo.CurrentUserAuthorityId = c.GetUint("currentUserAuthorityId")
	}

	// check query param
	if len(pageInfo.Projects) > 0 && pageInfo.Mine {
		response.FailWithMessage(fmt.Errorf("查询参数冲突").Error(), c)
		return
	}

	if list, total, err := projectService.GetProjectInfoList(pageInfo); err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败", c)
    } else {
        response.OkWithDetailed(response.PageResult{
            List:     list,
            Total:    total,
            Page:     pageInfo.Page,
            PageSize: pageInfo.PageSize,
        }, "获取成功", c)
    }
}

// SetProjectMembers 设置项目管理员或普通成员
//	@Tags		Project
//	@Summary	设置项目管理员或普通成员
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body 	jenkins_manageReq.SetProjectMemberReq	true	"设置成功"
//	@Success	200		{string}		string	"{"success":true,"data":{},"msg":"设置成功"}"
//	@Router		/project/setProjectMembers [post]
func (projectApi *ProjectApi) SetProjectMembers(c *gin.Context) {
	var setProjectMemberReq jenkins_manageReq.SetProjectMemberReq
	err := c.ShouldBindJSON(&setProjectMemberReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	currentUserAuthorityId := c.GetUint("currentUserAuthorityId")
	// test priv
	if !(projectService.IsProjectManager(setProjectMemberReq.ProjectId, currentUserAuthorityId) || utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("设置成员失败!", zap.Error(err))
		response.FailWithMessage("设置成员失败", c)
	}

	verify := utils.Rules{
		"ProjectId":{utils.NotEmpty()},
		"MemberType":{utils.NotEmpty()},
		"MemberIds":{utils.NotEmpty()},
	}
	if err := utils.Verify(setProjectMemberReq, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	currentUserId := c.GetUint("currentUserId")
	if err := projectService.SetProjectMembers(setProjectMemberReq, currentUserId); err != nil {
		global.GVA_LOG.Error("设置成员失败!", zap.Error(err))
		response.FailWithMessage("设置成员失败", c)
	} else {
		response.OkWithMessage( "设置成员成功", c)
	}
}
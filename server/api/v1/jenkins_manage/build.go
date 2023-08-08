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

type BuildApi struct {
}



// CreateBuild 创建Build
//	@Tags		Build
//	@Summary	创建Build
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.Build	true	"创建Build"
//	@Success	200		{string}	string					"{"success":true,"data":{},"msg":"获取成功"}"
//	@Router		/build/createBuild [post]
func (buildApi *BuildApi) CreateBuild(c *gin.Context) {
	var build jenkins_manage.Build
	err := c.ShouldBindJSON(&build)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    verify := utils.Rules{
        "AppId":{utils.NotEmpty()},
		"AppName":{utils.NotEmpty()},
		"ProjectId":{utils.NotEmpty()},
		"ProjectName":{utils.NotEmpty()},
		"BuildParamValues":{utils.NotEmpty()},
		"ApproveStatus": {utils.Eq("1")},
    }
	if err := utils.Verify(build, verify); err != nil {
    		response.FailWithMessage(err.Error(), c)
    		return
    }

    build.CreatedBy = utils.GetUserID(c)
    build.UpdatedBy = build.CreatedBy
	if err := buildService.CreateBuild(&build); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		go notificationService.ApproveBuildNotification(build.ID)
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteBuild 删除Build
//	@Tags		Build
//	@Summary	删除Build
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.Build	true	"删除Build"
//	@Success	200		{string}	string					"{"success":true,"data":{},"msg":"删除成功"}"
//	@Router		/build/deleteBuild [delete]
func (buildApi *BuildApi) DeleteBuild(c *gin.Context) {
	var build jenkins_manage.Build
	err := c.ShouldBindJSON(&build)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := buildService.DeleteBuild(build); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// ApproveBuild 审核Build
//	@Tags		Build
//	@Summary	审核Build
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	body		jenkins_manage.Build	true	"审核Build"
//	@Success	200		{string}	string					"{"success":true,"data":{},"msg":"审核成功"}"
//	@Router		/build/approveBuild [post]
func (buildApi *BuildApi) ApproveBuild(c *gin.Context) {
	var approveBuildReq jenkins_manageReq.ApproveBuildReq
	err := c.ShouldBindJSON(&approveBuildReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"ID":{utils.NotEmpty()},
		"ApproveStatus": {utils.NotEmpty()},
	}
	if err := utils.Verify(approveBuildReq, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	currentUserId := utils.GetUserID(c)
	currentUserAuthorityId := utils.GetUserAuthorityId(c)

	approveBuildReq.ApproveBy = currentUserId
	approveBuildReq.UpdatedBy = currentUserId

	build, err := buildService.GetBuildWithoutLog(approveBuildReq.ID)
	if err != nil {
		err = fmt.Errorf("参数错误")
		global.GVA_LOG.Error("审核发布失败!", zap.Error(err))
		response.FailWithMessage("审核发布失败", c)
	}

	//权限验证
	if !(projectService.IsProjectManager(build.ProjectId, currentUserAuthorityId) || utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("审核发布失败!", zap.Error(err))
		response.FailWithMessage("审核发布失败", c)
	}

	if err := buildService.ApproveBuild(approveBuildReq); err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败", c)
	} else {
		response.OkWithMessage("审核成功", c)
	}
}

// FindBuild 用id查询Build
//	@Tags		Build
//	@Summary	用id查询Build
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	query		jenkins_manage.Build	true	"用id查询Build"
//	@Success	200		{string}	string					"{"success":true,"data":{},"msg":"查询成功"}"
//	@Router		/build/findBuild [get]
func (buildApi *BuildApi) FindBuild(c *gin.Context) {
	var build_ jenkins_manage.Build
	err := c.ShouldBindQuery(&build_)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if build, err := buildService.GetBuild(build_.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"build": build}, c)
	}
}

// GetBuildList 分页获取Build列表
//	@Tags		Build
//	@Summary	分页获取Build列表
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	query		jenkins_manageReq.BuildSearch	true	"分页获取Build列表"
//	@Success	200		{string}	string							"{"success":true,"data":{},"msg":"获取成功"}"
//	@Router		/build/getBuildList [get]
func (buildApi *BuildApi) GetBuildList(c *gin.Context) {
	var pageInfo jenkins_manageReq.BuildSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	pageInfo.CurrentUserId = utils.GetUserID(c)
	pageInfo.CurrentUserAuthorityId = utils.GetUserAuthorityId(c)

	if list, total, err := buildService.GetBuildInfoList(pageInfo); err != nil {
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


// GetBuildList 分页获取Build列表
//	@Tags		Build
//	@Summary	获取发布详情
//	@Security	ApiKeyAuth
//	@accept		application/json
//	@Produce	application/json
//	@Param		data	query		jenkins_manageReq.BuildInfoReq	true	"获取发布详情"
//	@Success	200		{string}	string							"{"success":true,"data":{},"msg":"获取发布详情成功"}"
//	@Router		/build/buildInfo [post]
func (buildApi *BuildApi) BuildInfo(c *gin.Context) {
	var buildInfoReq jenkins_manageReq.BuildInfoReq
	err := c.ShouldBindJSON(&buildInfoReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	currentUserAuthorityId := utils.GetUserAuthorityId(c)
	//权限验证
	if !(utils.IsSuperAdmin(currentUserAuthorityId)) {
		err = fmt.Errorf("无权限")
		global.GVA_LOG.Error("获取发布详情失败!", zap.Error(err))
		response.FailWithMessage("获取发布详情失败", c)
	}

	if buidldInfo, err := buildService.GetJenkinsBuildInfo(buildInfoReq); err != nil {
		global.GVA_LOG.Error("获取发布详情失败!", zap.Error(err))
		response.FailWithMessage("获取发布详情失败", c)
	} else {
		response.OkWithDetailed(buidldInfo, "获取发布详情成功", c)
	}
}
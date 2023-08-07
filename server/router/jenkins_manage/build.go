package jenkins_manage

import (
	"jenkins-wrapper-ci/api/v1"
	"jenkins-wrapper-ci/middleware"
	"github.com/gin-gonic/gin"
)

type BuildRouter struct {
}

// InitBuildRouter 初始化 Build 路由信息
func (s *BuildRouter) InitBuildRouter(Router *gin.RouterGroup) {
	buildRouter := Router.Group("build").Use(middleware.OperationRecord())
	buildRouterWithoutRecord := Router.Group("build")
	var buildApi = v1.ApiGroupApp.Jenkins_manageApiGroup.BuildApi
	{
		buildRouter.POST("createBuild", buildApi.CreateBuild)   // 新建Build
		buildRouter.DELETE("deleteBuild", buildApi.DeleteBuild) // 删除Build
		buildRouter.POST("approveBuild", buildApi.ApproveBuild) //审核Build
		buildRouter.POST("buildInfo", buildApi.BuildInfo) //查看jenkins build详情
	}
	{
		buildRouterWithoutRecord.GET("findBuild", buildApi.FindBuild)        // 根据ID获取Build
		buildRouterWithoutRecord.GET("getBuildList", buildApi.GetBuildList)  // 获取Build列表
	}
}

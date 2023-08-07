package jenkins_manage

import (
	"jenkins-wrapper-ci/api/v1"
	"jenkins-wrapper-ci/middleware"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
}

// InitAppRouter 初始化 App 路由信息
func (s *AppRouter) InitAppRouter(Router *gin.RouterGroup) {
	appRouter := Router.Group("app").Use(middleware.OperationRecord())
	appRouterWithoutRecord := Router.Group("app")
	var appApi = v1.ApiGroupApp.Jenkins_manageApiGroup.AppApi
	{
		appRouter.POST("createApp", appApi.CreateApp)   // 新建App
		appRouter.DELETE("deleteApp", appApi.DeleteApp) // 删除App
		appRouter.PUT("updateApp", appApi.UpdateApp)    // 更新App
	}
	{
		appRouterWithoutRecord.GET("findApp", appApi.FindApp)        // 根据ID获取App
		appRouterWithoutRecord.GET("getAppList", appApi.GetAppList)  // 获取App列表
	}
}

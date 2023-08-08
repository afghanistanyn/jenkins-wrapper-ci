package jenkins_manage

import "jenkins-wrapper-ci/service"

type ApiGroup struct {
	ProjectApi
	AppApi
	BuildApi
}

var (
	appService = service.ServiceGroupApp.JenkinsManageServicegroup.AppService
	buildService = service.ServiceGroupApp.JenkinsManageServicegroup.BuildService
	projectService = service.ServiceGroupApp.JenkinsManageServicegroup.ProjectService
	notificationService = service.ServiceGroupApp.JenkinsManageServicegroup.NotificationService
)
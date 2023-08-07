package service

import (
	"jenkins-wrapper-ci/service/example"
	"jenkins-wrapper-ci/service/jenkins_manage"
	"jenkins-wrapper-ci/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup        system.ServiceGroup
	ExampleServiceGroup       example.ServiceGroup
	JenkinsManageServicegroup jenkins_manage.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)

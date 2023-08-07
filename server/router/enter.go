package router

import (
	"jenkins-wrapper-ci/router/example"
	"jenkins-wrapper-ci/router/jenkins_manage"
	"jenkins-wrapper-ci/router/system"
)

type RouterGroup struct {
	System         system.RouterGroup
	Example        example.RouterGroup
	Jenkins_manage jenkins_manage.RouterGroup
}

var RouterGroupApp = new(RouterGroup)

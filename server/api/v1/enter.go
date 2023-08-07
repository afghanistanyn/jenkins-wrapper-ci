package v1

import (
	"jenkins-wrapper-ci/api/v1/example"
	"jenkins-wrapper-ci/api/v1/jenkins_manage"
	"jenkins-wrapper-ci/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup         system.ApiGroup
	ExampleApiGroup        example.ApiGroup
	Jenkins_manageApiGroup jenkins_manage.ApiGroup
}

var ApiGroupApp = new(ApiGroup)

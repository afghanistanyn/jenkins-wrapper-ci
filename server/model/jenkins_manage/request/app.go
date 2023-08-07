package request

import (
	"jenkins-wrapper-ci/model/common/request"
	"jenkins-wrapper-ci/model/jenkins_manage"
	"time"
)

type AppSearch struct {
	jenkins_manage.App
	ProjectId uint `json:"projectId" form:"projectId"`
	BuildParam string `json:"buildParam" form:"buildParam"`
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}
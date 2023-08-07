package request

import (
	"jenkins-wrapper-ci/model/jenkins_manage"
	"jenkins-wrapper-ci/model/common/request"
	"time"
)

type ProjectSearch struct{
    jenkins_manage.Project
	Mine        bool   `json:"mine" form:"mine"`
    ProjectInfo string `json:"projectInfo" form:"projectInfo"`
	AppInfo     string `json:"appInfo" form:"appInfo"`
	CurrentUserId          uint
	CurrentUserAuthorityId uint
    Projects               []uint `json:"projects" form:"projects"`
    StartCreatedAt         *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt           *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    request.PageInfo
}


type SetProjectMemberReq struct {
	ProjectId uint		`json:"projectId"`
	MemberType string	`json:"memberType"`
	MemberIds []uint	`json:"memberIds"`
	updatedBy	uint
}

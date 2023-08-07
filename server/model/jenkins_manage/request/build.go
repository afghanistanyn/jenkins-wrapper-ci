package request

import (
	"jenkins-wrapper-ci/model/jenkins_manage"
	"jenkins-wrapper-ci/model/common/request"
	"time"
)

type BuildSearch struct{
    jenkins_manage.Build
	Mine bool 					`json:"mine" form:"mine"`
	CurrentUserId		 		 uint
	CurrentUserAuthorityId 		 uint
	Projects []uint 			`json:"projects" form:"projects"`
    Apps     []uint 			`json:"apps" form:"apps"`
    StartCreatedAt *time.Time 	`json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time 	`json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

type ApproveBuildReq struct {
	ID            uint `json:"id"`
	ApproveStatus int  `json:"approveStatus"`
	ApproveBy     uint
	UpdatedBy     uint
}


type BuildInfoReq struct {
	ID        uint `json:"id"`
}
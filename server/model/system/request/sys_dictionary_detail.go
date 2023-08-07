package request

import (
	"jenkins-wrapper-ci/model/common/request"
	"jenkins-wrapper-ci/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}

package request

import (
	"jenkins-wrapper-ci/model/common/request"
	"jenkins-wrapper-ci/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}

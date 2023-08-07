package request

import (
	"jenkins-wrapper-ci/model/common/request"
	"jenkins-wrapper-ci/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}

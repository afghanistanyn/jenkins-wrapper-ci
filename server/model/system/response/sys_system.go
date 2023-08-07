package response

import "jenkins-wrapper-ci/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}

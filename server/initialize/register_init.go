package initialize

import (
	_ "jenkins-wrapper-ci/source/example"
	_ "jenkins-wrapper-ci/source/system"
	_ "jenkins-wrapper-ci/source/jenkins_manage"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}

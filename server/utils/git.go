package utils

import (
	"fmt"
	"jenkins-wrapper-ci/global"
	"net/url"
	"strings"
)


func FindGitCredentialId(gitRepo string) (gitCredentialId string, err error) {
	credentials := global.GVA_CONFIG.GitCredentials
	if len(credentials) == 0 {
		return gitCredentialId,  fmt.Errorf("no git-credentials configs")
	}

		gitUrl, err := url.Parse(gitRepo)
	if err != nil {
		return gitCredentialId, err
	}

	for _, cred := range credentials {
		if ( strings.Contains(gitUrl.Host, cred.GitServer) || strings.Contains(gitRepo, cred.GitServer)) {
			return cred.GitCredentialId, nil
		}
	}

	return gitCredentialId,  fmt.Errorf("gitCredentialId for gitServer[%s] not found", gitUrl.Host)

}
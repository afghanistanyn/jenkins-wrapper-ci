package utils

import (
	"jenkins-wrapper-ci/global"
)

func contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IsSuperAdmin(authorityId uint) bool  {
	var SuperAdminRoles = global.GVA_CONFIG.SuperAdminRoles
	if contains(SuperAdminRoles, authorityId) {
		return true
	}
	return false
}

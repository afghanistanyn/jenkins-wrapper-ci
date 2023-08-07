package jenkins_manage

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"jenkins-wrapper-ci/global"
	"jenkins-wrapper-ci/model/common/request"
	"jenkins-wrapper-ci/model/jenkins_manage"
	jenkins_manageReq "jenkins-wrapper-ci/model/jenkins_manage/request"
	"jenkins-wrapper-ci/model/system"
	"jenkins-wrapper-ci/utils"
)

type ProjectService struct {
}

// CreateProject 创建Project记录
// Author [piexlmax](https://github.com/piexlmax)
func (projectService *ProjectService) CreateProject(project *jenkins_manage.Project) (err error) {
	err = global.GVA_DB.Create(project).Error
	return err
}

// DeleteProject 删除Project记录
func (projectService *ProjectService)DeleteProject(project jenkins_manage.Project) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var apps []jenkins_manage.App
		if err := tx.Model(&jenkins_manage.App{}).Where("project_id = ?", project.ID).Find(&apps).Error; err != nil {
			return err
		}

		// for delete app hook， delete app one by one
		for _, app := range apps {
			global.GVA_LOG.Info("delete app in project", zap.Any("project", app.ProjectId), zap.Any("app", app.AppName))
			if err := tx.Model(&jenkins_manage.App{}).Delete(&app).Error; err != nil {
				return err
			}
		}

        if err = tx.Model(&jenkins_manage.Project{}).Unscoped().Delete(&project).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteProjectByIds 批量删除Project记录
func (projectService *ProjectService)DeleteProjectByIds(ids request.IdsReq,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&jenkins_manage.App{}).Unscoped().Where("project_id in ?", ids).Delete(&jenkins_manage.App{}).Error; err != nil {
			return err
		}

        if err := tx.Model(&jenkins_manage.Project{}).Unscoped().Where("id in ?", ids.Ids).Delete(&jenkins_manage.Project{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateProject 更新Project记录
func (projectService *ProjectService)UpdateProject(project jenkins_manage.Project) (err error) {
	db := global.GVA_DB.Model(&jenkins_manage.Project{})
	// 不允许修改项目名
	err = db.Omit("ProjectName").Save(&project).Error
	return err
}

// GetProject 根据id获取Project记录
func (projectService *ProjectService)GetProject(id uint) (project jenkins_manage.Project, err error) {
	db := global.GVA_DB.Model(&jenkins_manage.Project{}).Preload("Managers").Preload("Members")
	err = db.Where("id = ?", id).First(&project).Error
	return
}

// GetProjectInfoList 分页获取Project记录
// Author [piexlmax](https://github.com/piexlmax)
func (projectService *ProjectService)GetProjectInfoList(info jenkins_manageReq.ProjectSearch) (list []jenkins_manage.Project, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&jenkins_manage.Project{}).Preload("Managers").Preload("Members").Preload("Apps").Order("created_at desc")
    var projects []jenkins_manage.Project
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }

    if info.ProjectInfo != "" {
		db = db.Where("name LIKE ?", "%"+info.ProjectInfo+"%").Or("name_cn LIKE ?", "%"+info.ProjectInfo+"%").Or("description LIKE ?", "%"+info.ProjectInfo+"%")
	}

	if info.AppInfo != "" {
		db = db.Where("id in  (select project_id from app where name LIKE ? or description LIKE ?)", "%"+info.AppInfo+"%", "%"+info.AppInfo+"%")
	}

	if info.ProjectName != "" {
		db = db.Where("name LIKE ?", "%"+info.ProjectName+"%")
	}
	if info.NameCn != "" {
		db = db.Where("nameCn LIKE ?", "%"+info.NameCn+"%")
	}
	if info.Description != "" {
		db = db.Where("description LIKE ?", "%"+info.Description+"%")
	}

	if info.Mine {
		// if ops or admin role
		isSuperAdmin := utils.IsSuperAdmin(info.CurrentUserAuthorityId)
		if !isSuperAdmin {
			db.Where("project.id in (select project_id from project_managers where sys_user_id = ? union select project_id from project_members where sys_user_id = ?)", info.CurrentUserId, info.CurrentUserId)
		}
	}

	if len(info.Projects) > 0 {
		db.Where("id in ?", info.Projects)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 && offset != 0 {
		err = db.Limit(limit).Offset(offset).Find(&projects).Error
	} else {
		err = db.Find(&projects).Error
	}

	//.Preload("Apps").Joins("right join app t1 on t1.project_id = project.id")
	// 手动填充 apps of project

	return  projects, total, err
}

// SetProjectMembers 设置项目管理员或普通成员
func (projectService *ProjectService)SetProjectMembers(setProjectMemberReq jenkins_manageReq.SetProjectMemberReq, updateBy uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var project jenkins_manage.Project
		if err := tx.Model(&jenkins_manage.Project{}).Preload("Managers").Preload("Members").Where("id = ?", setProjectMemberReq.ProjectId).Find(&project).Error; err != nil {
			return err
		}

		if (setProjectMemberReq.MemberType == "members") {
			var members []system.SysUser
			if err := tx.Model(&system.SysUser{}).Where("id in ?", setProjectMemberReq.MemberIds).Find(&members).Error; err != nil {
				return err
			}
			if err = tx.Model(&project).Association("Members").Replace(members); err != nil {
				return err
			}

		} else if (setProjectMemberReq.MemberType == "managers") {
			var managers []system.SysUser
			if err := tx.Model(&system.SysUser{}).Where("id in ?", setProjectMemberReq.MemberIds).Find(&managers).Error; err != nil {
				return err
			}
			if err = tx.Model(&project).Association("Managers").Replace(managers); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("error parameter MemberType")
		}

		if err := tx.Model(&project).UpdateColumn("updated_by", updateBy).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (projectService *ProjectService)IsProjectManager(projectId, authorityId uint) bool  {
	var project jenkins_manage.Project
	db := global.GVA_DB.Model(&jenkins_manage.Project{}).Preload("Managers")
	if err:= db.Where("id = ?", projectId).Find(&project).Error; err != nil {
		return false
	}
	for _, m := range project.Managers {
		if m.AuthorityId == authorityId {
			return true
		}
	}
	return false
}

func (projectService *ProjectService)IsProjectManagers(projectIds request.IdsReq, authorityId uint) bool  {
	var isProjectManagers bool = false
	var projects []jenkins_manage.Project
	db := global.GVA_DB.Model(&jenkins_manage.Project{}).Preload("Managers")
	if err:= db.Where("id in ?", projectIds).Find(&projects).Error; err != nil {
		return false
	}
	for _, project := range projects {
		for _, m := range project.Managers {
			if m.AuthorityId == authorityId {
				isProjectManagers = true
			} else {
				return false
			}
		}
	}
	return isProjectManagers
}

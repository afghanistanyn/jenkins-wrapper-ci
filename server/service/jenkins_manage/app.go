package jenkins_manage

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"jenkins-wrapper-ci/global"
	"jenkins-wrapper-ci/model/jenkins_manage"
	jenkins_manageReq "jenkins-wrapper-ci/model/jenkins_manage/request"
)
type AppService struct {
}


// CreateApp 创建App记录
// Author [piexlmax](https://github.com/piexlmax)
func (appService *AppService) CreateApp(app *jenkins_manage.App) (err error) {
	err = global.GVA_DB.Create(app).Error
	return err
}

// DeleteApp 删除App记录
// Author [piexlmax](https://github.com/piexlmax)
func (appService *AppService)DeleteApp(app jenkins_manage.App) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
        if err = tx.Unscoped().Delete(&app).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}


// UpdateApp 更新App记录
// Author [piexlmax](https://github.com/piexlmax)
func (appService *AppService)UpdateApp(app jenkins_manage.App) (err error) {
	// 不允许修改名称, 所属项目
	err = global.GVA_DB.Omit("ProjectId", "AppName").Save(&app).Error
	return err
}

// GetApp 根据id获取App记录
// Author [piexlmax](https://github.com/piexlmax)
func (appService *AppService)GetApp(id uint) (app jenkins_manage.App, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&app).Error
	return
}

// GetAppInfoList 分页获取App记录
// Author [piexlmax](https://github.com/piexlmax)
func (appService *AppService)GetAppInfoList(info jenkins_manageReq.AppSearch) (list []jenkins_manage.App, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&jenkins_manage.App{}).Preload("Project").Order("created_at desc")
    var apps []jenkins_manage.App
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.AppName != "" {
    	db = db.Where("name LIKE ?", "%"+info.AppName+"%")
	}
	if info.GitRepo != "" {
		db = db.Where("git_repo LIKE ?", "%"+info.GitRepo+"%")
	}
	if info.Image != "" {
		db = db.Where("image LIKE ?", "%"+info.Image+"%")
	}
	if info.ProjectId > 0 {
		db = db.Where("project_id = ?", info.ProjectId)
	}

	// query json_array contains
	if info.BuildParam != "" {
		db = db.Where(datatypes.JSONArrayQuery("build_params").Contains(info.BuildParam))
	}


	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
    if limit != 0 && offset != 0 {
		err = db.Limit(limit).Offset(offset).Find(&apps).Error
	} else {
		err = db.Find(&apps).Error
	}

	return  apps, total, err
}

func (appService *AppService)GetGitRepo(appId uint) (string, error) {
	var app jenkins_manage.App
	if err := global.GVA_DB.Model(&app).Where("id = ?", appId).Find(&app).Error; err != nil {
		return "", err
	}
	return app.GitRepo, nil
}
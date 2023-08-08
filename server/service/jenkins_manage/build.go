package jenkins_manage

import (
	"context"
	"encoding/json"
	"github.com/bndr/gojenkins"
	"go.uber.org/zap"
	"jenkins-wrapper-ci/global"
	"jenkins-wrapper-ci/model/jenkins_manage"
	jenkins_manageReq "jenkins-wrapper-ci/model/jenkins_manage/request"
	"jenkins-wrapper-ci/utils"
	"strconv"
	"time"
)

type BuildService struct {
}

// CreateBuild 创建Build记录
// Author [piexlmax](https://github.com/piexlmax)
func (buildService *BuildService) CreateBuild(build *jenkins_manage.Build) (err error) {
	err = global.GVA_DB.Create(build).Error
	return err
}

// DeleteBuild 删除Build记录
// Author [piexlmax](https://github.com/piexlmax)
func (buildService *BuildService)DeleteBuild(build jenkins_manage.Build) (err error) {
	err = global.GVA_DB.Unscoped().Delete(&build).Error
	return err
}


// ApproveBuild
func (buildService *BuildService)ApproveBuild(approveBuildReq jenkins_manageReq.ApproveBuildReq) (err error) {
	//mark approve update
	approveBuildReqJson, _ := json.Marshal(approveBuildReq)
	db := global.GVA_DB.Model(&jenkins_manage.Build{}).Set("approve_update", true).Set("approveBuildReq", approveBuildReqJson)
	now := time.Now()
	err = db.Where("id = ?", approveBuildReq.ID).Updates(jenkins_manage.Build{
		GVA_MODEL:     global.GVA_MODEL{
			UpdatedAt: now,
		},
		BuildAt: &now,
		ApproveStatus: approveBuildReq.ApproveStatus,
		ApproveBy:     approveBuildReq.ApproveBy,
		UpdatedBy:     approveBuildReq.UpdatedBy,
	}).Error
	return err
}

func (buildService *BuildService)UpdateBuild(build jenkins_manage.Build) (err error) {
	err = global.GVA_DB.Omit("approveStatus", "result").Save(&build).Error
	return err
}

func (buildService *BuildService)GetBuild(id uint) (build jenkins_manage.Build, err error) {
	db := global.GVA_DB.Model(&jenkins_manage.Build{}).Set("fetch_log", true)
	err = db.Where("id = ?", id).First(&build).Error
	return
}

func (buildService *BuildService)GetBuildWithoutLog(id uint) (build jenkins_manage.Build, err error) {
	err = global.GVA_DB.Omit("log").Where("id = ?", id).First(&build).Error
	return
}


// GetBuildInfoList 分页获取Build记录
func (buildService *BuildService)GetBuildInfoList(info jenkins_manageReq.BuildSearch) (list []jenkins_manage.Build, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&jenkins_manage.Build{}).Omit("log").Order("created_at desc")
    var builds []jenkins_manage.Build
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.ProjectName != "" {
		db = db.Where("project_name LIKE ?", "%"+info.ProjectName+"%")
	}
	if info.AppName != "" {
		db = db.Where("app_name LIKE ?", "%"+info.AppName+"%")
	}
	if info.ProjectId != 0 {
		db = db.Where("project_id = ?", info.ProjectId)
	}
	if info.AppId != 0 {
		db = db.Where("app_id = ?", info.AppId)
	}
	if info.Image != "" {
		db = db.Where("image LIKE ?", "%"+info.Image+"%")
	}
	if info.GitCommit != "" {
		db = db.Where("git_commit LIKE ?", "%"+info.GitCommit+"%")
	}
	if info.ApproveBy != 0 {
		db = db.Where("approve_by = ?", info.ApproveBy)
	}
	if info.ApproveStatus != 0 {
		db = db.Where("approve_status = ?", info.ApproveStatus)
	}
	if info.Result != 0 {
		db = db.Where("result = ?", info.Result)
	}
	if info.Mine {
		// if ops or admin role
		isSuperAdmin := utils.IsSuperAdmin(info.CurrentUserAuthorityId)
		if !isSuperAdmin {
			db.Where("project_id in (select project_id from project_managers where sys_user_id = ? union select project_id from project_members where sys_user_id = ?)", info.CurrentUserId, info.CurrentUserId)
		}
	}
	if len(info.Apps) > 0 {
		db.Where("app_id in ?", info.Apps)
	}
	if len(info.Projects) > 0 {
		db.Where("project_id in ?", info.Projects)
	}
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 && offset != 0 {
		err = db.Limit(limit).Offset(offset).Find(&builds).Error
	} else {
		err = db.Find(&builds).Error
	}
	return  builds, total, err
}


func (buildService *BuildService)GetJenkinsBuildInfo(buildInfoReq jenkins_manageReq.BuildInfoReq) (buildInfo *gojenkins.BuildResponse, err error) {

	var build jenkins_manage.Build
	err = global.GVA_DB.Model(&jenkins_manage.Build{}).Where("id = ?", buildInfoReq.ID).First(&build).Error
	if err != nil {
		return buildInfo, err
	}

	var app jenkins_manage.App
	err = global.GVA_DB.Model(&jenkins_manage.App{}).Preload("Project").Where("id = ?", build.AppId).First(&app).Error
	if err != nil {
		return buildInfo, err
	}

	jenkinsService := jenkins_manage.GetJenkinsService()
	buildInfo, err = jenkinsService.GetJenkinsJobBuildInfo(context.Background(), &app, int64(build.BuildNumber))
	return buildInfo, err
}

func PullRunningBuilds() {

	for {
		var runningBuilds []jenkins_manage.Build
		if err := global.GVA_DB.Model(&jenkins_manage.Build{}).Where("result = ?", jenkins_manage.BuildResultMap["BUILDING"]).Find(&runningBuilds).Error; err != nil {
			global.GVA_LOG.Warn("get running jenkins builds err", zap.Error(err))
			continue
		}

		jenkinsService := jenkins_manage.GetJenkinsService()
		for _,build := range runningBuilds {
			folderJob := jenkinsService.GetFolderJobObj(build.ProjectName, build.AppName)
			_, err := folderJob.Poll(context.Background())
			if err != nil {
				global.GVA_LOG.Warn("get jenkins job err", zap.Any("job", build.ProjectName + "/" + build.AppName),  zap.Error(err))
				continue
			}
			if(build.BuildNumber == 0) {
				global.GVA_LOG.Warn("jenkins job build number err, buildNumber can't be zero", zap.Any("build.ID", build.ID))
				continue
			}

			jenkinsBuild, err := folderJob.GetBuild(context.Background(), int64(build.BuildNumber))
			if err != nil {
				global.GVA_LOG.Warn("get jenkins job build detail err", zap.Any("job", build.ProjectName + "/" + build.AppName + "#" + strconv.Itoa(int(build.BuildNumber))), zap.Any("build.ID", build.ID),   zap.Error(err))
				continue
			}
			_, _ = jenkinsBuild.Poll(context.Background())

			if (jenkinsBuild.Raw.Building) {
				global.GVA_LOG.Debug("the jenkins job build still running", zap.Any("job", build.ProjectName + "/" + build.AppName),  zap.Error(err))
				continue
			}

			// build done
			if(jenkinsBuild.Raw.Result == "FAILURE" || jenkinsBuild.Raw.Result == "ABORTED" || jenkinsBuild.Raw.Result == "UNSTABLE") {
				build.Result = jenkins_manage.BuildResultMap["FAILIED"]
			} else if(jenkinsBuild.Raw.Result == "SUCCESS") {
				build.Result = jenkins_manage.BuildResultMap["SUCCESS"]
			}

			build.GitCommit = jenkinsBuild.GetRevision()


			changes, err := json.Marshal(jenkinsBuild.Raw.ChangeSets)
			if err != nil {
				global.GVA_LOG.Warn("get jenkins job build changes err", zap.Any("job", build.ProjectName + "/" + build.AppName),  zap.Error(err))
				build.Changes = ""
			}
			build.Changes = string(changes)
			build.Image = ""
			build.Log = jenkinsBuild.GetConsoleOutput(context.Background())

			// set update by admin
			build.UpdatedBy = 1
			build.UpdatedAt = time.Now()

			if err := global.GVA_DB.Model(&jenkins_manage.Build{}).Where("id = ?", build.ID).Save(&build).Error; err != nil {
				global.GVA_LOG.Warn("save build err", zap.Error(err))
				continue
			}

			var notificationService NotificationService
			go notificationService.BuildDoneNotification(build)
		}
		time.Sleep(4 * 1000 * time.Millisecond)
	}

}
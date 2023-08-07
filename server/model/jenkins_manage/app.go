// 自动生成模板App
package jenkins_manage

import (
      "context"
      "encoding/json"
      "fmt"
      "go.uber.org/zap"
      "gorm.io/datatypes"
      "gorm.io/gorm"
      "jenkins-wrapper-ci/global"
      "jenkins-wrapper-ci/utils"
      "strings"
)

type BuildParams []string
type BuildParamValues map[string]string
var DefaultBuildParams = datatypes.JSON([]byte(`["branch", "build_env"]`))

type JenkinsJobConfig struct {
      Name              string
      Description       string
      GitRepo           string
      GitCredentialId   string
      Jenkinsfile       string
      BuildParams       []string

}

type CustomGitAuth struct {
      User string
      Password   string
}

// App 结构体
type App struct {
      global.GVA_MODEL
      // form tag is required, it work for bind queryparams
      AppName      string           `json:"name" form:"name" gorm:"column:name;uniqueIndex:project_app_index;comment:应用名;"`
      Description  string           `json:"description" form:"description" gorm:"column:description;comment:应用描述;"`
      GitRepo      string           `json:"gitRepo" form:"gitRepo" gorm:"column:git_repo;comment:git仓库地址;"`
      Image        string           `json:"image" form:"image" gorm:"column:image;comment:镜像地址;"`
      Jenkinsfile  string           `json:"jenkinsFile" form:"jenkinsFile" gorm:"column:jenkins_file;default:Jenkinsfile;comment:jenkinsFile path"`
      BuildParams  datatypes.JSON   `json:"buildParams" form:"buildParams" gorm:"column:build_params;comment:构建参数"`
      // app belongs to proj
      ProjectId    uint             `json:"projectId" form:"projectId" gorm:"column:project_id;uniqueIndex:project_app_index;comment:所属项目Id;"`
      // same to jenkins folder
      Project      Project          `json:"-" form:"-" gorm:"foreignKey:ID;references:ProjectId;comment:所属项目"`
      //custom jenkins job config, format: xml
      CustomConfig  string          `json:"customConfig" form:"customConfig" gorm:"column:custom_config;type:text;comment:custom_jenkins_job_config.xml"`
      // 保留字段, 暂未使用
      CustomWebHook      string     `json:"customWebhook" form:"customWebhook" gorm:"column:custom_webhook;comment:customWebhook"`
      CreatedBy    uint             `json:"createdBy" form:"createdBy" gorm:"column:created_by;comment:创建者"`
      UpdatedBy    uint             `json:"updateBy" form:"updateBy" gorm:"column:updated_by;comment:更新者"`
      DeletedBy    uint             `json:"DeletedBy" form:"DeletedBy" gorm:"column:deleted_by;comment:删除者"`

}

// TableName App 表名
func (App) TableName() string {
  return "app"
}

func (app *App) GetJenkinsJobConfigData() (jenkinsJobConfig JenkinsJobConfig, err error)  {

      var buildParams []string
      if err := json.Unmarshal(app.BuildParams, &buildParams); err != nil {
            return JenkinsJobConfig{}, fmt.Errorf("Unmarshal App BuildParams err, %s", err.Error())
      }

      gitCredentialId, err := utils.FindGitCredentialId(app.GitRepo)
      if err != nil {
            return JenkinsJobConfig{}, fmt.Errorf("find gitRepo Credential err, %s", err.Error())
      }

      // 描述中包含中文字符, 转为10进制unicode, format: &#xxxx;
      var desc string
      for _, r := range app.Description {
            desc += utils.ChineseToHTMLEntity(r)
      }

      return JenkinsJobConfig{
            Name:            app.AppName,
            Description:     desc,
            GitRepo:         app.GitRepo,
            GitCredentialId: gitCredentialId,
            Jenkinsfile:     app.Jenkinsfile,
            BuildParams:     buildParams,
      }, nil
}



func (app *App) AfterCreate(tx *gorm.DB) (err error) {
      jenkinsService := GetJenkinsService()
      var project Project
      if err := tx.Model(&Project{}).Where("id = ?", app.ProjectId).First(&project).Error; err != nil {
            global.GVA_LOG.Error("create jenkins job, get app.Project error", zap.Any("app", app.AppName), zap.Error(err))
            return err
      }
      app.Project = project
      _, err = jenkinsService.CreateJenkinsJob(context.Background(), app)
      if err != nil {
            global.GVA_LOG.Error("create jenkins job error", zap.Any("app", app.AppName), zap.Error(err))
            return err
      }

      //parse buildParams
      folderJob := jenkinsService.GetFolderJobObj(app.Project.ProjectName, app.AppName)
      // ignore if buildParam not exist
      buildParameterDefinitions, err := folderJob.GetParameters(context.Background())
      if err != nil {
            if strings.Contains(err.Error(), "404") {
                  return nil
            }
            global.GVA_LOG.Info("create jenkins job get job buildParams err", zap.Any("app", app.AppName), zap.Error(err))
            // todo rollback delete job
            return err

      }
      var buildParams BuildParams
      for _, buildParameterDefinition := range buildParameterDefinitions {
            buildParams = append(buildParams, buildParameterDefinition.Name)
      }

      app.BuildParams, err = json.Marshal(buildParams)
      if err != nil {
            global.GVA_LOG.Error("create jenkins job marshal buildParams err", zap.Any("app", app.AppName), zap.Error(err))
            // todo rollback delete job
            return err
      }

      err = tx.Model(&app).UpdateColumn("BuildParams", app.BuildParams).Error
      if err != nil {
            global.GVA_LOG.Error("create jenkins job save buildParams err", zap.Any("app", app.AppName), zap.Error(err))
            // todo rollback delete job
            return err
      }
      global.GVA_LOG.Info("create jenkins job success", zap.Any("app", app.AppName))
      return
}

func (app *App) AfterUpdate(tx *gorm.DB) (err error) {
      jenkinsService := GetJenkinsService()
      var project Project
      if err := tx.Model(&Project{}).Where("id = ?", app.ProjectId).First(&project).Error; err != nil {
            global.GVA_LOG.Error("update jenkins job, get app.Project error", zap.Any("app", app.AppName), zap.Error(err))
            return err
      }
      app.Project = project

      err = jenkinsService.UpdateJenkinsJobConfig(context.Background(), app)
      if err != nil {
            global.GVA_LOG.Error("update jenkins job error", zap.Any("app", app.AppName), zap.Error(err))
            return err
      }
      global.GVA_LOG.Info("update jenkins job success", zap.Any("app", app.AppName))
      return
}


func (app *App) BeforeDelete(tx *gorm.DB) (err error) {
      jenkinsService := GetJenkinsService()
      if err = tx.Model(&App{}).Preload("Project").Where("id = ?", app.ID).First(&app).Error; err != nil {
            global.GVA_LOG.Error("delete jenkins job, get app error", zap.Any("app", app.ID), zap.Error(err))
            return err
      }

      _, err = jenkinsService.DeleteJenkinsJob(context.Background(), app)
      if err != nil {
            if strings.Contains(err.Error(), "404") {
                  global.GVA_LOG.Info("delete jenkins job that not exist", zap.Any("app", app.AppName))
                  return nil
            }
            global.GVA_LOG.Error("delete jenkins job error", zap.Any("app", app.AppName))
            return err
      }
      global.GVA_LOG.Info("delete jenkins job success", zap.Any("app", app.AppName))
      return nil
}
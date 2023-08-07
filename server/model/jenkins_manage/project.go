// 自动生成模板Project
package jenkins_manage

import (
      "context"
      "go.uber.org/zap"
      "gorm.io/gorm"
      "jenkins-wrapper-ci/global"
      "jenkins-wrapper-ci/model/system"
      "strings"
)

// Project 结构体
type Project struct {
      global.GVA_MODEL
      ProjectName string            `json:"name" form:"name" gorm:"column:name;uniqueIndex:project_unique_index;comment:英文名称;"`
      NameCn      string            `json:"nameCn" form:"naneCn" gorm:"column:name_cn;comment:中文名;"`
      Description string             `json:"description" form:"descripton" gorm:"column:description;comment:项目描述;"`
      // 企业微信机器人, 用于发送通知
      WeWorkWebHook string           `json:"weworkWebhook" form:"weworkWebhook" gorm:"column:wework_webhook;default:null;comment:企业微信机器人;"`
      Managers    []system.SysUser   `json:"managers" form:"managers" gorm:"many2many:project_managers;comment:项目管理人员"`
      Members     []system.SysUser	 `json:"members"  form:"members" gorm:"many2many:project_members;comment:普通成员"`
      Apps        []App              `json:"apps"`
      CreatedBy   uint               `json:"createdBy" form:"createdBy" gorm:"column:created_by;comment:创建者"`
      UpdatedBy   uint               `json:"updateBy"  form:"updatedBy" gorm:"column:updated_by;comment:更新者"`
      DeletedBy   uint               `json:"DeletedBy" form:"deletedBy" gorm:"column:deleted_by;comment:删除者"`
}


// TableName Project 表名
func (Project) TableName() string {
  return "project"
}

func (project *Project) AfterCreate(tx *gorm.DB) (err error) {
      jenkinsService := GetJenkinsService()
      _, err = jenkinsService.GetOrCreateJenkinsFolder(context.Background(), project.ProjectName)
      if err != nil {
            global.GVA_LOG.Error("create jenkins folder error", zap.Any("folder", project.ProjectName))
            return err
      }
      global.GVA_LOG.Info("create jenkins folder success", zap.Any("folder", project.ProjectName))
      return
}

func (project *Project) BeforeDelete(tx *gorm.DB) (err error) {
      if err = tx.Model(&Project{}).Where("id = ?", project.ID).First(&project).Error; err != nil {
            global.GVA_LOG.Error("delete jenkins folder get project error", zap.Any("project", project.ID), zap.Error(err))
            return err
      }

      jenkinsService := GetJenkinsService()
      err = jenkinsService.DeleteJenkinsFolder(context.Background(), project.ProjectName)
      if err != nil {
            if strings.Contains(err.Error(), "404") {
                  global.GVA_LOG.Info("delete jenkins folder that not exist", zap.Any("folder", project.ProjectName))
                  return nil
            }

            global.GVA_LOG.Error("delete jenkins folder error", zap.Any("folder", project.ProjectName))
            return err
      }
      global.GVA_LOG.Info("delete jenkins folder success", zap.Any("folder", project.ProjectName))
      return
}

// 自动生成模板Build
package jenkins_manage

import (
      "context"
      "encoding/json"
      "fmt"
      "go.uber.org/zap"
      "gorm.io/datatypes"
      "gorm.io/gorm"
      "jenkins-wrapper-ci/global"
      "strings"
      "time"
)

var BuildResultMap = map[string]int {
      "SUCCESS": 1,
      "FAILIED": 2,
      "BUILDING": 3,
      "UNKNOWN": 4,
}

var ApproveStatusMap = map[string]int {
      "待审批": 1,
      "已批准": 2,
      "已拒绝": 3,
}

type ApproveBuildReq struct {
      ID            uint `json:"id"`
      ApproveStatus int  `json:"approveStatus"`
}

// Build 结构体
type Build struct {
      global.GVA_MODEL
      ProjectId              uint               `json:"projectId" form:"projectId" gorm:"column:project_id;comment:项目Id;"`
      ProjectName            string             `json:"projectName" form:"projectName" gorm:"column:project_name;comment:项目名称;"`
      AppId                  uint               `json:"appId" form:"appId" gorm:"column:app_id;comment:应用Id;"`
      AppName                string             `json:"appName" form:"appName" gorm:"column:app_name;comment:应用名称;"`
      BuildParamValues       datatypes.JSON     `json:"buildParamValues" form:"buildParamValues" gorm:"column:build_param_values;comment:构建参数"`
      BuildInfo              string             `json:"buildInfo" form:"buildInfo" gorm:"column:build_info;comment:发布内容"`
      BuildNumber            uint               `json:"buildNumber" form:"buildNumber" gorm:"column:build_number;comment:jenkins_build_number"`
      GitCommit              string             `json:"gitCommit" form:"gitCommit" gorm:"column:git_commit;comment:git_commit_id;"`
      Image                  string             `json:"image" form:"image" gorm:"column:image;comment:镜像url;"`
      ApproveStatus          int                `json:"approveStatus" form:"approveStatus" gorm:"column:approve_status;default:0;comment:流程状态;"`
      Result                 int                `json:"result" form:"result" gorm:"column:result;default:0;comment:发布状态;"`
      Duration               float64            `json:"duration" form:"duration" gorm:"column:duration;default:0;comment:发布时间;"`
      Changes                string             `json:"changes" form:"changes" gorm:"column:changes;type:text;comment:本次变更列表;"`
      Log                    string             `json:"log" form:"log" gorm:"column:log;type:text;comment:构建日志;"`
      BuildAt                *time.Time          `json:"buildAt" form:"buildAt" gorm:"column:build_at;comment:构建时间"`
      ApproveBy              uint               `json:"approveBy" form:"approveBy" grom:"column:approve_by;comment:审核者"`
      CreatedBy              uint               `json:"createdBy" form:"createdBy" gorm:"column:created_by;comment:创建者"`
      UpdatedBy              uint               `json:"updateBy"  form:"updatedBy" gorm:"column:updated_by;comment:更新者"`
}

// TableName Build 表名
func (Build) TableName() string {
  return "build"
}

func (b *Build) String() string {
      return fmt.Sprintf("build[%s-%s]", b.AppName, b.BuildNumber)
}


func (b *Build) AfterFind(tx *gorm.DB) (err error) {
      // only for find by ID
      // fetch log for running build
      fetchlog, _ := tx.Statement.Get("fetch_log")
      fetchLog, _ := fetchlog.(bool)

      if(fetchLog && b.Result == BuildResultMap["BUILDING"]) {
            jenkinsService := GetJenkinsService()
            log, err := jenkinsService.GetJenkinsJobBuildConsole(context.Background(), b.ProjectName, b.AppName, int64(b.BuildNumber))
            if err != nil {
                  // if in quiet
                  if strings.Contains(err.Error(), "404") {
                        return nil
                  }
                  global.GVA_LOG.Error("get build log error", zap.Any("build", b.ID), zap.Error(err))
                  return err
            }

            b.Log = log
            if b.Result == BuildResultMap["BUILDING"] {
                  // jenkins build 静默期默认为5s
                  // jenkins duration 为毫秒
                  b.Duration = time.Since(b.BuildAt.Add(-5*time.Second)).Seconds() * 1000
            }
      }
      return

}

func (b *Build) AfterUpdate(tx *gorm.DB) (err error) {
      // 审核成功后
      isapproveupdate, _ := tx.Statement.Get("approve_update")
      isApproveUpdate, _ := isapproveupdate.(bool)

      if isApproveUpdate {
            approvebuildreq, _ := tx.Statement.Get("approveBuildReq")
            approvebuildreqS, _ := approvebuildreq.([]byte)
            var approveBuildReq ApproveBuildReq
            err = json.Unmarshal(approvebuildreqS, &approveBuildReq)
            if err != nil {
                  global.GVA_LOG.Error("approve build get approveBuildReq error", zap.Error(err))
                  return err
            }

            var build Build
            if err := tx.Model(&Build{}).Where("id = ?", approveBuildReq.ID).First(&build).Error; err != nil {
                  global.GVA_LOG.Error("approve build get build error", zap.Any("build", build.ID), zap.Error(err))
                  return err
            }
            global.GVA_LOG.Debug("build approved", zap.Any("build", build.ID))
            if(build.ApproveStatus == ApproveStatusMap["已批准"]) {
                  jenkinsService := GetJenkinsService()
                  JenkinsBuildId, err := jenkinsService.CreateJenkinsJobBuild(context.Background(), build)
                  if err != nil {
                        global.GVA_LOG.Error("start jenkins build error", zap.Any("projectName", build.ProjectName), zap.Any("appName", build.AppName), zap.Error(err))
                        return err
                  }
                  global.GVA_LOG.Info("start jenkins build success", zap.Any("projectName", build.ProjectName), zap.Any("appName", build.AppName), zap.Any("jenkinsBuildId", JenkinsBuildId))
                  build.BuildNumber = uint(JenkinsBuildId)
                  build.Result = BuildResultMap["BUILDING"]
                  err = tx.Model(&Build{}).Where("id = ?", build.ID).Save(&build).Error
                  return err
            } else if(build.ApproveStatus == ApproveStatusMap["已拒绝"]) {
                  build.Result = BuildResultMap["FAILIED"]
                  err = tx.Model(&Build{}).Where("id = ?", build.ID).Save(&build).Error
                  return err
            }
      }
      return nil
}

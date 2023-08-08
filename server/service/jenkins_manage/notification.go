package jenkins_manage

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"jenkins-wrapper-ci/global"
	"jenkins-wrapper-ci/model/jenkins_manage"
	"net/http"
	"time"
)

var projectService ProjectService
var appService AppService

type WeworkMarkdownMessage struct {
	MsgType	string			`json:"msgtype"`
	Markdown Markdown 		`json:"markdown"`
}

type Markdown struct {
	Content string 			`json:"content"`
}


type NotificationService struct {
}

func (notificationService *NotificationService) ApproveBuildNotification(buildId uint) (error)  {
	var build jenkins_manage.Build
	if err := global.GVA_DB.Model(&build).Where("id = ?", buildId).Find(&build).Error; err != nil {
		return err
	}
	weworkWebHook, err := projectService.GetWeWorkWebHook(build.ProjectId)
	if err != nil {
		return err
	}
	if weworkWebHook == "" {
		global.GVA_LOG.Error("project WeWorkWebHook did not config", zap.Any("project", build.ProjectName))
		return fmt.Errorf("project WeWorkWebHook did not config")
	}

	gitRepo, err := appService.GetGitRepo(build.AppId)
	if err != nil {
		global.GVA_LOG.Error("get gitRepo of build err", zap.Error(err))
		return err
	}

	var messageContent = fmt.Sprintf("## [发布申请#%d]\n您有新的发布申请,请项目管理员及时审批\n项目组: **%s**\n应用: **%s**\n发布参数: %s\ngit仓库: %s\n发布内容: <font color=\"comment\">%s</font>",
		buildId, build.ProjectName, build.AppName, build.BuildParamValues.String(), gitRepo, build.BuildInfo)

	markdownMessage := WeworkMarkdownMessage{
		MsgType:  "markdown",
		Markdown: Markdown{
			Content: messageContent,
		},
	}

	message, err := json.Marshal(markdownMessage)
	if err != nil {
		return err
	}

	tr := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: &tr,
		Timeout:   15* time.Second,
	}
	_, err = client.Post(weworkWebHook, "application/json", bytes.NewReader(message))
	if err != nil {
		global.GVA_LOG.Error("send ApproveBuildNotification err", zap.Any("project", build.ProjectName), zap.Any("app", build.AppName),  zap.Any("project.WeWorkWebHook", weworkWebHook), zap.Error(err))
	}
	return err
}


func (notificationService *NotificationService)BuildDoneNotification(build jenkins_manage.Build) (error) {
	weworkWebHook, err := projectService.GetWeWorkWebHook(build.ProjectId)
	if err != nil {
		return err
	}
	if weworkWebHook == "" {
		global.GVA_LOG.Error("project WeWorkWebHook did not config", zap.Any("project", build.ProjectName))
		return fmt.Errorf("project WeWorkWebHook did not config")
	}

	gitRepo, err := appService.GetGitRepo(build.AppId)
	if err != nil {
		global.GVA_LOG.Error("get gitRepo of build err", zap.Error(err))
		return err
	}

	var messageContent string
	if build.Result == jenkins_manage.BuildResultMap["SUCCESS"] {
		messageContent = fmt.Sprintf("#### <font color=\"info\">[发布成功]</font>\n项目组: **%s**\n应用: **%s**\n发布参数: %s\ngit仓库: %s\n发布内容: <font color=\"comment\">%s</font>",
			build.ProjectName, build.AppName, build.BuildParamValues.String(), gitRepo, build.BuildInfo)
	} else if build.Result == jenkins_manage.BuildResultMap["FAILIED"] {
		messageContent = fmt.Sprintf("#### <font color=\"warning\">[发布失败]</font>\n项目组: **%s**\n应用: **%s**\n发布参数: %s\ngit仓库: %s\n发布内容: <font color=\"comment\">%s</font>",
			build.ProjectName, build.AppName, build.BuildParamValues.String(), gitRepo, build.BuildInfo)
	} else {
		global.GVA_LOG.Warn("send BuildDoneNotification err, unknown build result")
		return nil
	}

	markdownMessage := WeworkMarkdownMessage{
		MsgType:  "markdown",
		Markdown: Markdown{
			Content: messageContent,
		},
	}

	message, err := json.Marshal(markdownMessage)
	if err != nil {
		return err
	}

	tr := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: &tr,
		Timeout:   15* time.Second,
	}
	_, err = client.Post(weworkWebHook, "application/json", bytes.NewReader(message))
	if err != nil {
		global.GVA_LOG.Error("send BuildDoneNotification err", zap.Any("project", build.ProjectName), zap.Any("app", build.AppName),  zap.Any("project.WeWorkWebHook", weworkWebHook), zap.Error(err))
	}
	return err
}

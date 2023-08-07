package jenkins_manage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bndr/gojenkins"
	"go.uber.org/zap"
	"jenkins-wrapper-ci/global"
	"jenkins-wrapper-ci/utils"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
)


type JenkinsView *gojenkins.View
type JenkinsJob *gojenkins.Job
type JenkinsBuild *gojenkins.Build
type JenkinsBuildParams []gojenkins.ParameterDefinition


const JobTemplatePath  = "resource/jenkins_template/job.xml.tpl"

var (
	jenkinsService	*JenkinsService
	once 	sync.Once
)

func GetJenkinsService() *JenkinsService {
	once.Do(func() {
		s, err := NewJenkinsService(nil, global.GVA_CONFIG.JenkinsConfig.Url, global.GVA_CONFIG.JenkinsConfig.UserName, global.GVA_CONFIG.JenkinsConfig.Password, global.GVA_CONFIG.JenkinsConfig.Timeout)
		if err != nil {
			global.GVA_LOG.Error("init jenkins client error", zap.Error(err))
			panic(err)
		}
		// test plugins
		_, err = s.JCli.HasPlugin(context.Background(), "folder")
		if err != nil {
			global.GVA_LOG.Error("jenkins server did not install folder plugin")
			panic(err)
		}
		jenkinsService = s
	})
	return jenkinsService
}


type JenkinsService struct {
	JCli *gojenkins.Jenkins
}

func NewJenkinsService(client *http.Client, url, username, password string, timeout uint) (*JenkinsService, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	jenkins, err := gojenkins.CreateJenkins(client, url, username, password).Init(ctx)
	if err != nil {
		return &JenkinsService{}, err
	}
	return &JenkinsService{
		JCli: jenkins,
	}, nil
}

func (JenkinsService *JenkinsService) GenJenkinsJobConfig(app *App) (string, error) {
	if(app.CustomConfig != "") {
		// config描述中包含中文字符, 转为10进制unicode, format: &#xxxx;
		// 只转换description中的中文
		descPattern := regexp.MustCompile(`<description>(.*)</description>`)
		descMatched := descPattern.FindAllString(app.CustomConfig, 1)
		var desc string
		for _, r := range descMatched[0] {
			desc += utils.ChineseToHTMLEntity(r)
		}
		config  := descPattern.ReplaceAll([]byte(app.CustomConfig), []byte(desc))
		return string(config), nil
	}

	jobTemplate, err := template.ParseFiles(JobTemplatePath)
	if err != nil {
		return "", fmt.Errorf("parse jenkins job template err:", err.Error())
	}
	cfg := strings.Builder{}
	jobConfigData, err := app.GetJenkinsJobConfigData()
	if err != nil {
		return "", err
	}
	err = jobTemplate.Execute(&cfg, jobConfigData)
	if err != nil {
		return "", fmt.Errorf("render jenkins job template err:", err.Error())
	}

	//config = strings.ToValidUTF8(strings.ReplaceAll(cfg.String(), string(rune(0x8b)), ""), "")
	//exist := strings.ContainsRune(config, 0x8b)

	return cfg.String(), nil
}

func (jenkinsService *JenkinsService) GetJenkinsFolder(ctx context.Context, name string) (jenkinsFolder *gojenkins.Folder, err error) {
	return jenkinsService.JCli.GetFolder(ctx, name)
}

func (jenkinsService *JenkinsService) CreateJenkinsFolder(ctx context.Context, name string) (jenkinsFolder *gojenkins.Folder, err error) {
	return jenkinsService.JCli.CreateFolder(ctx, name)
}


func (jenkinsService *JenkinsService) GetOrCreateJenkinsFolder(ctx context.Context, name string) (jenkinsFolder *gojenkins.Folder, err error) {
	jenkinsFolder, err = jenkinsService.JCli.GetFolder(ctx, name)
	if err != nil {
		jenkinsFolder, err = jenkinsService.JCli.CreateFolder(ctx, name)
		return jenkinsFolder, err
	}
	return jenkinsFolder, err
}

func (jenkinsService *JenkinsService) DeleteJenkinsFolder(ctx context.Context, name string) (err error) {
	// folder is a job
	folder := gojenkins.Job{
		Jenkins: jenkinsService.JCli,
		Raw: new(gojenkins.JobResponse),
		Base: "/job/" + name,
	}

	_, err = folder.Delete(ctx)
	return err
}



func (jenkinsService *JenkinsService) GetFolderJobObj(folderName, jobName string) (gojenkins.Job) {
	folders := []string{folderName}
	return  gojenkins.Job{
		Jenkins: jenkinsService.JCli,
		Raw: new(gojenkins.JobResponse),
		Base: "/job/" + strings.Join(append(folders, jobName), "/job/")}
}

func (jenkinsService *JenkinsService) CreateJenkinsJob(ctx context.Context, app *App) (jenkinsJob *gojenkins.Job, err error) {
	folder, err := jenkinsService.GetOrCreateJenkinsFolder(ctx, app.Project.ProjectName)
	if err != nil {
		return jenkinsJob, err
	}

	config, err := jenkinsService.GenJenkinsJobConfig(app)
	if err != nil {
		return jenkinsJob, err
	}

	folderJob := jenkinsService.GetFolderJobObj(folder.GetName(), app.AppName)

	qr := map[string]string{
		"name": app.AppName,
	}

	jenkinsJob, err = folderJob.Create(ctx, config, qr)
	if err != nil {
		return jenkinsJob, err
	}
	return jenkinsJob, nil
}

func (jenkinsService *JenkinsService) UpdateJenkinsJobConfig(ctx context.Context, app *App) (err error) {
	folder, err := jenkinsService.GetOrCreateJenkinsFolder(ctx, app.Project.ProjectName)
	if err != nil {
		return err
	}

	config, err := jenkinsService.GenJenkinsJobConfig(app)
	if err != nil {
		return err
	}

	folderJob := jenkinsService.GetFolderJobObj(folder.GetName(), app.AppName)
	_, err = folderJob.Poll(ctx)
	if err != nil {
		return err
	}
	err = folderJob.UpdateConfig(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

func (jenkinsService *JenkinsService) DeleteJenkinsJob(ctx context.Context, app *App) (success bool, err error) {
	folder, err := jenkinsService.GetOrCreateJenkinsFolder(ctx, app.Project.ProjectName)
	if err != nil {
		return success, err
	}
	folderJob := jenkinsService.GetFolderJobObj(folder.GetName(), app.AppName)

	success, err = folderJob.Delete(ctx)
	if err != nil {
		return success, err
	}
	return success, nil
}

func (jenkinsService *JenkinsService) CreateJenkinsJobBuild(ctx context.Context, build Build) (buildId int64, err error) {
	var buildParamValues BuildParamValues
	err = json.Unmarshal(build.BuildParamValues, &buildParamValues)
	if err != nil {
		return buildId, err
	}
	folderJob := jenkinsService.GetFolderJobObj(build.ProjectName, build.AppName)
	_, err = folderJob.Poll(ctx)
	if err != nil {
		return buildId, err
	}
	jobDetails := folderJob.GetDetails()
	nextBuildNumber := jobDetails.NextBuildNumber

	_, err = folderJob.InvokeSimple(ctx, buildParamValues)
	if err != nil {
		// log
		return buildId, err
	}
	return nextBuildNumber, nil

	// pull the build number from queued task
	//task, err := jenkinsService.JCli.GetQueueItem(ctx, queueId)
	//if err != nil {
	//	//	log,
	//	return buildId, err
	//}
	//for (task.Raw.Executable.Number == 0) {
	//	_, _ = task.Poll(ctx)
	//}
	//return task.Raw.Executable.Number, err


}

func (jenkinsService *JenkinsService) GetJenkinsJobBuild(ctx context.Context, app *App, buildId int64) (jenkinsBuild *gojenkins.Build, err error) {
	folderJob := jenkinsService.GetFolderJobObj(app.Project.ProjectName, app.AppName)
	_, err = folderJob.Poll(ctx)
	if err != nil {
		return jenkinsBuild, err
	}

	jenkinsBuild, err = folderJob.GetBuild(ctx, buildId)
	if err != nil {
		return jenkinsBuild, err
	}
	return jenkinsBuild, nil
}

func (jenkinsService *JenkinsService) GetJenkinsJobBuilds(ctx context.Context, app *App) (jenkinsBuilds []gojenkins.JobBuild, err error) {
	folderJob := jenkinsService.GetFolderJobObj(app.Project.ProjectName, app.AppName)

	jenkinsBuilds, err = folderJob.GetAllBuildIds(ctx)
	if err != nil {
		return jenkinsBuilds, err
	}
	return jenkinsBuilds, nil
}

func (jenkinsService *JenkinsService) GetJenkinsJobLastBuild(ctx context.Context, app *App) (jenkinsBuild *gojenkins.Build, err error) {
	folderJob := jenkinsService.GetFolderJobObj(app.Project.ProjectName, app.AppName)

	jenkinsBuild, err = folderJob.GetLastBuild(ctx)
	if err != nil {
		return jenkinsBuild, err
	}
	return jenkinsBuild, nil
}

func (jenkinsService *JenkinsService) GetJenkinsJobFirstBuild(ctx context.Context, app *App) (jenkinsBuild *gojenkins.Build, err error) {
	folderJob := jenkinsService.GetFolderJobObj(app.Project.ProjectName, app.AppName)

	jenkinsBuild, err = folderJob.GetFirstBuild(ctx)
	if err != nil {
		return jenkinsBuild, err
	}
	return jenkinsBuild, nil
}


func (jenkinsService *JenkinsService) GetJenkinsJobBuildConsole(ctx context.Context, projectName,  appName string, buildId int64) (output string, err error) {
	folderJob := jenkinsService.GetFolderJobObj(projectName, appName)
	_, err = folderJob.Poll(ctx)
	if err != nil {
		return output, err
	}
	jobBuild, err := folderJob.GetBuild(ctx, buildId)
	if err != nil {
		return output, err
	}
	return jobBuild.GetConsoleOutput(ctx), nil
}

func (jenkinsService *JenkinsService) GetJenkinsJobBuildRevision(ctx context.Context, app *App, buildId int64) (version string, err error) {
	jobBuild, err := jenkinsService.GetJenkinsJobBuild(ctx, app, buildId)
	if err != nil {
		return version, err
	}
	return jobBuild.GetRevision(), nil
}

func (jenkinsService *JenkinsService) GetJenkinsJobBuildDuration(ctx context.Context, app *App, buildId int64) (duration float64, err error) {
	jobBuild, err := jenkinsService.GetJenkinsJobBuild(ctx, app, buildId)
	if err != nil {
		return duration, err
	}
	return jobBuild.GetDuration(), nil
}

func (jenkinsService *JenkinsService) GetJenkinsJobBuildInfo(ctx context.Context, app *App, buildId int64) (buildInfo *gojenkins.BuildResponse, err error) {
	jobBuild, err := jenkinsService.GetJenkinsJobBuild(ctx, app, buildId)
	if err != nil {
		return buildInfo, err
	}
	return jobBuild.Info(), nil
}


func (jenkinsService *JenkinsService) GetJenkinsJobConfig(ctx context.Context, app *App) (config string, err error) {
	folderJob := jenkinsService.GetFolderJobObj(app.Project.ProjectName, app.AppName)
	config, err  = folderJob.GetConfig(ctx)
	if err != nil {
		//	log
		return config, err
	}
	return config, nil
}

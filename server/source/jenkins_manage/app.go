package jenkins_manage

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"jenkins-wrapper-ci/global"
	jenkinsManageModel "jenkins-wrapper-ci/model/jenkins_manage"
	"jenkins-wrapper-ci/service/system"
	"time"
)

const initOrderApp = initOrderProject + 1

type initApp struct{}

// auto run
func init() {
	system.RegisterInit(initOrderApp, &initApp{})
}

func (i *initApp) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&jenkinsManageModel.App{})
}

func (i *initApp) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&jenkinsManageModel.App{})
}

func (i initApp) InitializerName() string {
	return jenkinsManageModel.App{}.TableName()
}

func (i *initApp) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	projects, ok := ctx.Value(jenkinsManageModel.Project{}.TableName()).([]jenkinsManageModel.Project)
	entities := []jenkinsManageModel.App{
		{
			GVA_MODEL: global.GVA_MODEL{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			AppName:   "test_proj_app1",
			Description: "测试项目应用1",
			GitRepo:     "github.com/jenkins_wrapper_ci/app1.git",
			Image:       "jenkins_wrapper_ci/app1",
			BuildParams: jenkinsManageModel.DefaultBuildParams,
			CreatedBy: 1,
			UpdatedBy: 1,
			DeletedBy: 0,
			ProjectId: projects[0].ID,
		},
		{
			GVA_MODEL: global.GVA_MODEL{
				ID:        2,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			AppName:      "test_proj_app2",
			Description: "测试项目应用2",
			Image: "jenkins_wrapper_ci/app2",
			GitRepo: "github.com/jenkins_wrapper_ci/app2.git",
			BuildParams: jenkinsManageModel.DefaultBuildParams,
			CreatedBy: 1,
			UpdatedBy: 1,
			ProjectId: projects[1].ID,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, jenkinsManageModel.App{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, err
}

func (i *initApp) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var record jenkinsManageModel.App
	if errors.Is(db.Where("name = ?", "test_proj_app1").
		Preload("Project").First(&record).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return record.AppName == "test_proj_app1"
}

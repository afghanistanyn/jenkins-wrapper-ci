package jenkins_manage

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"jenkins-wrapper-ci/global"
	jenkinsManageModel "jenkins-wrapper-ci/model/jenkins_manage"
	sysModel "jenkins-wrapper-ci/model/system"
	"jenkins-wrapper-ci/service/system"
	"time"
)

const initOrderProject = system.InitOrderJenkinsManage + 1

type initProject struct{}

// auto run
func init() {
	system.RegisterInit(initOrderProject, &initProject{})
}

func (i *initProject) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&jenkinsManageModel.Project{})
}

func (i *initProject) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&jenkinsManageModel.Project{})
}

func (i initProject) InitializerName() string {
	return jenkinsManageModel.Project{}.TableName()
}

func (i *initProject) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []jenkinsManageModel.Project{
		{
			GVA_MODEL: global.GVA_MODEL{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ProjectName: "test_proj",
			NameCn:    "测试项目",
			Description: "测试项目",
			CreatedBy: 1,
			UpdatedBy: 1,
		},
		{
			GVA_MODEL: global.GVA_MODEL{
				ID:        2,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ProjectName:   "test_proj_2",
			NameCn:    "测试项目2",
			Description: "测试项目2",
			CreatedBy: 1,
			UpdatedBy: 1,
		},
		{
			GVA_MODEL: global.GVA_MODEL{
				ID:        3,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ProjectName:  "test_proj_3",
			NameCn:    "测试项目3",
			Description: "测试项目3",
			CreatedBy: 1,
			UpdatedBy: 1,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, jenkinsManageModel.Project{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	users, ok := ctx.Value(sysModel.SysUser{}.TableName()).([]sysModel.SysUser)
	if !ok {
		return next, errors.Wrap(system.ErrMissingDependentContext, "创建 [项目-成员] 关联失败, 未找到用户表初始化数据")
	}
	// set members
	if err = db.Model(&entities[0]).Association("Members").Replace(users[3:4]); err != nil {
		return next, err
	}
	if err = db.Model(&entities[0]).Association("Managers").Replace(users[2:3]); err != nil {
		return next, err
	}
	if err = db.Model(&entities[1]).Association("Members").Replace(users[3:4]); err != nil {
		return next, err
	}
	if err = db.Model(&entities[1]).Association("Managers").Replace(users[2:3]); err != nil {
		return next, err
	}
	if err = db.Model(&entities[2]).Association("Managers").Replace(users[3:4]); err != nil {
		return next, err
	}

	return next, err
}

func (i *initProject) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var record jenkinsManageModel.Project
	if errors.Is(db.Where("name = ?", "test_proj").
		Preload("Managers").Preload("Members").First(&record).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return record.ProjectName == "test_proj"
}

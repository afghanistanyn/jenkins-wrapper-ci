package system

import (
	"context"
	sysModel "jenkins-wrapper-ci/model/system"
	"jenkins-wrapper-ci/service/system"
	"jenkins-wrapper-ci/utils"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderUser = initOrderAuthority + 1

type initUser struct{}

// auto run
func init() {
	system.RegisterInit(initOrderUser, &initUser{})
}

func (i *initUser) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysUser{})
}

func (i *initUser) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysUser{})
}

func (i initUser) InitializerName() string {
	return sysModel.SysUser{}.TableName()
}

func (i *initUser) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	password := utils.BcryptHash("123456")

	entities := []sysModel.SysUser{
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "admin",
			Password:    password,
			NickName:    "管理员",
			HeaderImg:   "https://qmplusimg.henrongyi.top/gva_header.jpg",
			AuthorityId: 888,
			Phone:       "18502073575",
			Email:       "333333333@qq.com",
		},
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "ops_1",
			Password:    password,
			NickName:    "运维用户1",
			HeaderImg:   "https:///qmplusimg.henrongyi.top/1572075907logo.png",
			AuthorityId: 999,
			Phone:       "17611111111",
			Email:       "333333333@qq.com",
		},
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "project_manager_1",
			Password:    password,
			NickName:    "项目管理员1",
			HeaderImg:   "https:///qmplusimg.henrongyi.top/1572075907logo.png",
			AuthorityId: 1000,
			Phone:       "17611111111",
			Email:       "333333333@qq.com",
		},
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "dev_1",
			Password:    password,
			NickName:    "开发者1",
			HeaderImg:   "https:///qmplusimg.henrongyi.top/1572075907logo.png",
			AuthorityId: 1000,
			Phone:       "17611111111",
			Email:       "333333333@qq.com",
		},

	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysUser{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	authorityEntities, ok := ctx.Value(initAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return next, errors.Wrap(system.ErrMissingDependentContext, "创建 [用户-权限] 关联失败, 未找到权限表初始化数据")
	}
	if err = db.Model(&entities[0]).Association("Authorities").Replace(authorityEntities[0:1]); err != nil {
		return next, err
	}
	if err = db.Model(&entities[1]).Association("Authorities").Replace(authorityEntities[1:2]); err != nil {
		return next, err
	}
	// assign ops
	if err = db.Model(&entities[2]).Association("Authorities").Replace(authorityEntities[2:3]); err != nil {
		return next, err
	}
	// assign dev
	if err = db.Model(&entities[3]).Association("Authorities").Replace(authorityEntities[2:3]); err != nil {
		return next, err
	}

	return next, err
}

func (i *initUser) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var record sysModel.SysUser
	if errors.Is(db.Where("username = ?", "admin").
		Preload("Authorities").First(&record).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return len(record.Authorities) > 0 && record.Authorities[0].AuthorityId == 888
}

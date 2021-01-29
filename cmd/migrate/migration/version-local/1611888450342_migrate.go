package version_local

import (
	"fiy/app/cmdb/models/model"
	"runtime"
	"time"

	"gorm.io/gorm"

	"fiy/cmd/migrate/migration"
	common "fiy/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1611888450342Migrate)
}

func _1611888450342Migrate(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		relatedTypeList := []model.RelatedType{
			{Id: 1, Identifies: "belong", Name: "属于", SourceDescribe: "属于", TargetDescribe: "包含", Direction: 1, BaseModel: common.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{Id: 2, Identifies: "group", Name: "组成", SourceDescribe: "组成", TargetDescribe: "组成于", Direction: 1, BaseModel: common.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{Id: 3, Identifies: "bk_mainline", Name: "拓扑组成", SourceDescribe: "组成", TargetDescribe: "组成于", Direction: 1, BaseModel: common.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{Id: 4, Identifies: "run", Name: "运行", SourceDescribe: "运行于", TargetDescribe: "运行", Direction: 1, BaseModel: common.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{Id: 5, Identifies: "connect", Name: "上联", SourceDescribe: "上联", TargetDescribe: "下联", Direction: 1, BaseModel: common.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{Id: 6, Identifies: "default", Name: "默认关联", SourceDescribe: "关联", TargetDescribe: "关联", Direction: 1, BaseModel: common.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		}

		err := tx.Create(relatedTypeList).Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}

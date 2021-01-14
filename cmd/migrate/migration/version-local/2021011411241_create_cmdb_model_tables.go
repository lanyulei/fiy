package version_local

import (
	"runtime"

	"gorm.io/gorm"

	cmdbModelModels "fiy/app/cmdb/models/model"
	"fiy/cmd/migrate/migration"
	common "fiy/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _20210114112401_create_cmdb_model_tables)
}

func _20210114112401_create_cmdb_model_tables(db *gorm.DB, version string) error {
	err := db.Debug().Migrator().AutoMigrate(
		// Cmdb 模型相关表
		new(cmdbModelModels.FieldGroup),
		new(cmdbModelModels.Fields),
		new(cmdbModelModels.Group),
		new(cmdbModelModels.Info),
		new(cmdbModelModels.RelatedType),
		new(cmdbModelModels.InfoRelatedType),
	)
	if err != nil {
		return err
	}
	return db.Create(&common.Migration{
		Version: version,
	}).Error
}

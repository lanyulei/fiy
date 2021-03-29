package version

import (
	"gorm.io/gorm"

	"fiy/app/admin/models/system"
	"fiy/cmd/migrate/migration"
	common "fiy/common/models"

	"runtime"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1603516925109Migrate)
}

func _1603516925109Migrate(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if tx.Migrator().HasColumn(&system.SysLoginLog{}, "info_id") {
			err := tx.Migrator().RenameColumn(&system.SysLoginLog{}, "info_id", "id")
			if err != nil {
				return err
			}
		}

		if tx.Migrator().HasColumn(&system.SysOperaLog{}, "oper_id") {
			err := tx.Migrator().RenameColumn(&system.SysOperaLog{}, "oper_id", "id")
			if err != nil {
				return err
			}
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}

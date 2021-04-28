package version_local

import (
	"fiy/app/admin/models"
	"fiy/common/global"
	"runtime"

	"gorm.io/gorm"

	"fiy/cmd/migrate/migration"
	common "fiy/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1619578282843K8sMigrate)
}

func _1619578282843K8sMigrate(db *gorm.DB, version string) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {

		if err = models.InitDb(tx, "config/sql/k8s.sql"); err != nil {
			global.Logger.Errorf("同步k8s初始数据失败, %v", err)
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}

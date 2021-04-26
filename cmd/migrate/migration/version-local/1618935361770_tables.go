package version_local

import (
	"runtime"

	"gorm.io/gorm"

	k8sModelsProject "fiy/app/k8s/models/project"
	k8sModelsVersion "fiy/app/k8s/models/version"
	"fiy/cmd/migrate/migration"
	common "fiy/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1618935361770Tables)
}

func _1618935361770Tables(db *gorm.DB, version string) error {
	err := db.Debug().Migrator().AutoMigrate(
		// k8s
		new(k8sModelsProject.Info),
		new(k8sModelsVersion.Manifest),
	)
	if err != nil {
		return err
	}
	return db.Create(&common.Migration{
		Version: version,
	}).Error
}

package version_local

import (
	"runtime"

	"gorm.io/gorm"

	k8sModelsProject "fiy/app/k8s/models/project"
	k8sModelsSettings "fiy/app/k8s/models/settings"
	k8sModelsVersion "fiy/app/k8s/models/version"
	"fiy/cmd/migrate/migration"
	common "fiy/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1618935361770K8sTables)
}

func _1618935361770K8sTables(db *gorm.DB, version string) error {
	err := db.Debug().Migrator().AutoMigrate(
		// k8s
		new(k8sModelsProject.Info),
		new(k8sModelsVersion.Manifest),
		new(k8sModelsSettings.Registry),
		new(k8sModelsSettings.Credential),
		new(k8sModelsSettings.NTP),
	)
	if err != nil {
		return err
	}
	return db.Create(&common.Migration{
		Version: version,
	}).Error
}

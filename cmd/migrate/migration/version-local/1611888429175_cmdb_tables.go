package version_local

import (
	"runtime"

	"gorm.io/gorm"

	cmdbAnalysisModels "fiy/app/cmdb/models/analysis"
	cmdbBusinessModels "fiy/app/cmdb/models/business"
	cmdbModelModels "fiy/app/cmdb/models/model"
	cmdbResourceModels "fiy/app/cmdb/models/resource"
	"fiy/cmd/migrate/migration"
	common "fiy/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1611888429175CmdbTables)
}

func _1611888429175CmdbTables(db *gorm.DB, version string) error {
	err := db.Debug().Migrator().AutoMigrate(
		// Cmdb 模型相关表
		new(cmdbModelModels.FieldGroup),
		new(cmdbModelModels.Fields),
		new(cmdbModelModels.Group),
		new(cmdbModelModels.Info),
		new(cmdbModelModels.RelatedType),
		new(cmdbModelModels.InfoRelatedType),

		// Cmdb 资源
		new(cmdbResourceModels.Data),
		new(cmdbResourceModels.CloudAccount),
		new(cmdbResourceModels.CloudDiscovery),
		new(cmdbResourceModels.DataRelated),

		// Cmdb 业务
		new(cmdbBusinessModels.ServiceClassify),
		new(cmdbBusinessModels.ClusterTemplate),
		new(cmdbBusinessModels.ServiceCluster),
		new(cmdbBusinessModels.ServiceTemplate),
		new(cmdbBusinessModels.ServiceTemplateProcess),

		// Cmdb 运营分析
		new(cmdbAnalysisModels.Audit),
	)
	if err != nil {
		return err
	}
	return db.Create(&common.Migration{
		Version: version,
	}).Error
}

package database

import (
	"database/sql"
	"fiy/common/log"
	. "log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"fiy/common/config"
	"fiy/common/global"
	goAdminLogger "fiy/pkg/core/logger"
	"fiy/tools"
	toolsConfig "fiy/tools/config"
)

type PgSql struct {
}

func (e *PgSql) Setup() {
	var err error

	global.Source = e.GetConnect()
	log.Info(global.Source)
	db, err := sql.Open("postgresql", global.Source)
	if err != nil {
		global.Logger.Fatal(tools.Red(e.GetDriver()+" connect error :"), err)
	}
	global.Cfg.SetDb(&config.DBConfig{
		Driver: "mysql",
		DB:     db,
	})
	global.Eloquent, err = e.Open(db, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("%s connect error %v", e.GetDriver(), err)
	} else {
		log.Infof("%s connect success!", e.GetDriver())
	}

	if global.Eloquent.Error != nil {
		log.Fatalf("database error %v", global.Eloquent.Error)
	}

	if toolsConfig.LoggerConfig.EnabledDB {
		global.Eloquent.Logger = logger.New(
			New(goAdminLogger.DefaultLogger.Options().Out, "\r\n", LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					goAdminLogger.DefaultLogger.Options().Level.LevelForGorm()),
			})
	}
}

// 打开数据库连接
func (e *PgSql) Open(db *sql.DB, cfg *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{Conn: db}), cfg)
}

func (e *PgSql) GetConnect() string {
	return toolsConfig.DatabaseConfig.Source
}

func (e *PgSql) GetDriver() string {
	return toolsConfig.DatabaseConfig.Driver
}

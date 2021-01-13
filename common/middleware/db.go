package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"fiy/common/config"
	"fiy/common/global"
	"fiy/tools"
)

func WithContextDb(dbMap map[string]*gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db, ok := dbMap["*"]; ok {
			c.Set("db", db)
		} else {
			c.Set("db", dbMap[c.Request.Host])
		}
		c.Next()
	}
}

func GetGormFromConfig(cfg config.Conf) map[string]*gorm.DB {
	gormDB := make(map[string]*gorm.DB)
	if cfg.GetSaas() {
		var err error
		for k, v := range cfg.GetDbs() {
			gormDB[k], err = getGormFromDb(v.Driver, v.DB, &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
			if err != nil {
				global.Logger.Fatal(tools.Red(k+" connect error :"), err)
			}
		}
		return gormDB
	}
	c := cfg.GetDb()
	db, err := getGormFromDb(c.Driver, c.DB, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		global.Logger.Fatal(tools.Red(c.Driver+" connect error :"), err)
	}
	gormDB["*"] = db
	return gormDB
}

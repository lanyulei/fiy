package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"fiy/common/config"
	"fiy/pkg/logger"
)

const (
	// fiy Version Info
	Version = "1.0.0"
)

var Cfg config.Conf = config.DefaultConfig()

var GinEngine *gin.Engine
var CasbinEnforcer *casbin.SyncedEnforcer
var Eloquent *gorm.DB

var GADMCron *cron.Cron

var (
	Source string
	Driver string
	_      string
)

var (
	Logger        = &logger.Logger{}
	JobLogger     = &logger.Logger{}
	RequestLogger = &logger.Logger{}
)

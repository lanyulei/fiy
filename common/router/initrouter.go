package router

import (
	"os"

	"fiy/common/global"
	"fiy/common/log"
	"fiy/common/middleware"
	"fiy/common/middleware/handler"
	_ "fiy/pkg/jwtauth"
	"fiy/tools"
	"fiy/tools/config"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	var r *gin.Engine
	h := global.Cfg.GetEngine()
	if h == nil {
		h = gin.New()
		global.Cfg.SetEngine(h)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}
	if config.SslConfig.Enable {
		r.Use(handler.TlsHandler())
	}
	r.Use(middleware.WithContextDb(middleware.GetGormFromConfig(global.Cfg)))

	r.Use(middleware.Sentinel())
	middleware.InitMiddleware(r)
	// the jwt middleware
	var err error
	authMiddleware, err := middleware.AuthInit()
	tools.HasError(err, "JWT Init Error", 500)

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	// 注册业务路由

	//return r
}

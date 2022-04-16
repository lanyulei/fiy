package api

import (
	"context"
	"fiy/common/router"
	"fiy/pkg/sync_cloud"
	"fiy/tools/trace"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fiy/common/database"
	"fiy/common/global"
	"fiy/common/log"
	myCasbin "fiy/pkg/casbin"
	"fiy/pkg/es"
	"fiy/pkg/logger"
	"fiy/tools"
	"fiy/tools/config"
)

var (
	configYml  string
	port       string
	mode       string
	traceStart bool
	StartCmd   = &cobra.Command{
		Use:          "server",
		Short:        "Start API server",
		Example:      "fiy server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

var AppRouters = make([]func(), 0)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
	StartCmd.PersistentFlags().BoolVarP(&traceStart, "traceStart", "t", false, "start traceStart app dash")

	//注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
}

func setup() {

	//1. 读取配置
	config.Setup(configYml)
	//2. 设置日志
	global.Logger.Logger = logger.SetupLogger(config.LoggerConfig.Path, "bus")
	global.JobLogger.Logger = logger.SetupLogger(config.LoggerConfig.Path, "job")
	global.RequestLogger.Logger = logger.SetupLogger(config.LoggerConfig.Path, "request")
	//3. 初始化数据库链接
	database.Setup(config.DatabaseConfig.Driver)
	//4. 接口访问控制加载
	global.CasbinEnforcer = myCasbin.Setup(global.Eloquent, "sys_")

	usageStr := `starting api server`
	log.Info(usageStr)

}

func run() error {
	if err := global.InitTrans("zh"); err != nil {
		panic(fmt.Sprintf("init trans failed, err:%v\n", err))
	}

	if viper.GetString("settings.application.mode") == string(tools.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := global.Cfg.GetEngine()
	if engine == nil {
		engine = gin.New()
	}

	if mode == "dev" {
		//监控
		AppRouters = append(AppRouters, router.Monitor)
	}

	// Es 连接
	if viper.GetBool("settings.es.enabled") {
		es.Init()
	}

	// 同步云资源
	go func() {
		err := sync_cloud.Start()
		if err != nil {
			log.Fatal("sync cloud err: ", err)
		}
	}()

	for _, f := range AppRouters {
		f()
	}

	srv := &http.Server{
		Addr:    config.ApplicationConfig.Host + ":" + config.ApplicationConfig.Port,
		Handler: global.Cfg.GetEngine(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if traceStart {
		//链路追踪, fixme 页面显示需要自备梯子
		trace.Start()
		defer trace.Stop(ctx)
	}

	go func() {
		// 服务连接
		if config.SslConfig.Enable {
			if err := srv.ListenAndServeTLS(config.SslConfig.Pem, config.SslConfig.KeyStr); err != nil && err != http.ErrServerClosed {
				log.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("listen: ", err)
			}
		}
	}()
	tip()
	fmt.Println(tools.Green("Server run at:"))
	fmt.Printf("-  Local:   http://localhost:%s/ \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: http://%s:%s/ \r\n", tools.GetLocaHonst(), config.ApplicationConfig.Port)
	fmt.Println(tools.Green("Swagger run at:"))
	fmt.Printf("-  Local:   http://localhost:%s/swagger/index.html \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: http://%s:%s/swagger/index.html \r\n", tools.GetLocaHonst(), config.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", tools.GetCurrentTimeStr())
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", tools.GetCurrentTimeStr())

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")

	return nil
}

func tip() {
	usageStr := `欢迎使用 ` + tools.Green(`fiy `+global.Version) + ` 可以使用 ` + tools.Red(`-h`) + ` 查看命令`
	fmt.Printf(" \n%s\n\n", usageStr)
}

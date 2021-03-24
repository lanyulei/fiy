package agent

import (
	"fiy/common/global"
	"fiy/pkg/logger"
	"fiy/tools"
	"fiy/tools/config"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:          "agent",
		Short:        "Start agent",
		Example:      "fiy agent -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
}

func setup() {
	//1. 读取配置
	config.Setup(configYml)
	//2. 设置日志
	global.Logger.Logger = logger.SetupLogger(config.LoggerConfig.Path, "bus")
}

func run() (err error) {
	fmt.Println(tools.Green("Agent run at:"))
	fmt.Printf("-  Server: http://%s:%s/ \r\n", tools.GetLocaHonst(), config.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Agent \r\n", tools.GetCurrentTimeStr())
	return nil
}

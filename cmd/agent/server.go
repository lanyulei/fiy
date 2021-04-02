package agent

import (
	"fiy/common/global"
	"fiy/common/log"
	"fiy/pkg/grpc/client"
	"fiy/pkg/logger"
	"fiy/tools"
	"fiy/tools/config"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	server   string
	interval int
	uuidPath string
	StartCmd = &cobra.Command{
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
	StartCmd.PersistentFlags().StringVarP(&server, "server", "s", "localhost:50051", "Server address, example: localhost:50051")
	StartCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 5, "Resource reporting interval, unit: minutes, default: 5 minutes")
	StartCmd.PersistentFlags().StringVarP(&uuidPath, "uuid", "u", "~/.uuid", "UUID save path, default: ~/.uuid")
}

func setup() {
	global.Logger.Logger = logger.SetupLogger(config.LoggerConfig.Path, "bus")
	log.Info("初始化完成")
}

func run() (err error) {
	log.Info("启动agent...")
	fmt.Println(tools.Green("Agent run at:"))
	fmt.Printf("-  Server: %s/ \r\n", server)
	fmt.Printf("%s Enter Control + C Shutdown Agent \r\n", tools.GetCurrentTimeStr())

	c := client.NewRpcClient(uuidPath)
	c.RunClient(server, interval)

	return nil
}

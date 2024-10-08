package cmd

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/server"
	"github.com/spf13/cobra"

	"github.com/infraboard/mflow/cmd/initial"
	"github.com/infraboard/mflow/cmd/start"
)

var (
	// pusher service config option
	confType string
	confFile string
)

var vers bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mflow",
	Short: "微服务发布平台",
	Long:  "微服务发布平台",
	Run: func(cmd *cobra.Command, args []string) {
		if vers {
			fmt.Println(application.FullVersion())
			return
		}
		cmd.Help()
	},
}

func initail() {
	req := ioc.NewLoadConfigRequest()
	switch confType {
	case "file":
		req.ConfigFile.Enabled = true
		req.ConfigFile.Path = confFile
	default:
		req.ConfigEnv.Enabled = true
	}
	server.DefaultConfig = req
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 补充初始化设置
	cobra.OnInitialize(initail)
	RootCmd.AddCommand(start.Cmd)
	RootCmd.AddCommand(initial.Cmd)

	err := RootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env]")
	RootCmd.PersistentFlags().StringVarP(&confFile, "config-file", "f", "etc/config.toml", "the service config from file")
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the mflow version")
}

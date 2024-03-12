package start

import (
	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/v2/ioc/server"

	// 注册所有服务
	_ "github.com/infraboard/mcenter/clients/rpc/middleware"
	_ "github.com/infraboard/mflow/apps"
)

// Cmd represents the start command
var Cmd = &cobra.Command{
	Use:   "start",
	Short: "mflow API服务",
	Long:  "mflow API服务",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(server.Run(cmd.Context()))
	},
}

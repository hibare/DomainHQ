package cmd

import (
	"os"

	"github.com/hibare/DomainHQ/internal/api"
	"github.com/hibare/DomainHQ/internal/config"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "GoWebFinger",
	Short: "A WebFinger server implementation in Golang",
	Run: func(cmd *cobra.Command, args []string) {
		app := &api.App{}
		app.Init()
		app.Serve()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(commonLogger.InitDefaultLogger, config.LoadConfig)
}

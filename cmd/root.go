package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

var rootCmd = &cobra.Command{
	Use:   "go-sensortag",
	Short: "go-sensortag can be used to connect to TI Sensortags",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// logger setup
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		defer logger.Sync() // flushes buffer, if any
		log = logger.Sugar()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

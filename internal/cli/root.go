/*
Copyright © 2023 Mohammed Zaki zs84907@gmail.com
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zakisk/dock-stats/internal/cli/show"
	"github.com/zakisk/dock-stats/pkg/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dock-stats",
	Short: "Docker containers stats rendering tool",
	Long: `As an utility tool it renders container real-time stats more visualy 
Usage:

`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.IsDockerInstalled(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(show.ShowCmd)
}

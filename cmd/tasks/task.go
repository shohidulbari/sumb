package tasks

import (
	"github.com/spf13/cobra"
)

var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long:  `Easily create, track, update and remove tasks`,
	RunE: func(cmd *cobra.Command, args []string) error {	
		return cmd.Help()
	},
}



func init() {
	TaskCmd.AddCommand(createCmd)
	TaskCmd.AddCommand(listCmd)
	TaskCmd.AddCommand(statusCmd)
	TaskCmd.AddCommand(searchCmd)
	TaskCmd.AddCommand(deleteCmd)
	TaskCmd.AddCommand(deleteMultipleCmd)	
} 
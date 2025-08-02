package pomodoro

import (
	"github.com/spf13/cobra"
)

var PdCmd = &cobra.Command{
	Use:   "pd",
	Short: "Alias for pomodoro command",
	Long:  `Alias for the pomodoro command. Use 'pd' instead of 'pomodoro' for shorter typing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Delegate to the main pomodoro command
		return PomodoroCmd.RunE(cmd, args)
	},
}

func init() {
	// Add all the same subcommands to the alias
	PdCmd.AddCommand(startCmd)
	PdCmd.AddCommand(statusCmd)
	PdCmd.AddCommand(timerCmd)
	PdCmd.AddCommand(stopCmd)
	PdCmd.AddCommand(listCmd)
	
	// Add the same flags
	PdCmd.Flags().StringP("create", "c", "", "Quick start a pomodoro with title")
	PdCmd.Flags().IntP("session", "s", 0, "Session duration in minutes (for quick start)")
} 
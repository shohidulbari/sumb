package pomodoro

import (
	"fmt"

	"github.com/spf13/cobra"
	db "github.com/sumb/db"
)

var PomodoroCmd = &cobra.Command{
	Use:   "pomodoro",
	Short: "Manage pomodoro timers",
	Long:  `Manage your pomodoro timers with various operations like start, status, timer, stop, and list.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle quick start with -c flag for title and -s flag for session duration
		title, _ := cmd.Flags().GetString("create")
		session, _ := cmd.Flags().GetInt("session")
		
		if title != "" && session > 0 {
			pm, err := db.NewPomodoroManager()
			if err != nil {
				return fmt.Errorf("failed to initialize pomodoro manager: %w", err)
			}
			defer pm.Close()

			// Check if there's already an active pomodoro
			active, err := pm.GetActivePomodoro()
			if err != nil {
				return fmt.Errorf("failed to check active pomodoro: %w", err)
			}

			if active != nil {
				fmt.Printf("⚠️  Stopping existing pomodoro (ID: %d) to start new one...\n", active.ID)
			}

			// Create new pomodoro (this will automatically stop any existing ones)
			pomodoro, err := pm.CreatePomodoro(title, session)
			if err != nil {
				return fmt.Errorf("failed to create pomodoro: %w", err)
			}

			fmt.Printf("🍅 Pomodoro started! Title: %s\n", title)
			fmt.Printf("Duration: %d minutes\n", session)
			fmt.Printf("Started at: %s\n", pomodoro.StartedAt.Format("15:04:05"))
			fmt.Printf("Use 'sumb pomodoro status' to check remaining time\n")
			fmt.Printf("Use 'sumb pomodoro timer' to see live countdown\n")

			return nil
		}
		return cmd.Help()
	},
}

func init() {
	PomodoroCmd.AddCommand(startCmd)
	PomodoroCmd.AddCommand(statusCmd)
	PomodoroCmd.AddCommand(timerCmd)
	PomodoroCmd.AddCommand(stopCmd)
	PomodoroCmd.AddCommand(listCmd)
	
	PomodoroCmd.Flags().StringP("create", "c", "", "Quick start a pomodoro with title")
	PomodoroCmd.Flags().IntP("session", "s", 0, "Session duration in minutes (for quick start)")
} 
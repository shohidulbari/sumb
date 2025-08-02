package pomodoro

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/sumb/cmd/styles"
	db "github.com/sumb/db"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a pomodoro timer",
	Long:  `Start a pomodoro timer with specified title and duration in minutes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		duration, _ := cmd.Flags().GetInt("session")

		if title == "" {
			return fmt.Errorf("title is required")
		}

		if duration <= 0 {
			return fmt.Errorf("duration must be greater than 0")
		}

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
		pomodoro, err := pm.CreatePomodoro(title, duration)
		if err != nil {
			return fmt.Errorf("failed to create pomodoro: %w", err)
		}

		fmt.Printf("🍅 Pomodoro started! Title: %s\n", title)
		fmt.Printf("Duration: %d minutes\n", duration)
		fmt.Printf("Started at: %s\n", pomodoro.StartedAt.Format("15:04:05"))
		fmt.Printf("Use 'sumb pomodoro status' to check remaining time\n")
		fmt.Printf("Use 'sumb pomodoro timer' to see live countdown\n")

		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current pomodoro status",
	Long:  `Display the status and remaining time of the current active pomodoro.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := db.NewPomodoroManager()
		if err != nil {
			return fmt.Errorf("failed to initialize pomodoro manager: %w", err)
		}
		defer pm.Close()

		active, err := pm.GetActivePomodoro()
		if err != nil {
			return fmt.Errorf("failed to get active pomodoro: %w", err)
		}

		if active == nil {
			fmt.Println("🍅 No active pomodoro found")
			return nil
		}

		// Calculate remaining time
		endTime := active.StartedAt.Add(time.Duration(active.Duration) * time.Minute)
		remaining := time.Until(endTime)

		if remaining <= 0 {
			// Pomodoro should be completed
			if err := pm.CompletePomodoro(active.ID); err != nil {
				return fmt.Errorf("failed to complete pomodoro: %w", err)
			}
			fmt.Println("🎉 Pomodoro completed!")
			return nil
		}

		// Display status
		minutes := int(remaining.Minutes())
		seconds := int(remaining.Seconds()) % 60
		progress := float64(active.Duration*60-int(remaining.Seconds())) / float64(active.Duration*60)

		fmt.Println(styles.Separator)
		fmt.Printf("🍅 Active Pomodoro (ID: %d)\n", active.ID)
		fmt.Printf("Title: %s\n", active.Title)
		fmt.Printf("Duration: %d minutes\n", active.Duration)
		fmt.Printf("Started at: %s\n", active.StartedAt.Format("15:04:05"))
		fmt.Printf("Remaining: %02d:%02d\n", minutes, seconds)
		fmt.Printf("Progress: %.1f%%\n", progress*100)
		fmt.Println(styles.Separator)

		return nil
	},
}

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Show live pomodoro timer",
	Long:  `Display a live countdown timer for the current active pomodoro.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := db.NewPomodoroManager()
		if err != nil {
			return fmt.Errorf("failed to initialize pomodoro manager: %w", err)
		}
		defer pm.Close()

		active, err := pm.GetActivePomodoro()
		if err != nil {
			return fmt.Errorf("failed to get active pomodoro: %w", err)
		}

		if active == nil {
			return fmt.Errorf("no active pomodoro found")
		}

		// Start the timer display
		return runTimer(pm, active.ID, active.Duration)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the active pomodoro timer",
	Long:  `Stop the currently active pomodoro timer.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := db.NewPomodoroManager()
		if err != nil {
			return fmt.Errorf("failed to initialize pomodoro manager: %w", err)
		}
		defer pm.Close()

		active, err := pm.GetActivePomodoro()
		if err != nil {
			return fmt.Errorf("failed to get active pomodoro: %w", err)
		}

		if active == nil {
			return fmt.Errorf("no active pomodoro found")
		}

		if err := pm.StopPomodoro(active.ID); err != nil {
			return fmt.Errorf("failed to stop pomodoro: %w", err)
		}

		fmt.Println(styles.Separator)
		fmt.Printf("⏹️  Pomodoro stopped! ID: %d\n", active.ID)
		fmt.Printf("Title: %s\n", active.Title)
		fmt.Printf("Duration: %d minutes\n", active.Duration)
		fmt.Printf("Started at: %s\n", active.StartedAt.Format("15:04:05"))
		fmt.Printf("Stopped at: %s\n", time.Now().Format("15:04:05"))
		fmt.Println(styles.Separator)

		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List pomodoro history",
	Long:  `Display pomodoro history with their status and details. Shows max 10 latest pomodoros by default. Use --skip to paginate through older pomodoros.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skip, _ := cmd.Flags().GetInt("skip")
		
		pm, err := db.NewPomodoroManager()
		if err != nil {
			return fmt.Errorf("failed to initialize pomodoro manager: %w", err)
		}
		defer pm.Close()

		pomodoroList, err := pm.ListPomodoros(10, skip)
		if err != nil {
			return fmt.Errorf("failed to list pomodoros: %w", err)
		}

		if len(pomodoroList) == 0 {
			if skip > 0 {
				fmt.Printf("🍅 No more pomodoros found. You've reached the end of your pomodoro history.\n")
			} else {
				fmt.Println("🍅 No pomodoros found. Start your first pomodoro with 'sumb pomodoro start -t \"Work Session\" -d 25'")
			}
			return nil
		}

		fmt.Println()
		fmt.Printf("📋 Found %d pomodoro(s)", len(pomodoroList))
		if skip > 0 {
			fmt.Printf(" (skipped %d)", skip)
		}
		fmt.Printf(":\n\n")
		
		for idx, pomodoro := range pomodoroList {
			status := "⏳"
			if pomodoro.Status == "completed" {
				status = "✅"
			} else if pomodoro.Status == "stopped" {
				status = "⏹️"
			}

			if idx == 0 {
				fmt.Println(styles.Separator)
			}

			fmt.Printf("%s [%d] %s (%d minutes)\n", status, pomodoro.ID, pomodoro.Title, pomodoro.Duration)
			fmt.Printf("   Started: %s\n", pomodoro.StartedAt.Format("2006-01-02 15:04:05"))
			
			if pomodoro.CompletedAt != nil {
				fmt.Printf("   Completed: %s\n", pomodoro.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			
			fmt.Printf("   Status: %s\n", pomodoro.Status)
			fmt.Println(styles.Separator)
		}

		// Show pagination hint if there might be more pomodoros
		if len(pomodoroList) == 10 {
			fmt.Printf("💡 To see more pomodoros, use: sumb pomodoro list --skip %d\n", skip+10)
		}

		return nil
	},
}

func runTimer(pm *db.PomodoroManager, pomodoroID, duration int) error {
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Calculate end time based on the pomodoro's start time
	active, err := pm.GetActivePomodoro()
	if err != nil {
		return fmt.Errorf("failed to get active pomodoro: %w", err)
	}
	if active == nil {
		return fmt.Errorf("no active pomodoro found")
	}

	endTime := active.StartedAt.Add(time.Duration(duration) * time.Minute)
	
	// Timer loop
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("🍅 POMODORO TIMER: %s\n", active.Title)
	fmt.Println(strings.Repeat("=", 50))

	for {
		select {
		case <-ticker.C:
			remaining := time.Until(endTime)
			if remaining <= 0 {
				// Timer completed
				if err := pm.CompletePomodoro(pomodoroID); err != nil {
					return fmt.Errorf("failed to complete pomodoro: %w", err)
				}

				fmt.Println("\n" + strings.Repeat("=", 50))
				fmt.Println("🎉 POMODORO COMPLETED!")
				fmt.Println(strings.Repeat("=", 50))
				fmt.Printf("Title: %s\n", active.Title)
				fmt.Printf("Duration: %d minutes\n", duration)
				fmt.Printf("Started at: %s\n", active.StartedAt.Format("15:04:05"))
				fmt.Printf("Completed at: %s\n", time.Now().Format("15:04:05"))
				return nil
			}

			// Display countdown with bigger text
			minutes := int(remaining.Minutes())
			seconds := int(remaining.Seconds()) % 60
			
			// Clear line and show big countdown
			fmt.Printf("\r\033[K") // Clear current line
			fmt.Printf("⏰  %02d:%02d  ", minutes, seconds)
			
			// Add visual progress bar
			progress := float64(duration*60-int(remaining.Seconds())) / float64(duration*60)
			barLength := 30
			filled := int(progress * float64(barLength))
			
			fmt.Printf("[")
			for i := 0; i < barLength; i++ {
				if i < filled {
					fmt.Printf("█")
				} else {
					fmt.Printf("░")
				}
			}
			fmt.Printf("] %d%%", int(progress*100))

		case <-sigChan:
			// User interrupted timer display - just exit, don't stop pomodoro
			fmt.Println("\n\n" + strings.Repeat("=", 50))
			fmt.Println("⏹️  TIMER DISPLAY STOPPED")
			fmt.Println(strings.Repeat("=", 50))
			fmt.Println("Pomodoro is still running in the background!")
			fmt.Printf("Use 'sumb pomodoro status' to check remaining time\n")
			fmt.Printf("Use 'sumb pomodoro timer' to see live countdown again\n")
			return nil
		}
	}
}

func init() {
	startCmd.Flags().StringP("title", "t", "", "Pomodoro title (required)")
	startCmd.Flags().IntP("session", "s", 25, "Session duration in minutes (default: 25)")
	startCmd.MarkFlagRequired("title")
	
	// Add pagination flag for list command
	listCmd.Flags().IntP("skip", "s", 0, "Number of pomodoros to skip (for pagination)")
} 
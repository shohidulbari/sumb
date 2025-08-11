package notes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sumb/cmd/styles"
	db "github.com/sumb/db"
)

var NoteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes",
	Long:  `Manage your notes with various operations like create, list, and delete.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle quick create with -c flag
		body, _ := cmd.Flags().GetString("create")
		interactive, _ := cmd.Flags().GetBool("interactive")
		
		if interactive {
			return createNoteInteractive()
		}
		
		if body != "" {
			nm, err := db.NewNoteManager()
			if err != nil {
				return fmt.Errorf("failed to initialize note manager: %w", err)
			}
			defer nm.Close()

			if err := nm.CreateNote(body); err != nil {
				return fmt.Errorf("failed to create note: %w", err)
			}

			fmt.Println("--------------------------------")

			fmt.Printf("🌟 Note created!\n")
			fmt.Printf("Body: %s\n", body)
			
			fmt.Println("--------------------------------")

			return nil
		}
		return cmd.Help()
	},
}

// createNoteInteractive creates a note by reading from stdin interactively
func createNoteInteractive() error {
	fmt.Println("📝 Interactive Note Creation")
	fmt.Println("Enter your note content (press Enter twice to finish):")
	
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		
		line := scanner.Text()
		if line == "" && len(lines) > 0 && lines[len(lines)-1] == "" {
			// Two empty lines in a row - finish input
			break
		}
		
		lines = append(lines, line)
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}
	
	// Remove the last empty line and join all lines
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	
	body := strings.Join(lines, "\n")
	
	if body == "" {
		return fmt.Errorf("no note content provided")
	}
	
	nm, err := db.NewNoteManager()
	if err != nil {
		return fmt.Errorf("failed to initialize note manager: %w", err)
	}
	defer nm.Close()

	if err := nm.CreateNote(body); err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Println(styles.Separator)
	fmt.Printf("🌟 Note created!\n")
	fmt.Printf("Body:\n%s\n", body)
	fmt.Println(styles.Separator)

	return nil
}

// isJSON checks if a string is valid JSON
func isJSON(str string) bool {
	str = strings.TrimSpace(str)
	return (strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")) ||
		   (strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]"))
}

// formatJSON formats JSON with proper indentation
func formatJSON(jsonStr string) (string, error) {
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return "", err
	}
	
	formatted, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(formatted), nil
}

func init() {
	NoteCmd.AddCommand(createCmd)
	NoteCmd.AddCommand(listCmd)
	NoteCmd.AddCommand(deleteCmd)
	NoteCmd.AddCommand(deleteMultipleCmd)
	
	NoteCmd.Flags().StringP("create", "c", "", "Quick create a note with body")
	NoteCmd.Flags().BoolP("interactive", "i", false, "Create a note interactively")
} 
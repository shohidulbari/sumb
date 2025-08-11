package notes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/shohidulbari/sumb/cmd/styles"
	db "github.com/shohidulbari/sumb/db"
	"github.com/spf13/cobra"
)

var NoteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes",
	Long:  `Easily create, list, search and delete notes`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

			fmt.Println(styles.Separator)

			fmt.Printf("🌟 Note created!\n")
			fmt.Printf("%s\n", body)
			
			fmt.Println(styles.Separator)

			return nil
		}
		return cmd.Help()
	},
}

func createNoteInteractive() error {
	fmt.Println("📝 Interactive Note Creation")
	fmt.Println("Enter your note content (press Ctrl+D to finish):")
	fmt.Println(styles.Separator)
	
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			// This happens when Ctrl+D is pressed (EOF)
			break
		}
		
		line := scanner.Text()
		lines = append(lines, line)
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
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
	fmt.Printf("%s\n", body)
	fmt.Println(styles.Separator)

	return nil
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

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
	NoteCmd.AddCommand(searchCmd)
	
	NoteCmd.Flags().StringP("create", "c", "", "Quick create a note with body")
	NoteCmd.Flags().BoolP("interactive", "i", false, "Create a note interactively")
} 
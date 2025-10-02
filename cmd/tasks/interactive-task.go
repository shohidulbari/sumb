package tasks

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	db "github.com/shohidulbari/sumb/db"
)

type Task struct {
	Title string
	Description string
	Deadline string
	Status string
}

type TaskModel struct {
	form   *huh.Form
	cancelled bool
	task *Task 
}

func InteractiveTaskForm(task *Task) *TaskModel {
	currentTask := &Task{}
	if task != nil {
		currentTask = task
	}

	taskModel := &TaskModel{
		task: currentTask,
		cancelled: false,
	}
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title (required)").
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("title is required")
					}
					return nil
				}).Value(&taskModel.task.Title),

			huh.NewText().
				Key("description").
				Title("Description").
				CharLimit(200).Value(&taskModel.task.Description),

			huh.NewInput().
				Title("Deadline").
				Description("Enter deadline (YYYY-MM-DD, 'today', 'tomorrow', or press Enter to skip)").
				Placeholder("2024-01-15").
				Validate(func(str string) error {
					if str == "" || str == "today" || str == "tomorrow" {
						return nil
					}	
					_, err := time.Parse("2006-01-02", str)
					if err != nil {
						return fmt.Errorf("invalid date format. Use YYYY-MM-DD")
					}	
					return nil
				}).Value(&taskModel.task.Deadline),

			huh.NewSelect[string]().
				Title("Status").
				Description("Initial status for the task").
				Options(
					huh.NewOption("TODO", string(db.StatusTODO)),
					huh.NewOption("In Progress", string(db.StatusInProgress)),
				).
				Value(&taskModel.task.Status),
		),
	).WithTheme(huh.ThemeCatppuccin())

	taskModel.form = form

	return taskModel
}

func (m *TaskModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *TaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "ctrl+d", "ctrl+z":
			m.cancelled = true
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}	

	if m.form.State == huh.StateCompleted {
		return m, tea.Quit
	}

	return m, cmd
}

func (m *TaskModel) View() string {
	if m.form.State == huh.StateCompleted {
		return ""
	}
	return m.form.View()
}
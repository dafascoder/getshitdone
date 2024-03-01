package components

import (
	"getshitdone/db"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	Table table.Model
}

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#d0d0d0"))

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.Table.SelectedRow()[1]),
			)
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func TableColumnSturcture() table.Model {

	column := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Description", Width: 20},
		{Title: "Project", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Created At", Width: 10},
	}

	database, err := db.OpenDB()
	if err != nil {
		panic(err)
	}
	defer database.Close()

	tasks, err := database.GetTasks()
	if err != nil {
		panic(err)
	}

	var rows []table.Row
	for _, task := range tasks {
		rows = append(rows, table.Row{
			(string(task.ID)),
			task.Name,
			task.Description,
			task.Project,
			task.Status,
			task.CreatedAt.Format("2006-01-02"),
		})
	}

	t := table.New(
		table.WithColumns(column),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t

}

func (m TableModel) View() string {
	return baseStyle.Render(m.Table.View()) + "\n"
}

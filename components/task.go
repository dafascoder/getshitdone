package components

import (
	"fmt"
	"getshitdone/db"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	statusNormal state = iota
	stateDone
)

type Model struct {
	lg     *lipgloss.Renderer
	form   *huh.Form
	state  state
	width  int
	styles *Styles
}

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

const maxWidth = 80

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.Copy().
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter the Task Name").Key("name"),
			huh.NewInput().Title("Enter the Project Name").Key("project"),
			huh.NewText().Title("Enter the Task Description").Key("description"),
			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("%s", "Get your shit together man")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		t, err := db.OpenDB()
		if err != nil {
			log.Fatal("Failed to connect to DB")
		}
		defer t.Close()

		t.AddTask(m.form.GetString("name"), m.form.GetString("project"), m.form.GetString("description"))
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles

	switch m.form.State {
	case huh.StateCompleted:
		return s.Base.Render("Task created successfully!")
	case huh.StateAborted:
		return s.StatusHeader.Render("Aborted")
	default:
		var name string
		if m.form.GetString("name") != "" {
			name = "Name: " + m.form.GetString("name")
		}

		var project string
		if m.form.GetString("project") != "" {
			project = "Project: " + m.form.GetString("project")
		}

		var description string
		if m.form.GetString("description") != "" {
			description = "Description: " + m.form.GetString("description")
		}

		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{
			var (
				buildInfo = "(None)"
			)

			const statusWidth = 28
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.Copy().
				Height(lipgloss.Height(form)).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(s.StatusHeader.Render("Current Task") + "\n" +
					buildInfo +
					name + "\n" +
					project + "\n" +
					description + "\n",
				)
		}

		errors := m.form.Errors()
		header := m.appBoundaryView("Get Shit Done!")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}

}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

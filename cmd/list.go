/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"getshitdone/components"
	"getshitdone/db"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := components.TableModel{Table: components.TableColumnSturcture()}
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		return nil
	},
}

func setupTable(tasks []db.Task) *table.Table {
	columns := []string{"ID", "Name", "Description", "Project", "Status", "Created At", "Updated At"}
	var rows [][]string
	for _, task := range tasks {
		rows = append(rows, []string{
			fmt.Sprintf("%d", task.ID),
			task.Name,
			task.Description,
			task.Status,
			task.Project,
			task.CreatedAt.Format("2006-01-02"),
			task.UpdatedAt.Format("2006-01-02"),
		})
	}
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers(columns...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("64")).
					Border(lipgloss.NormalBorder()).
					BorderTop(false).
					BorderLeft(false).
					BorderRight(false).
					BorderBottom(true).
					Bold(true)
			}
			if row%2 == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("146"))
			}
			return lipgloss.NewStyle()
		})
	return t
}

func init() {
	rootCmd.AddCommand(listCmd)
}

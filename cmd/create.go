/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"getshitdone/components"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Add a new task to the list of tasks.",
	RunE: func(cmd *cobra.Command, args []string) error {

		createTask()

		return nil
	},
}

func createTask() {
	p := tea.NewProgram(components.NewModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	createCmd.Flags().StringP("project", "p", "", "Project the task belongs to")
	createCmd.Flags().StringP("description", "d", "", "Description of the task")
	rootCmd.AddCommand(createCmd)
}

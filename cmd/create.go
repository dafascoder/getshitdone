/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"getshitdone/db"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create NAME",
	Short: "Add a new task to the list of tasks.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := db.OpenDB()
		if err != nil {
			return err
		}
		defer t.Close()
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}
		return t.AddTask(args[0], description, project)
	},
}

func init() {
	createCmd.Flags().StringP("project", "p", "", "Project the task belongs to")
	createCmd.Flags().StringP("description", "d", "", "Description of the task")
	rootCmd.AddCommand(createCmd)
}

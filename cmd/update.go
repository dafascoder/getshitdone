/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"getshitdone/db"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("update called")
		t, err := db.OpenDB()
		if err != nil {
			return err
		}
		defer t.Close()
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		progress, err := cmd.Flags().GetInt("status")
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		var status string
		switch progress {
		case int(db.InProgress):
			status = db.InProgress.String()
		case int(db.Done):
			status = db.Done.String()
		default:
			status = db.Done.String()
		}
		newTask := db.Task{
			ID:          id,
			Name:        name,
			Description: description,
			Project:     project,
			Status:      status,
			UpdatedAt:   time.Now(),
		}

		return t.UpdateTask(newTask)
	},
}

func init() {
	updateCmd.Flags().StringP("name", "n", "", "Name of the task")
	updateCmd.Flags().StringP("description", "d", "", "Description of the task")
	updateCmd.Flags().StringP("project", "p", "", "Project the task belongs to")
	updateCmd.Flags().IntP(
		"status",
		"s",
		int(db.Todo),
		"specify a status for your task",
	)
	rootCmd.AddCommand(updateCmd)
}

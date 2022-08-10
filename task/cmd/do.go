package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/db"
	"os"
	"strconv"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the argument:", arg)
				os.Exit(1)
			} else {
				ids = append(ids, id)
			}
		}
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number:", id)
			} else {
				task := tasks[id-1]
				err := db.DeleteTasks(task.Key)
				if err != nil {
					fmt.Printf("Failed to mark \"%d\" as complete. Error \"%s\"", id, err)
				}
				fmt.Printf("You have completed the \"%s\" task.\n", task.Value)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}

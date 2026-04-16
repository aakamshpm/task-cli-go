package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aakamshpm/task-cli-go/internal/storage/jsonstore"
	"github.com/aakamshpm/task-cli-go/internal/task"
)

func main() {

	store := jsonstore.NewStore("data/tasks.json")
	service := task.NewService(store, nil)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	command := os.Args[1]

	switch command {
	case "help", "-h", "--help":
		printUsage()
	case "add":
		handleAdd(service)
	case "list":
		handleList(service)
	case "delete":
		handleDelete(service)
	case "done":
		handleDone(service)
	default:
		fmt.Printf("Unknown command: %q\n\n", command)
		printUsage()
		os.Exit(2)
	}
}

func handleAdd(service *task.Service) {
	if len(os.Args) < 3 {
		exitWithError(errors.New("usage: taskcli add <title>"))
	}

	// slice from index 2 onwards
	//
	title := strings.Join(os.Args[2:], " ")
	created, err := service.Add(title)

	if err != nil {
		exitWithError((err))
	}

	fmt.Printf("Added task %d: %s\n", created.ID, created.Title)
}

func handleList(service *task.Service) {
	tasks, err := service.List()

	if err != nil {
		exitWithError(err)
	}

	if len(tasks) == 0 {
		fmt.Println("no tasks found")
		return
	}

	for _, t := range tasks {
		status := " "
		if t.Done {
			status = "x" // for tasks that are done, we mark them as 'x' while printing
		}
		fmt.Printf("[%s] %d - %s\n", status, t.ID, t.Title)
	}
}

func handleDone(service *task.Service) {
	id, err := parseIDArg()
	if err != nil {
		exitWithError(err)
	}

	updated, err := service.MarkDone(id)
	if err != nil {
		exitWithError(err)
	}

	fmt.Printf("Task %d marked done: %s\n", updated.ID, updated.Title)

}

func handleDelete(service *task.Service) {
	id, err := parseIDArg()
	if err != nil {
		exitWithError(err)
	}

	if err := service.Delete(id); err != nil {
		exitWithError(err)
	}

	fmt.Printf("Deleted task %d\n", id)
}

func printUsage() {
	fmt.Println("Task CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  taskcli add <title>")
	fmt.Println("  taskcli list")
	fmt.Println("  taskcli done <id>")
	fmt.Println("  taskcli delete <id>")
}

func parseIDArg() (int, error) {
	if len(os.Args) != 3 {
		return 0, errors.New("usage: taskcli <done|delete> <id>")
	}

	// 'Atoi' stands for ASCII to integer which converts the respective id string to integer
	id, err := strconv.Atoi(os.Args[2])

	if err != nil {
		return 0, fmt.Errorf("invalid id %q", os.Args[2])
	}

	if id <= 0 {
		return 0, errors.New("id must be greater than 0")
	}

	return id, nil
}

func exitWithError(err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
}

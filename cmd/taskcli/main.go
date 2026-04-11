package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	command := os.Args[1]

	switch command {
	case "help", "-h", "--help":
		printUsage()
	case "add", "list", "delete":
		fmt.Printf("TODO: implement command %q\n", command)
	default:
		fmt.Printf("Unknown command: %q\n\n", command)
		printUsage()
		os.Exit(2)
	}
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

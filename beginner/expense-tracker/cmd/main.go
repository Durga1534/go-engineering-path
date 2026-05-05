package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	desc := addCmd.String("description", "", "Description of expense")
	amount := addCmd.Float64("amount", 0, "Amount of expense")

	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	month := summaryCmd.Int("month", 0, "Month of filter by")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', 'summary', or 'delete', subcommands")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *desc == "" || *amount == 0 {
			fmt.Println("Please provide a valid description and amount > 0")
			return
		}
		fmt.Printf("Expense add successfully: %s ($%.2f)\n", *desc, *amount)
	case "summary":
		summaryCmd.Parse(os.Args[2:])
		fmt.Printf("Total expenses of month %d: ...\n", *month)
	}
}

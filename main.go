package main

import (
	"fmt"
	"hozon/cli"
	"hozon/postgres"
	"log"
)

func main() {
	fmt.Println("Hozon started...")

	dbSettings, err := cli.GetArguments()
	if err != nil {
		log.Fatalf("Failed to get arguments: %v", err)
	}

	fmt.Println("Starting backup process...")
	err = postgres.Backup(dbSettings)

	if err != nil {
		log.Fatalf("Failed to backup database: %v", err)
	}

	fmt.Println("Backup process completed successfully.")
}

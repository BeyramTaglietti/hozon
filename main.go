package main

import (
	"fmt"
	"hozon/cli"
	"hozon/postgres"
	"log"
	"time"
)

func main() {
	fmt.Println("Hozon started...")

	cliSettings, err := cli.GetArguments()
	if err != nil {
		log.Fatalf("Failed to get arguments: %v", err)
	}

	log.Println("Running first backup...")
	err = postgres.Backup(cliSettings)
	if err != nil {
		log.Fatalf("Failed to backup database: %v", err)
	}
	log.Println("First backup completed successfully.")

	ticker := time.NewTicker(time.Duration(cliSettings.BackupFrequency) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Starting scheduled backup process...")
		err = postgres.Backup(cliSettings)

		if err != nil {
			log.Fatalf("Failed to backup database: %v", err)
		}

		log.Println("Scheduled backup process completed successfully.")
	}
}

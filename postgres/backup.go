package postgres

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type PGBackupSettings struct {
	DbName          string
	DbUser          string
	DbPass          string
	DbHost          string
	DbPort          int
	BackupFrequency int
}

func runBackup(dbSettings PGBackupSettings, backupDir string) error {

	backupFile := fmt.Sprintf("%s/%s.dump", backupDir, time.Now().Format("HozonBackup_2006-01-02__15_04_05"))

	cmd := exec.Command("pg_dump",
		"-h", dbSettings.DbHost,
		"-p", fmt.Sprintf("%d", dbSettings.DbPort),
		"-U", dbSettings.DbUser,
		"-F", "c", // Custom format
		"-b",             // Include large objects
		"-v",             // Verbose mode
		"-f", backupFile, // Output file
		dbSettings.DbName,
	)

	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", dbSettings.DbPass))

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute pg_dump: %v", err)
	}

	return nil
}

func InitBackupProcess(dbSettings PGBackupSettings) {
	log.Println("Running first backup...")

	backupDir := "./backups"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			log.Default().Fatalf("failed to create backup directory: %v", err)
		}
	}

	err := runBackup(dbSettings, backupDir)
	if err != nil {
		log.Fatalf("Failed to backup database: %v", err)
	}

	log.Println("First backup completed successfully.")

	ticker := time.NewTicker(time.Duration(dbSettings.BackupFrequency) * time.Millisecond * 400)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Starting scheduled backup process...")
		err = runBackup(dbSettings, backupDir)

		if err != nil {
			log.Fatalf("Failed to backup database: %v", err)
		}

		log.Println("Scheduled backup process completed successfully.")
	}
}

package postgres

import (
	"fmt"
	"hozon/cli"
	"os"
	"os/exec"
	"time"
)

func Backup(dbSettings cli.OsArguments) error {

	backupDir := "./backups"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create backup directory: %v", err)
		}
	}

	backupFile := fmt.Sprintf("%s/%s.dump", backupDir, time.Now().Format("2006-01-02__15_04_05"))

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

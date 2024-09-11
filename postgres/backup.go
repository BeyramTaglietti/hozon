package postgres

import (
	"fmt"
	"hozon/telegram"
	"log"
	"os"
	"os/exec"
	"time"
)

type BackupSettings struct {
	BackupFrequency int
	CleanDirectory  bool
}

type PostgresSettings struct {
	DbName string
	DbUser string
	DbPass string
	DbHost string
	DbPort int
}

func InitBackupProcess(postgresSettings PostgresSettings, telegramSettings telegram.TelegramSettings, backupSettings BackupSettings) {

	telegram.SendGreeting(telegramSettings.TGBotToken, telegramSettings.TGChatID)
	logBackup(telegramSettings.TGBotToken, telegramSettings.TGChatID, "Starting backup process...", false)

	backupDir := "./backups"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			logBackup(telegramSettings.TGBotToken, telegramSettings.TGChatID, fmt.Sprintf("Failed to create backup directory!\n%s", err), true)
			os.Exit(1)
		}
	}

	filePath, err := runBackup(postgresSettings, backupDir, backupSettings.CleanDirectory)
	if err != nil {
		logBackup(telegramSettings.TGBotToken, telegramSettings.TGChatID, fmt.Sprintf("Failed to backup database!\n%s", err), true)
		os.Exit(1)
	}

	logBackup(telegramSettings.TGBotToken, telegramSettings.TGChatID, "First backup completed successfully.", false)
	sendFile(telegramSettings.TGBotToken, telegramSettings.TGChatID, filePath)

	ticker := time.NewTicker(time.Duration(backupSettings.BackupFrequency) * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		logBackup(telegramSettings.TGBotToken, telegramSettings.TGChatID, "Starting scheduled backup process...", false)

		filePath, err := runBackup(postgresSettings, backupDir, backupSettings.CleanDirectory)

		if err != nil {
			logBackup(telegramSettings.TGBotToken, telegramSettings.TGChatID, fmt.Sprintf("Failed to run scheduled database backup! Process will not stop\n%s", err), true)
		} else {
			logBackup(telegramSettings.TGBotToken, telegramSettings.TGChatID, "Scheduled backup process completed successfully.", false)
			sendFile(telegramSettings.TGBotToken, telegramSettings.TGChatID, filePath)
		}
	}
}

func runBackup(dbSettings PostgresSettings, backupDir string, cleanDir bool) (string, error) {

	if cleanDir {
		// delete previous backups
		err := cleanDirectory(backupDir)
		if err != nil {
			return "", fmt.Errorf("failed to clean backup directory: %v", err)
		}
	}

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
		return "", fmt.Errorf("failed to execute pg_dump: %v", err)
	}

	return backupFile, nil
}

func cleanDirectory(backupDir string) error {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := os.Remove(fmt.Sprintf("%s/%s", backupDir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to delete file: %v", err)
		}
	}

	return nil
}

func logBackup(token string, chatid string, message string, error bool) {
	if error {
		message = fmt.Sprintf("Error: %s", message)
	}

	if error {
		telegram.SendMessage(
			token, telegram.CreateTelegramTextRequest(
				chatid, "There was an error in the backup process, please check your logs",
			),
		)
		log.Fatal(message)
	} else {
		telegram.SendMessage(token,
			telegram.CreateTelegramTextRequest(chatid, message),
		)
		log.Println(message)
	}

}

func sendFile(token string, chatid string, filePath string) {
	go telegram.SendFile(token, telegram.CreateTelegramDocumentRequest(chatid, filePath))
}

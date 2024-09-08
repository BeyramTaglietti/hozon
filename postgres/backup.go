package postgres

import (
	"fmt"
	"hozon/telegram"
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
	TGBotToken      string
	TGChatID        string
}

func runBackup(dbSettings PGBackupSettings, backupDir string) (string, error) {

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

func InitBackupProcess(backupSettings PGBackupSettings) {
	// telegram.SendGreeting(backupSettings.TGBotToken, backupSettings.TGChatID)
	logBackup(backupSettings.TGBotToken, backupSettings.TGChatID, "Starting backup process...", false)

	backupDir := "./backups"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			logBackup(backupSettings.TGBotToken, backupSettings.TGChatID, "Failed to create backup directory!", true)
			os.Exit(1)
		}
	}

	filePath, err := runBackup(backupSettings, backupDir)
	if err != nil {
		logBackup(backupSettings.TGBotToken, backupSettings.TGChatID, "Failed to backup database!", true)
		os.Exit(1)
	}

	logBackup(backupSettings.TGBotToken, backupSettings.TGChatID, "First backup completed successfully.", false)
	sendFile(backupSettings.TGBotToken, backupSettings.TGChatID, filePath)

	ticker := time.NewTicker(time.Duration(backupSettings.BackupFrequency) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		logBackup(backupSettings.TGBotToken, backupSettings.TGChatID, "Starting scheduled backup process...", false)

		filePath, err := runBackup(backupSettings, backupDir)

		if err != nil {
			logBackup(backupSettings.TGBotToken, backupSettings.TGChatID, "Failed to run scheduled database backup! Process will not stop", true)
		} else {
			logBackup(backupSettings.TGBotToken, backupSettings.TGChatID, "Scheduled backup process completed successfully.", false)
			sendFile(backupSettings.TGBotToken, backupSettings.TGChatID, filePath)
		}
	}
}

func logBackup(token string, chatid string, message string, error bool) {
	if error {
		message = fmt.Sprintf("Error: %s", message)
	}

	if error {
		log.Fatal(message)
	} else {
		log.Println(message)
	}

	go telegram.SendMessage(token,
		telegram.CreateTelegramTextRequest(chatid, message),
	)
}

func sendFile(token string, chatid string, filePath string) {
	go telegram.SendFile(token, telegram.CreateTelegramDocumentRequest(chatid, filePath))
}

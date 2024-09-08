package main

import (
	"hozon/postgres"
	"hozon/telegram"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppSettings struct {
	PostgresSettings postgres.PostgresSettings
	TelegramSettings telegram.TelegramSettings
	BackupSettings   postgres.BackupSettings
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := requireEnv("DB_USER")
	dbName := requireEnv("DB_NAME")
	dbPassword := requireEnv("DB_PASSWORD")
	dbPort := requireEnv("DB_PORT")
	dbPortNum, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatalf("Error converting DB_PORT to int: %v", err)
		os.Exit(1)
	}

	dbHost := requireEnv("DB_HOST")

	tgToken := requireEnv("TELEGRAM_TOKEN")
	tgChatID := requireEnv("TELEGRAM_CHATID")

	backupFrequency := requireEnv("BACKUP_FREQUENCY")

	backupFrequencyNum, err := strconv.Atoi(backupFrequency)
	if err != nil {
		log.Fatalf("Error converting BACKUP_FREQUENCY to int: %v", err)
		os.Exit(1)
	}

	cleanDirectory := requireEnv("CLEAN_DIRECTORY")
	cleanDirectoryBoolean, err := strconv.ParseBool(cleanDirectory)
	if err != nil {
		log.Fatalf("Error converting CLEAN_DIRECTORY to bool: %v", err)
		os.Exit(1)
	}

	settings := AppSettings{
		PostgresSettings: postgres.PostgresSettings{
			DbUser: dbUser,
			DbName: dbName,
			DbPass: dbPassword,
			DbHost: dbHost,
			DbPort: dbPortNum,
		},

		TelegramSettings: telegram.TelegramSettings{
			TGBotToken: tgToken,
			TGChatID:   tgChatID,
		},

		BackupSettings: postgres.BackupSettings{
			BackupFrequency: backupFrequencyNum,
			CleanDirectory:  cleanDirectoryBoolean,
		},
	}

	postgres.InitBackupProcess(settings.PostgresSettings, settings.TelegramSettings, settings.BackupSettings)
}

func requireEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is required", key)
		os.Exit(1)
	}

	return value
}

package cmd

import (
	"hozon/postgres"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "Hozon",
	Short: `Hozon allows you to backup your PostgreSQL database and store it using Telegram's BOT API.

Please use the suggested flags to interact with the application and provide your database settings.`,
	Run: func(cmd *cobra.Command, args []string) {

		name, err := cmd.Flags().GetString("dbname")
		if err != nil {
			log.Println("Error reading flag:", err)
			os.Exit(1)
		}

		user, err := cmd.Flags().GetString("dbuser")
		if err != nil {
			log.Println("Error reading flag:", err)
			os.Exit(1)
		}

		password, err := cmd.Flags().GetString("dbpass")
		if err != nil {
			log.Println("Error reading flag:", err)
			os.Exit(1)
		}

		host, err := cmd.Flags().GetString("dbhost")
		if err != nil {
			log.Println("Error reading flag:", err)
			os.Exit(1)
		}

		port, err := cmd.Flags().GetInt("dbport")
		if err != nil {
			log.Println("Error reading flag:", err)
			os.Exit(1)
		}

		frequency, err := cmd.Flags().GetInt("frequency")
		if err != nil {
			log.Println("Error reading flag:", err)
			os.Exit(1)
		}

		token, err := cmd.Flags().GetString("tgtoken")
		if err != nil {
			log.Println("Error reading flag:", err)
			os.Exit(1)
		}

		chatid, err := cmd.Flags().GetString("tgchatid")
		if err != nil {
			log.Println("Error reading flag:", err)
			os.Exit(1)
		}

		if name == "" || user == "" || password == "" || host == "" || frequency == 0 || token == "" || chatid == "" {
			log.Println("Please provide all the required flags")
			log.Println("Use --help for more information")
			os.Exit(1)
		}

		postgres.InitBackupProcess(postgres.PGBackupSettings{
			DbName:          name,
			DbUser:          user,
			DbPass:          password,
			DbHost:          host,
			DbPort:          port,
			BackupFrequency: frequency,
			TGBotToken:      token,
			TGChatID:        chatid,
		})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("dbname", "", "Your database name")
	rootCmd.Flags().String("dbuser", "", "Your database user")
	rootCmd.Flags().String("dbpass", "", "Your database password")
	rootCmd.Flags().String("dbhost", "", "Your database host")
	rootCmd.Flags().Int("dbport", 5432, "Your database port")
	rootCmd.Flags().Int("frequency", 1, "Backup frequency in Minutes")
	rootCmd.Flags().String("tgtoken", "", "Telegram Bot Token")
	rootCmd.Flags().String("tgchatid", "", "Telegram Chat ID")
}

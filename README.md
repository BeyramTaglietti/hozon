# Hozon (保存)

<img src="assets/hozon.jpg" width=200>

- [Hozon (保存)](#hozon-保存)
  - [Description](#description)
  - [How does it work](#how-does-it-work)
  - [Limitations](#limitations)
  - [Size](#size)
  - [Features](#features)
  - [Installation and Usage (From Source)](#installation-and-usage-from-source)
    - [Prerequisites](#prerequisites)
    - [Steps (From Source)](#steps-from-source)
  - [Installation and Usage (Using Docker/Compose)](#installation-and-usage-using-dockercompose)
    - [Prerequisites (Using Docker/Compose)](#prerequisites-using-dockercompose)
    - [Steps (Using Docker/Compose)](#steps-using-dockercompose)
  - [Contributing](#contributing)
  - [License](#license)

## Description

Hozon (保存) – hoh-zohn – is a command-line interface (CLI) tool designed to simplify the process of backing up PostgreSQL databases and securely sending the backup file via [Telegram](https://telegram.org/).

<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/8/82/Telegram_logo.svg/512px-Telegram_logo.svg.png" alt="Telegram logo" width=50>

With Hozon, you can ensure that your database backups are handled efficiently and shared privately with ease utilizing the Telegram messaging platform and its bot API.

## How does it work

You can Run Hozon on your local machine or in a Docker container, depending on your preference and it will run in the background, creating backups at the specified interval and sending them to your Telegram chat.

## Limitations

Telegram allows a maximum file size of [**50MB** for sending files](https://core.telegram.org/bots/api#senddocument), so please ensure that your database backups are within this limit.

I will eventually add the ability to interact with a [Telegram BOT hosted locally](https://core.telegram.org/bots/api#senddocument) which has no file size limit in the future.

## Size

- The binary size is around **8MB (8,000KB)**
- The Docker image size is around **22MB (22,000KB)**

## Features

- **Automated PostgreSQL Backup:** Easily create backups of your PostgreSQL databases using pg_dump.
- **Secure Transmission:** Send your backup file directly to a private Telegram by providing the chat id.
- **Easy Configuration:** Simple setup and configuration to get started quickly, it takes one (long) command to have a fully functioning backup system.
- **Cross-Platform:** Works on various operating systems, including Linux, macOS, and Windows.

## Installation and Usage (From Source)

### Prerequisites

- Golang 1.23 or higher
- pg_dump (part of PostgreSQL)
- Telegram bot token and chat ID

### Steps (From Source)

1. #### Clone the repository

   ```bash
   git clone https://github.com/BeyramTaglietti/hozon.git
   cd hozon
   ```

2. #### Build the binary

   ```bash
   go build -o hozon main.go
   ```

3. #### Configure your telehram bot token and chat ID

   - Create a new bot on Telegram via BotFather and obtain your bot token. [Learn more](https://core.telegram.org/bots#how-do-i-create-a-bot)
   - Get your chat ID (you can use any bot you'd like, I purposefully did not include this functionality in Hozon to keep it simple, I'd suggest using [@getmyid_bot](https://t.me/getmyid_bot), just add it to any group/channel and receive the chat ID, you can even use your userID as the chatID and have Hozon send messages privately to you).

4. #### Set the environment variables

   ```bash
   DB_USER="..."
   DB_NAME="..."
   DB_PASSWORD="..."
   DB_PORT="..."
   DB_HOST="..."

   TELEGRAM_TOKEN="..."
   TELEGRAM_CHATID="..."

   BACKUP_FREQUENCY=... # (int) in hours
   CLEAN_DIRECTORY=... # (bool) CLEAN_DIRECTORY=true will keep the backup directory clean, false will keep all backups
   ```

5. #### Run the binary

   ```bash
    ./hozon
   ```

## Installation and Usage (Using Docker/Compose)

### Prerequisites (Using Docker/Compose)

- Docker
- Telegram bot token and chat ID

### Steps (Using Docker/Compose)

1. #### Copy the Docker Compose file you can find in the repository

2. #### Provide the environment variables in a `.env` file

   ```bash
   DB_USER="..."
   DB_NAME="..."
   DB_PASSWORD="..."
   DB_PORT="..."
   DB_HOST="..."

   TELEGRAM_TOKEN="..."
   TELEGRAM_CHATID="..."

   BACKUP_FREQUENCY=... # (int) in hours
   CLEAN_DIRECTORY=... # (bool) CLEAN_DIRECTORY=true will keep the backup directory clean, false will keep all backups
   ```

3. #### Run docker compose up (make sure to provide the desidered POSTGRES VERSION in the Docker image ARGs in order for pg_dump to work correctly - Default is 15)

   ```bash
    docker compose up -d
   ```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the **MIT License**. See the [LICENSE](./LICENSE) file for details.

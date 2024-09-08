package cli

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

const (
	db_name          string = "--dbname"
	db_user          string = "--dbuser"
	db_pass          string = "--dbpass"
	db_host          string = "--dbhost"
	db_port          string = "--dbport"
	backup_frequency string = "--frequency"
)

type OsArguments struct {
	DbName          string
	DbUser          string
	DbPass          string
	DbHost          string
	DbPort          int
	BackupFrequency int
}

func GetArguments() (OsArguments, error) {

	name, ok := getSettingValue(os.Args, db_name)
	if !ok {
		return OsArguments{}, errors.New("missing database name argument")
	}

	user, ok := getSettingValue(os.Args, db_user)
	if !ok {
		return OsArguments{}, errors.New("missing database user argument")
	}

	password, ok := getSettingValue(os.Args, db_pass)
	if !ok {
		return OsArguments{}, errors.New("missing database password argument")
	}

	host, ok := getSettingValue(os.Args, db_host)
	if !ok {
		return OsArguments{}, errors.New("missing database host argument")
	}

	port, ok := getSettingValue(os.Args, db_port)
	if !ok {
		return OsArguments{}, errors.New("missing database port argument")
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return OsArguments{}, errors.New("invalid database port argument")
	}

	backupFrequency, ok := getSettingValue(os.Args, backup_frequency)
	if !ok {
		backupFrequency = "1"
	}

	frequency, err := strconv.Atoi(backupFrequency)
	if err != nil {
		return OsArguments{}, errors.New("invalid backup frequency argument")
	}

	arguments := OsArguments{
		DbName:          name,
		DbUser:          user,
		DbPass:          password,
		DbHost:          host,
		DbPort:          portNumber,
		BackupFrequency: frequency,
	}

	return arguments, nil
}

func getSettingValue(cmdArgs []string, setting string) (string, bool) {
	for _, arg := range cmdArgs {
		if strings.HasPrefix(arg, setting) {

			if !strings.Contains(arg, "=") {
				return "", true
			}

			return strings.Split(arg, "=")[1], true
		}
	}

	return "", false
}

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/satori/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// readEnvAndSetDefault устанавливает переменные конфигурации viper по умолчанию. Используется для случая,
// когда файл конфигурации не был найден. Использует переменные окружения для настройки.
func readEnvAndSetDefault(logger *zap.Logger) {
	// POSTGRES
	if port := os.Getenv("P_POSTGRES_PORT"); port != "" {
		psqlPort, err := strconv.Atoi(port)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'P_POSTGRES_PORT', so it will be with default value 5432")
			viper.SetDefault("postgres.port", 5432)
		} else {
			viper.SetDefault("postgres.port", psqlPort)
		}
	} else {
		viper.SetDefault("postgres.port", 5432)
	}

	if host := os.Getenv("P_POSTGRES_CONNECTION_HOST"); host != "" {
		viper.SetDefault("postgres.connectionHost", host)
	} else {
		viper.SetDefault("postgres.connectionHost", "postgres_privelege")
	}

	viper.SetDefault("postgres.sslmode", "disable")

	// SERVER
	if address := os.Getenv("PS_SERVER_ADDRESS"); address != "" {
		viper.SetDefault("server.address", address)
	} else {
		viper.SetDefault("server.address", ":8010")
	}

	if writeTimeout := os.Getenv("PS_SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		timeout, err := time.ParseDuration(writeTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'PS_SERVER_WRITE_TIMEOUT', so it will be with default value 5s")
			viper.SetDefault("server.write_timeout", 5*time.Second)
		} else {
			viper.SetDefault("server.write_timeout", timeout)
		}
	} else {
		viper.SetDefault("server.write_timeout", 5*time.Second)
	}

	if readTimeout := os.Getenv("PS_SERVER_READ_TIMEOUT"); readTimeout != "" {
		timeout, err := time.ParseDuration(readTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'PS_SERVER_READ_TIMEOUT', so it will be with default value 5s")
			viper.SetDefault("server.read_timeout", 5*time.Second)
		} else {
			viper.SetDefault("server.read_timeout", timeout)
		}
	} else {
		viper.SetDefault("server.read_timeout", 5*time.Second)
	}

	if idleTimeout := os.Getenv("SERVER_IDLE_TIMEOUT"); idleTimeout != "" {
		timeout, err := time.ParseDuration(idleTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'SERVER_IDLE_TIMEOUT', so it will be with default value 3s")
			viper.SetDefault("server.idle_timeout", 3*time.Second)
		} else {
			viper.SetDefault("server.idle_timeout", timeout)
		}
	} else {
		viper.SetDefault("server.idle_timeout", 3*time.Second)
	}

	if shutdownDuration := os.Getenv("SERVER_SHUTDOWN_DURATION"); shutdownDuration != "" {
		duration, err := time.ParseDuration(shutdownDuration)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'PS_SERVER_SHUTDOWN_DURATION', so it will be with default value 10s")
			viper.SetDefault("server.shutdown_duration", 10*time.Second)
		} else {
			viper.SetDefault("server.shutdown_duration", duration)
		}
	} else {
		viper.SetDefault("server.shutdown_duration", 10*time.Second)
	}

	viper.SetDefault("secret_key", uuid.NewV4().String())
}

// Read получает переменные из среды и файла конфигурации
func Read(configFilePath string, logger *zap.Logger) {
	readEnvAndSetDefault(logger)
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			logger.Fatal(fmt.Sprintf("fatal error config file: %v", err))
		}
		logger.Warn(fmt.Sprintf("configuration file is not found, programm will be executed within default configuration: %v", err))
		return
	}
	logger.Info("successful read of configuration")
}

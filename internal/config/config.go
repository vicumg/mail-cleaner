package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	IMAPServer string
	IMAPPort   int
	Email      string
	Password   string
}

func LoadConfig(service_name string) *Config {
	// load from env
	godotenv.Load(".env")
	server := os.Getenv(service_name + "_IMAP_SERVER")
	portStr := os.Getenv(service_name + "_IMAP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Invalid IMAP_PORT")
		panic("Invalid IMAP_PORT")
	}

	email := os.Getenv(service_name + "_EMAIL")
	password := os.Getenv(service_name + "_PASSWORD")

	return &Config{
		IMAPServer: server,
		IMAPPort:   port,
		Email:      email,
		Password:   password,
	}
}

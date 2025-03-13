package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func (c Config) SetDBUrl(dbUrl string) error {
	c.DBUrl = dbUrl
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}

func write(c Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	pathConfig := homeDir + "/" + configFileName
	file, err := os.OpenFile(pathConfig, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// Write struct to file as JSON
	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	// Read the config file
	config := &Config{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return *config, err
	}
	pathConfig := homeDir + "/" + configFileName
	file, err := os.Open(pathConfig)
	if err != nil {
		return *config, err
	}
	defer file.Close()

	// Parse the config file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return *config, err
	}
	return *config, nil
}

func ReadCfgFile() error {
	// Print the file content
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filePath := homeDir + "/" + configFileName
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	fmt.Println(string(content))
	return nil
}

func Hello() {
	fmt.Println("Hello, World!")
}

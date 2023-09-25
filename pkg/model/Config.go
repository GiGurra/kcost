package model

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Prices     PriceTable `yaml:"prices"`
	Namespaces Namespaces `yaml:"namespaces"`
}

type Namespaces struct {
	Excluded []string `yaml:"excluded"`
}

type PriceTable struct {
	GKE GKE `yaml:"gke"`
}

type GKE struct {
	Autopilot GkePrice `yaml:"autopilot"`
}

type GkePrice struct {
	Spot    Resource `yaml:"spot"`
	Regular Resource `yaml:"regular"`
}

type Resource struct {
	CPU float64 `yaml:"cpu"`
	RAM float64 `yaml:"ram"`
}

func ReadFile(path string) ([]byte, error) {

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Close file on exit
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("Error closing file: %v\n", err))
		}
	}(file)

	// Read file into bytes
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func NewConfigFromFile(path string) (Config, error) {

	// Read bytes from file
	bytes, err := ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file: %v", err)
	}

	// parse bytes into price table
	table := Config{}
	err = yaml.Unmarshal(bytes, &table)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing yaml: %v", err)
	}

	return table, nil
}
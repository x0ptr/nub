package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Sources         []string `json:"sources"`
	LLMAPIKey       string   `json:"llm_api_key"`
	LLMAPIURL       string   `json:"llm_api_url"`
	LLMAPIModel     string   `json:"llm_api_model"`
	ScheduleMinutes int      `json:"schedule_minutes"`
	SummaryPrompt   string   `json:"summary_prompt"`
	FocusTopics     string   `json:"focus_topics"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "nub", "config.json"), nil
}

func GetDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "nub"), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := &Config{
			Sources:         []string{},
			ScheduleMinutes: 15,
			SummaryPrompt:   "Summarize the key news topics and main stories from this website. Focus on the most important headlines and provide a concise overview in markdown format.",
		}
		if err := SaveConfig(config); err != nil {
			return nil, err
		}
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.SummaryPrompt == "" {
		config.SummaryPrompt = "Summarize the key news topics and main stories from this website. Focus on the most important headlines and provide a concise overview in markdown format."
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return err
	}

	fmt.Printf("Config saved to: %s\n", configPath)
	return nil
}

func AddSource(config *Config, url string) error {
	for _, source := range config.Sources {
		if source == url {
			return fmt.Errorf("source already exists: %s", url)
		}
	}
	config.Sources = append(config.Sources, url)
	return SaveConfig(config)
}

func RemoveSource(config *Config, idOrURL string) error {
	idx := -1
	
	for i := 0; i < 10 && i < len(idOrURL); i++ {
		if idOrURL[i] < '0' || idOrURL[i] > '9' {
			idx = -1
			break
		}
	}
	
	if idx == -1 {
		num := 0
		for i := 0; i < len(idOrURL); i++ {
			if idOrURL[i] >= '0' && idOrURL[i] <= '9' {
				num = num*10 + int(idOrURL[i]-'0')
			} else {
				num = -1
				break
			}
		}
		if num >= 1 && num <= len(config.Sources) {
			idx = num - 1
		}
	}
	
	if idx == -1 {
		for i, source := range config.Sources {
			if source == idOrURL {
				idx = i
				break
			}
		}
	}
	
	if idx == -1 {
		return fmt.Errorf("source not found: %s", idOrURL)
	}
	
	config.Sources = append(config.Sources[:idx], config.Sources[idx+1:]...)
	return SaveConfig(config)
}

func ListSources(config *Config) {
	if len(config.Sources) == 0 {
		fmt.Println("No sources configured")
		return
	}
	
	fmt.Println("Sources:")
	for i, source := range config.Sources {
		fmt.Printf("  [%d] %s\n", i+1, source)
	}
}

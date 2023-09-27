package config

import (
	"encoding/json"
	"ipmanlk/ani2mal/models"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type AppConfig struct {
	configDir         string
	malConfigPath     string
	anilistConfigPath string
	excludesFilePath  string
}

var (
	once     sync.Once
	instance *AppConfig
)

func GetAppConfig() *AppConfig {
	once.Do(
		func() {
			configDir, err := getConfigDir()

			if err != nil {
				log.Fatal("Failed to locate the configuration directory.", err)
			}

			instance = &AppConfig{
				configDir:         configDir,
				malConfigPath:     filepath.Join(configDir, "mal.json"),
				anilistConfigPath: filepath.Join(configDir, "anilist.json"),
				excludesFilePath:  filepath.Join(configDir, "excludes.json"),
			}
		})

	return instance
}

func (cfg *AppConfig) SaveMalConfig(malConfig *models.MalConfig) {
	jsonData, err := json.MarshalIndent(malConfig, "", " ")

	if err != nil {
		log.Fatal("Failed to marshal mal config", err)
	}

	err = os.WriteFile(cfg.malConfigPath, jsonData, 0644)

	if err != nil {
		log.Fatal("Error writing MAL config", err)
	}
}

func (cfg *AppConfig) GetMalConfig() *models.MalConfig {
	_, err := os.Stat(cfg.malConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Please login to MyAnimeList first")
		}
		log.Fatalf("Failed to read MyAnimeList configuration file. Check if file permissions are correct %+v", err)
	}

	content, _ := os.ReadFile(cfg.malConfigPath)

	var malConfig models.MalConfig
	json.Unmarshal(content, &malConfig)

	return &malConfig
}

func (cfg *AppConfig) SaveAnilistConfig(anilistConfig *models.AnilistConfig) {
	jsonData, err := json.MarshalIndent(anilistConfig, "", " ")

	if err != nil {
		log.Fatal("Failed to marshal Anilist config", err)
	}

	err = os.WriteFile(cfg.anilistConfigPath, jsonData, 0644)

	if err != nil {
		log.Fatal("Error writing Anilist config", err)
	}
}

func (cfg *AppConfig) GetAnilistConfig() *models.AnilistConfig {
	_, err := os.Stat(cfg.anilistConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Please configure Anilist first")
		}
		log.Fatalf("Failed to read Anilist configuration file. Check if file permissions are correct %+v", err)
	}

	content, _ := os.ReadFile(cfg.anilistConfigPath)

	var anilistConfig models.AnilistConfig
	json.Unmarshal(content, &anilistConfig)

	return &anilistConfig
}

func getConfigDir() (string, error) {
	var configDir string
	switch currentOs := runtime.GOOS; currentOs {
	case "windows":
		// On Windows, use %APPDATA%
		appdata := os.Getenv("APPDATA")
		configDir = filepath.Join(appdata, "ani2mal")
	case "linux", "darwin":
		// On Unix/Linux or macOS, use ~/.config
		home := os.Getenv("HOME")
		configDir = filepath.Join(home, ".config", "ani2mal")
	default:
		exePath, err := os.Executable()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(exePath, "ani2mal")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

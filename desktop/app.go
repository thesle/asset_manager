package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

// DesktopConfig holds the desktop app configuration
type DesktopConfig struct {
	APIUrl string `yaml:"api_url"`
	Token  string `yaml:"token"`
}

// App struct
type App struct {
	ctx    context.Context
	config *DesktopConfig
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadConfig()
}

// GetConfigPath returns the config file path
func (a *App) GetConfigPath() string {
	configPath, _ := xdg.ConfigFile(filepath.Join(AppName, "config.yaml"))
	return configPath
}

// loadConfig loads the configuration from file
func (a *App) loadConfig() {
	configPath := a.GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		a.config = &DesktopConfig{}
		return
	}

	var cfg DesktopConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		a.config = &DesktopConfig{}
		return
	}

	a.config = &cfg
}

// SaveConfig saves the configuration to file
func (a *App) SaveConfig(apiUrl string, token string) error {
	a.config = &DesktopConfig{
		APIUrl: apiUrl,
		Token:  token,
	}

	configPath := a.GetConfigPath()

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(a.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (a *App) GetConfig() *DesktopConfig {
	if a.config == nil {
		return &DesktopConfig{}
	}
	return a.config
}

// HasConfig returns true if a valid config exists
func (a *App) HasConfig() bool {
	return a.config != nil && a.config.APIUrl != ""
}

// ClearConfig clears the stored configuration
func (a *App) ClearConfig() error {
	a.config = &DesktopConfig{}
	configPath := a.GetConfigPath()
	return os.Remove(configPath)
}

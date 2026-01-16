package backend

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

// Startup - Wails v2 lifecycle hook
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// GetContext - Expose context for runtime calls
func (a *App) GetContext() context.Context {
	return a.ctx
}

// CheckDeviceByIP checks if a device is reachable on port 80
func (a *App) CheckDeviceByIP(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":80", 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// LoadConfig loads the configuration file from the default path
func (a *App) LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist, return empty config
		return &Config{Devices: []Device{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// OpenConfig opens a file dialog to select a JSON configuration file
func (a *App) OpenConfig() (string, error) {
	// Open file dialog and get the selected file path
	selectedFile, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a JSON file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JSON",
				Pattern:     "*.json",
			},
		},
	})

	if err != nil {
		return "", err
	}

	// If user cancelled, return empty string
	if selectedFile == "" {
		return "", nil
	}

	// Read the selected file
	data, err := os.ReadFile(selectedFile)
	if err != nil {
		return "", err
	}

	// Save to default config location
	path, err := configPath()
	if err != nil {
		return "", err
	}
	ncDir := filepath.Dir(path)
	os.MkdirAll(ncDir, 0755)
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}

	return selectedFile, nil
}

// configPath returns the default config file path
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Documents", "network-checker", "config.json"), nil
}

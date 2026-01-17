package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	logFile *os.File
}

// Startup - Wails v2 lifecycle hook
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.initLogger()
}

func (a *App) initLogger() {
	path, err := logPath()
	if err != nil {
		return
	}

	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	a.logFile = f
	log.SetOutput(f)
	log.Println("--- Application Started ---")
}

func (a *App) log(message string) {
	if a.logFile != nil {
		log.Println(message)
	}
	fmt.Println(message)
}

// GetContext - Expose context for runtime calls
func (a *App) GetContext() context.Context {
	return a.ctx
}

// CheckDeviceByIP checks if a device responds to ICMP Ping
func (a *App) CheckDeviceByIP(ip string) bool {
	a.log(fmt.Sprintf("Pinging device: %s", ip))

	pinger, err := probing.NewPinger(ip)
	if err != nil {
		a.log(fmt.Sprintf("  ERROR creating pinger for %s: %v", ip, err))
		return false
	}

	// Set to non-privileged (UDP) ping. Works on macOS and most modern OSs without root.
	pinger.SetPrivileged(false)
	pinger.Count = 1
	pinger.Timeout = 2 * time.Second

	err = pinger.Run() // Blocks until finished
	if err != nil {
		a.log(fmt.Sprintf("  FAILED pinging %s: %v", ip, err))
		return false
	}

	stats := pinger.Statistics()
	if stats.PacketLoss > 0 {
		a.log(fmt.Sprintf("  OFFLINE %s: Packet loss 100%%", ip))
		return false
	}

	a.log(fmt.Sprintf("  ONLINE %s: RTT %v", ip, stats.AvgRtt))
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

func logPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Documents", "network-checker", "app.log"), nil
}

# Network Checker 🌐

A professional, high-performance network monitoring application built with **Wails** (Go backend) and **Vanilla JavaScript** (Frontend).

## Features
- **Real-time Monitoring**: Check the availability of multiple devices across your network simultaneously.
- **Auto-Polling**: configurable check intervals (default 5s) to stay updated on device status.
- **Modern UI**: Sleek, glassmorphism-inspired dark interface with real-time status badges.
- **Configurable**: Easily load your infrastructure mapping via JSON.
- **Cross-Platform**: Built for Mac, Windows, and Linux.

## Project Structure
- `backend/`: Go source files for logic, networking, and configuration management.
- `frontend/`: 
    - `src/`: Core logic (`main.js`) and styling (`style.css`).
    - `wailsjs/`: Auto-generated Go-to-JavaScript bindings.
- `examples/`: Sample configuration files for testing.
- `design/`: Design assets including the official logo.

## Getting Started

### 1. Requirements
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- Go 1.18+
- Node.js & npm

### 2. Development
Run the application in development mode with live reload:
```bash
wails dev
```

### 3. Loading Configuration
The application looks for a configuration file in `~/Documents/network-checker/config.json`. You can load a new configuration using **File → Open Configuration** or by pressing **Cmd+O**.

Sample JSON structure:
```json
{
  "devices": [
    {
      "name": "Primary Server",
      "nics": [
        { "name": "eth0", "ip": "192.168.1.10" }
      ]
    }
  ]
}
```

## Building
To build a production-ready package for your OS:
```bash
wails build
```

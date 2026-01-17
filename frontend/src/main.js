// Import Wails v2 bindings
import { CheckDeviceByIP, LoadConfig, OpenConfig } from '../wailsjs/go/backend/App.js'
import { EventsOn } from '../wailsjs/runtime/runtime.js'

// Go Wails v2 calls
let devicesData = []
const lastChecks = {}
let checkIntervals = []

// Time formatting
function formatTime(date) {
  return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false })
}

// Render UI
function render() {
  console.log('Rendering UI with devices:', devicesData)
  const grid = document.getElementById('device-grid')
  if (!grid) {
    console.error('Grid element #device-grid not found')
    return
  }

  if (!devicesData || devicesData.length === 0) {
    showEmptyState('No devices configured')
    return
  }

  grid.innerHTML = devicesData.map(device => `
    <div class="device-card">
      <div class="device-card-header">
        <div class="device-info">
          <div class="device-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-server">
              <rect width="20" height="8" x="2" y="2" rx="2" ry="2"></rect>
              <rect width="20" height="8" x="2" y="14" rx="2" ry="2"></rect>
              <line x1="6" x2="6.01" y1="6" y2="6"></line>
              <line x1="6" x2="6.01" y1="18" y2="18"></line>
            </svg>
          </div>
          <div>
            <h3 class="device-name">${device.name}</h3>
            <p class="device-interfaces">${device.nics ? device.nics.length : 0} interface${device.nics && device.nics.length !== 1 ? 's' : ''}</p>
          </div>
        </div>
      </div>
      <div class="device-card-body">
        ${(device.nics || []).map(nic => {
    const key = `${device.name}-${nic.name}`
    if (!lastChecks[key]) lastChecks[key] = new Date()
    return `
            <div class="nic-card">
              <div class="nic-info">
                <div class="nic-header">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-network nic-icon">
                    <rect x="16" y="16" width="6" height="6" rx="1"></rect>
                    <rect x="2" y="16" width="6" height="6" rx="1"></rect>
                    <rect x="9" y="2" width="6" height="6" rx="1"></rect>
                    <path d="M5 16v-3a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1v3"></path>
                    <path d="M12 12V8"></path>
                  </svg>
                  <span class="nic-name">${nic.name}</span>
                  <span class="badge" id="status-${key}">?</span>
                </div>
                <p class="nic-ip">${nic.ip}</p>
              </div>
              <div class="nic-status">
                <p>Last check: <span id="check-${key}">${formatTime(lastChecks[key])}</span></p>
              </div>
            </div>
          `
  }).join('')}
      </div>
    </div>
  `).join('')
}

// Check status for a single NIC
async function checkStatus(device, nic) {
  const key = `${device.name}-${nic.name}`
  try {
    const online = await CheckDeviceByIP(nic.ip)
    const statusEl = document.getElementById(`status-${key}`)
    const checkEl = document.getElementById(`check-${key}`)

    if (statusEl) {
      statusEl.textContent = online ? 'Online' : 'Offline'
      statusEl.className = online ? 'badge' : 'badge offline'
    }

    lastChecks[key] = new Date()
    if (checkEl) {
      checkEl.textContent = formatTime(lastChecks[key])
    }
  } catch (err) {
    console.error(`Error checking status for ${key}:`, err)
  }
}

// Device polling
function setupIntervals() {
  console.log('Setting up intervals for devices:', devicesData)
  // Clear existing intervals
  checkIntervals.forEach(clearInterval)
  checkIntervals = []

  if (!devicesData) return

  devicesData.forEach(device => {
    if (!device.nics) return
    device.nics.forEach(nic => {
      // Immediate initial check
      checkStatus(device, nic)

      // Set interval for subsequent checks (every 5 seconds)
      const interval = setInterval(() => checkStatus(device, nic), 5000)
      checkIntervals.push(interval)
    })
  })
}

// Load config
async function loadConfigData() {
  try {
    console.log('Loading config data...')
    const cfg = await LoadConfig()
    console.log('Received result from LoadConfig:', cfg)

    if (!cfg) {
      console.error('LoadConfig returned null or undefined')
      showEmptyState('Error: Backend failed to provide configuration.')
      return
    }

    if (!cfg.devices || !Array.isArray(cfg.devices)) {
      console.log('No devices found or invalid structure, showing empty state')
      devicesData = []
      showEmptyState('No configuration found or file is empty')
      return
    }

    devicesData = cfg.devices
    console.log('Applied devicesData:', devicesData)

    // Ensure all devices have a nics array
    devicesData.forEach(device => {
      if (!device.nics || !Array.isArray(device.nics)) {
        device.nics = []
      }
    })

    if (devicesData.length === 0) {
      showEmptyState('Configuration file is empty')
      return
    }

    render()
    setupIntervals()
  } catch (err) {
    console.error('Error loading config:', err)
    showEmptyState(`Error loading configuration: ${err.message || err}`)
  }
}

// Show empty state message
function showEmptyState(message) {
  const grid = document.getElementById('device-grid')
  if (!grid) return

  grid.innerHTML = `
    <div style="grid-column: 1 / -1; text-align: center; padding: 3rem;">
      <div style="margin-bottom: 1rem; opacity: 0.3; display: flex; justify-content: center;">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"></path>
          <polyline points="14 2 14 8 20 8"></polyline>
          <line x1="12" x2="12" y1="18" y2="12"></line>
          <line x1="9" x2="15" y1="15" y2="15"></line>
        </svg>
      </div>
      <h2 style="font-size: 1.5rem; margin-bottom: 0.5rem; color: #a3a3a3;">${message}</h2>
      <p style="color: #737373; margin-bottom: 1rem;">Use <strong>File → Open Configuration</strong> or press <strong>Cmd+O</strong> to load a JSON file</p>
    </div>
  `
}

// Open config dialog
async function openConfigDialog() {
  try {
    const result = await OpenConfig()
    if (result) {
      console.log('Config file selected:', result)
      await loadConfigData()
    }
  } catch (err) {
    console.error('Error opening config:', err)
  }
}

// Events - keyboard shortcut only, button no longer exists
window.addEventListener('keydown', (e) => {
  if ((e.metaKey || e.ctrlKey) && e.key === 'o') {
    e.preventDefault()
    openConfigDialog()
  }
})

// Listen for config-loaded event from native menu
EventsOn('config-loaded', () => {
  console.log('Config loaded via menu, reloading...')
  loadConfigData()
})

// Initialize
loadConfigData()

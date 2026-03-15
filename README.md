# Shodan GUI

A full-featured desktop GUI for the [Shodan API](https://developer.shodan.io/api), built with Go + Wails v2 + Vue 3.

![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- **Host Lookup** — Query any IP for open ports, services, banners, SSL certs, location
- **Search** — Full Shodan query syntax with facets, pagination, table/card/raw views
- **Honeypot Filter** — Auto-filter results tagged as honeypots (toggle on/off)
- **Export** — Select fields and export current page or all results (up to 10,000) as JSON or CSV
- **DNS Tools** — Domain info, DNS resolve, reverse DNS, My IP
- **Scans** — Submit on-demand scans, track status, list active scans
- **Alerts** — Create/delete network alerts, enable/disable triggers
- **Notifiers** — View and manage notification integrations
- **Settings** — API key management with live plan/credit status

## Download

Pre-built binaries are available on the [Releases](../../releases) page.

| Platform | File |
|----------|------|
| Windows x64 | `shodan-gui-windows-amd64.exe` |
| macOS Apple Silicon | `shodan-gui-darwin-arm64` |
| macOS Intel | `shodan-gui-darwin-amd64` |
| Linux x64 | `shodan-gui-linux-amd64` |

## Build from Source

**Requirements:**
- [Go 1.21+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Wails v2](https://wails.io/docs/gettingstarted/installation): `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Linux only:** `libgtk-3-dev libwebkit2gtk-4.0-dev`

```bash
git clone https://github.com/justinwoo280/shodan-gui.git
cd shodan-gui
wails build
```

The binary will be at `build/bin/shodan-gui` (or `shodan-gui.exe` on Windows).

## Usage

1. Launch the app
2. Go to **Settings** and enter your [Shodan API key](https://account.shodan.io/)
3. Start searching

### Search Tips

```
# Find Apache servers in Germany
apache country:DE

# Find Nginx on port 80
product:nginx port:80

# Find open MongoDB instances
product:MongoDB -authentication

# Find specific CVE
vuln:CVE-2021-44228
```

### Export All Results

In the Search view, click **Export** → set **Scope** to **All results** → configure max pages (each page = 100 records = 1 query credit) → click **Fetch & Export**.

## API Credits

| Operation | Cost |
|-----------|------|
| Search (per page) | 1 query credit |
| Host Lookup | 1 query credit |
| Count | Free |
| Scan (per IP) | 1 scan credit |
| DNS queries | Free |

## Tech Stack

- **Backend:** Go + [Wails v2](https://wails.io/)
- **Frontend:** Vue 3 + TypeScript + Pinia + Vue Router
- **Theme:** GitHub Dark

## License

 [MIT](./LICENSE)

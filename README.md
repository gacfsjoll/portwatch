# port open ports and alerts on unexpected changes.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with default settings (scans every 60 seconds):

```bash
portwatch start
```

Specify a custom scan interval and alert on any change:

```bash
portwatch start --interval 30s --notify
```

Define a baseline of expected open ports:

```bash
portwatch baseline --ports 22,80,443
```

View current port status:

```bash
portwatch status
```

When an unexpected port opens or closes, `portwatch` logs an alert and optionally sends a desktop or webhook notification.

---

## Configuration

A config file can be placed at `~/.config/portwatch/config.yaml`:

```yaml
interval: 60s
notify: true
baseline:
  - 22
  - 80
  - 443
webhook: "https://hooks.example.com/alert"
```

---

## License

MIT © 2024 Your Name
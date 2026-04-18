# portwatch

Lightweight CLI to monitor and alert on open port changes on a host.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start monitoring open ports on the local host:

```bash
portwatch watch
```

Specify a target host and polling interval:

```bash
portwatch watch --host 192.168.1.10 --interval 30s
```

Alert output when a change is detected:

```
[ALERT] New port opened: 8080/tcp
[ALERT] Port closed: 3306/tcp
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `localhost` | Target host to monitor |
| `--interval` | `60s` | Polling interval |
| `--alert` | stdout | Alert destination (`stdout`, `webhook`) |
| `--webhook` | | Webhook URL for notifications |

### Example with Webhook

```bash
portwatch watch --host 10.0.0.5 --interval 15s --alert webhook --webhook https://hooks.example.com/notify
```

## Requirements

- Go 1.21+
- Linux / macOS / Windows

## License

MIT © 2024 yourusername
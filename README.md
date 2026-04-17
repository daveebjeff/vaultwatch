# vaultwatch

A CLI tool to monitor and alert on HashiCorp Vault secret expiration and lease renewals.

## Installation

```bash
go install github.com/yourusername/vaultwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultwatch.git && cd vaultwatch && go build -o vaultwatch .
```

## Usage

Set your Vault address and token, then run vaultwatch against a path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-token-here"

# Monitor secrets at a path and alert if expiring within 7 days
vaultwatch monitor --path secret/myapp --threshold 7d

# Watch lease renewals and get notified on expiration
vaultwatch leases --renew --alert-on-expire

# Output status of all monitored secrets
vaultwatch status --output json
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to monitor | `secret/` |
| `--threshold` | Alert threshold before expiration | `24h` |
| `--renew` | Automatically renew leases | `false` |
| `--output` | Output format (`text`, `json`) | `text` |
| `--interval` | Polling interval | `5m` |

## Configuration

vaultwatch respects standard Vault environment variables (`VAULT_ADDR`, `VAULT_TOKEN`, `VAULT_NAMESPACE`) and can also be configured via a `vaultwatch.yaml` file:

```yaml
vault_addr: https://vault.example.com
threshold: 48h
interval: 10m
paths:
  - secret/production
  - secret/staging
```

## Requirements

- Go 1.21+
- HashiCorp Vault 1.12+

## License

MIT © 2024 — see [LICENSE](LICENSE) for details.
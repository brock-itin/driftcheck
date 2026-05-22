# driftcheck

Detect config drift between running containers and their source Compose/Helm definitions.

---

## Installation

```bash
go install github.com/yourusername/driftcheck@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/driftcheck/releases).

---

## Usage

Point `driftcheck` at your Compose file or Helm chart and let it compare against your running containers:

```bash
# Check drift against a Docker Compose file
driftcheck --compose docker-compose.yml

# Check drift against a Helm release
driftcheck --helm my-release --namespace production

# Output results as JSON
driftcheck --compose docker-compose.yml --output json
```

Example output:

```
[DRIFT] web: environment variable PORT expected "8080", got "9090"
[DRIFT] api: image expected "myapp:v1.2.0", got "myapp:v1.1.3"
[OK]    db: no drift detected
```

---

## How It Works

`driftcheck` queries the Docker daemon (or Kubernetes API) to inspect running containers and compares their actual configuration — image tags, environment variables, port bindings, volume mounts, and resource limits — against the values defined in your source Compose or Helm files.

---

## Configuration

| Flag | Description | Default |
|------|-------------|---------|
| `--compose` | Path to a Docker Compose file | — |
| `--helm` | Helm release name to inspect | — |
| `--namespace` | Kubernetes namespace | `default` |
| `--output` | Output format: `text`, `json` | `text` |
| `--strict` | Exit with non-zero code on any drift | `false` |

---

## License

MIT © 2024 yourusername
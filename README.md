# Friendle

Self-hosted daily Wordle for friend groups. One person picks the word each day; everyone else guesses. Independent project — not affiliated with the New York Times Wordle.

## Stack

- **Backend:** Go, Fiber, SQLite
- **Frontend:** SvelteKit, Tailwind (embedded in the Go binary)
- **Deploy:** Single binary or Docker

## Quick start

### Build

```bash
make build
```

### Run

```bash
./friendle serve --enable-signup
```

On first run, open the bootstrap URL printed in the console to create the site admin account.

### Docker

```bash
docker build -t friendle -f docker/Dockerfile .
docker run -p 8080:80 -v friendle_data:/friendle_data friendle
```

## CLI flags

| Flag | Default | Description |
|------|---------|-------------|
| `--http` | `127.0.0.1:9081` | Listen address |
| `--data-dir` | `./friendle_data` | SQLite and data files |
| `--enable-signup` | off | Allow new user registration |
| `--dev` | off | Proxy UI to Vite dev server |
| `--debug` | off | Debug logging |

## Development

Terminal 1 (UI):

```bash
cd ui && pnpm install && pnpm run dev
```

Terminal 2 (API):

```bash
go mod download
air   # or: go run . serve --dev --enable-signup
```

Tests (no UI build required):

```bash
go test -tags dev ./...
```

## License

MIT

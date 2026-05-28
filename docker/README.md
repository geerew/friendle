# Docker

Build from the repository root:

```bash
docker build -t friendle -f docker/Dockerfile .
```

Run with persisted data:

```bash
docker run -p 8080:80 -v friendle_data:/friendle_data friendle
```

Or use the example compose file:

```bash
docker compose -f docker/docker-compose.example.yml up -d
```

Sign-up is enabled in the default image entrypoint. Mount `/friendle_data` so SQLite and the word pepper file survive restarts.

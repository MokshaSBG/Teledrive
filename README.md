# 📦 Teldrive — Telegram as Cloud Storage

> A self-hosted cloud storage service that uses Telegram as its storage backend, built with Go.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Status](https://img.shields.io/badge/Status-Early%20Development-orange)

Teldrive turns your Telegram account into a personal cloud drive — unlimited storage, no subscription fees, accessible from anywhere through a REST API, CLI, and web interface.

---

## ✨ Features

- 🔐 **Telegram Account Authentication** — Login with your own Telegram account, not a bot
- 📤 **Unlimited File Uploads** — Leverage Telegram's free unlimited storage
- 📥 **Fast Downloads** — Multi-stream download support for maximum speed
- 🗂️ **Virtual Folders** — Organize files with a database-backed folder structure
- 🔍 **Search & Filter** — Find files by name, type, size, or date
- 🔗 **Share Links** — Generate public links with optional expiry
- ♻️ **Deduplication** — SHA-256 checksums prevent duplicate uploads
- 🌐 **REST API** — Integrate with any application
- 💻 **CLI Client** — Manage files from the terminal
- 🎨 **Web UI** — Clean browser interface (coming soon)
- 🐳 **Docker Support** — One-command deployment
- 🔒 **Session Persistence** — Stay logged in across restarts

---

## 🏗️ Architecture

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   User      │─────▶│  Teldrive    │─────▶│  Telegram   │
│ (Web/CLI)   │◀─────│   Server     │◀─────│  Servers    │
└─────────────┘      └──────────────┘      └─────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │  Database    │
                     │ (SQLite/PG)  │
                     └──────────────┘
```

**How it works:**
- **Telegram** stores the actual file bytes (as messages in Saved Messages or a private channel)
- **Your database** tracks metadata (filename, size, checksum, message ID, folders)
- **Teldrive server** coordinates between users, storage, and metadata

For a deeper dive, see [ARCHITECTURE.md](./docs/ARCHITECTURE.md).

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** — [Download](https://go.dev/dl/)
- **Git**
- **Telegram API credentials** — Get them from [my.telegram.org/apps](https://my.telegram.org/apps)
- **SQLite** (bundled) or **PostgreSQL 14+** (optional)

### Installation

```bash
# Clone the repository
git clone https://github.com/MokshaSBG/teldrive.git
cd teldrive

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
nano .env
```

### Configuration

Edit `.env` with your Telegram credentials:

```env
# Telegram API (from https://my.telegram.org/apps)
TELEGRAM_API_ID=1234567
TELEGRAM_API_HASH=abcdef1234567890abcdef1234567890
TELEGRAM_PHONE=+1234567890

# Storage location: "self" for Saved Messages or a channel ID
TELEGRAM_STORAGE_PEER=self

# Database
DATABASE_URL=./teldrive.db

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Security
JWT_SECRET=change-me-to-a-random-string
ENCRYPTION_KEY=change-me-to-a-32-byte-key
```

### First Run

```bash
# Start the server (first run will prompt for Telegram verification code)
go run ./cmd/server

# Or build and run
make build
./bin/teldrive
```

On first run, you'll be prompted to enter the verification code sent to your Telegram app.

### Using Docker

```bash
# Build the image
docker build -t teldrive .

# Run the container
docker run -d \
  --name teldrive \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --env-file .env \
  teldrive
```

---

## 📖 Usage

### REST API

#### Upload a file

```bash
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@vacation.jpg"
```

#### List files

```bash
curl http://localhost:8080/api/files?limit=20 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### Download a file

```bash
curl -O http://localhost:8080/api/files/{id}/download \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### Delete a file

```bash
curl -X DELETE http://localhost:8080/api/files/{id} \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### CLI

```bash
# Login
teldrive login

# Upload
teldrive upload ./document.pdf

# List files
teldrive ls

# Download
teldrive download <file-id> ./downloads/

# Delete
teldrive rm <file-id>
```

### Web UI

Once the server is running, open:

```
http://localhost:8080
```

---

## 📁 Project Structure

```
teldrive/
├── cmd/
│   ├── server/              # HTTP API server
│   └── cli/                 # Command-line client
│
├── internal/
│   ├── telegram/            # Telegram integration (gotd/td)
│   ├── storage/             # Storage abstraction layer
│   ├── db/                  # Database models and repositories
│   ├── api/                 # HTTP routes, handlers, middleware
│   ├── service/             # Business logic
│   └── config/              # Configuration loading
│
├── pkg/
│   ├── types/               # Shared types
│   └── utils/               # Utilities
│
├── web/                     # Frontend (React)
├── migrations/              # Database migrations
├── scripts/                 # Development scripts
├── docs/                    # Documentation
│
├── .env.example
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

---

## 🛠️ Development

### Build

```bash
make build
```

### Run tests

```bash
make test
```

### Run with hot reload

```bash
make dev
```

### Format code

```bash
make fmt
```

### Lint

```bash
make lint
```

### Database migrations

```bash
make migrate-up
make migrate-down
```

---

## ⚠️ Important Considerations

### Telegram is NOT a Real Cloud Storage Provider

Before using this project, understand Telegram's limits:

| Limit | Value |
|:---|:---|
| **Max file size** | 2 GB (4 GB with Premium) |
| **Daily uploads** | ~100 files (varies by account age) |
| **Flood waits** | Can be hours if you upload too fast |
| **Terms of Service** | Automated file storage may violate ToS |

**Recommendations:**
- Use a **dedicated Telegram account**, not your main one
- **Enable 2FA** on that account
- Add **rate limiting** to your server
- Don't share your API credentials with anyone
- Test with small files first

### Security

- 🔐 Never commit `.env` files or session data to Git
- 🔒 Use strong `JWT_SECRET` and `ENCRYPTION_KEY` values
- 🛡️ Run behind a reverse proxy (Nginx/Caddy) with HTTPS
- 🚫 Don't expose the server directly to the internet without auth

### Legal

This project is for **personal use only**. Automated use of Telegram's API for mass storage may violate their Terms of Service. Use responsibly and at your own risk.

---

## 🗺️ Roadmap

### Phase 1 — Core (MVP)
- [x] Telegram account authentication
- [x] Upload files to Saved Messages
- [x] Download files
- [x] List files
- [x] SQLite metadata storage
- [x] Basic REST API

### Phase 2 — Usability
- [ ] Delete files
- [ ] Virtual folders
- [ ] Search & filtering
- [ ] File deduplication
- [ ] Upload/download progress
- [ ] Chunked uploads for large files
- [ ] User authentication (JWT)

### Phase 3 — Advanced
- [ ] Share links with expiry
- [ ] Multi-user support
- [ ] Quota management
- [ ] Web UI
- [ ] File previews
- [ ] Trash & restore

### Phase 4 — Production
- [ ] Docker deployment
- [ ] PostgreSQL support
- [ ] Multiple Telegram accounts
- [ ] Encryption at rest
- [ ] Backups
- [ ] Monitoring (Prometheus/Grafana)

---

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

```bash
# Fork the repo, then:
git checkout -b feature/amazing-feature
git commit -m "Add amazing feature"
git push origin feature/amazing-feature
```

Then open a Pull Request.

### Development Guidelines

- Write tests for new features
- Follow Go best practices ([Effective Go](https://go.dev/doc/effective_go))
- Run `make fmt` and `make lint` before committing
- Keep commits atomic and descriptive

---

## 📚 Documentation

- [Architecture Overview](./docs/ARCHITECTURE.md)
- [API Reference](./docs/API.md)
- [Configuration Guide](./docs/CONFIGURATION.md)
- [Deployment Guide](./docs/DEPLOYMENT.md)
- [Troubleshooting](./docs/TROUBLESHOOTING.md)

---

## 🙏 Acknowledgments

- [gotd/td](https://github.com/gotd/td) — Pure Go Telegram MTProto client
- [Teldrive](https://github.com/teldrive) — Original inspiration
- [Telegram](https://telegram.org) — For the unlimited storage


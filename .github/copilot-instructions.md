# SRCGI Copilot Instructions

## Project Overview
SRCGI is a Go web server/CGI application for tracking and visualizing SHOWROOM event participant scores and contribution points. It serves as both a standalone web server and can operate under Apache2/nginx as CGI.

## Architecture

### Multi-Component System
- **SRCGI**: Web interface (this repo) - event configuration, data visualization, contribution rankings
- **SRGSE5M**: Daemon for score data acquisition (separate repo)
- **SRGPC**: Daemon for per-slot listener contribution calculation (separate repo)
- **SRGCE**: Cron job for new event data and block event expansion (separate repo)

### Core Libraries (v2.x)
- `github.com/Chouette2100/srapi/v2`: SHOWROOM API wrapper
- `github.com/Chouette2100/srdblib/v2`: Database access library
- `github.com/Chouette2100/exsrapi/v2`: Common utilities
- `github.com/Chouette2100/srhandler/v2`: Additional web handlers

Replace paths with `replace` directive in `go.mod` for local development.

### Database
MySQL 8.0+ with optional SSH tunnel support via `sshql`. Key tables: `user`, `event`, `eventuser`, `points`, `contribution`, `accesslog`. Uses `gorp` ORM with `ExpandSliceArgs: true` for slice expansion.

Connection tuning in `main.go`:
```go
srdblib.Db.SetMaxOpenConns(8)
srdblib.Db.SetMaxIdleConns(12)
srdblib.Db.SetConnMaxLifetime(time.Minute * 5)
```

## Configuration

### Environment-Based Config
All config files use YAML with environment variable substitution (e.g., `${DBHOST}`).

**Critical**: Source `my_script.env` before running:
```bash
source ./my_script.env
./SRCGI
```

Key files:
- `ServerConfig.yml`: Web server, SSL, bot filtering, rate limiting
- `DBConfig.yml`: Database connection, SSH tunnel settings
- `Env.yml`: API rate limits (Lmin, Waitmsec)
- `bots.yml`: Bot user-agent patterns (regex-compiled in `ReadBots()`)
- `nontargetentry.yml`: Handlers excluded from bot filtering when `LvlBots == 2`

### Version Tracking
Version format: `YYMMPP` (6-digit, e.g., `200502`). Update in both `main.go` constant and `ShowroomCGIlib/ShowroomCGIlib.go`. Displayed as composite: `main_ShowroomCGIlib_srdblib_srapi`.

## Development Workflows

### Building & Running
```bash
# Build
go build -v .

# Run (requires my_script.env sourced)
source ./my_script.env
./SRCGI

# Or use run.sh
./run.sh
```

**Port binding**: For ports <1024 without root, use `setcap`:
```bash
sudo setcap cap_net_bind_service=+ep SRCGI
```

### Testing
Run tests requiring DB connection:
```bash
go test -v ./ShowroomCGIlib/... -run TestDrawLineGraph
```
Tests initialize DB with same connection settings as main app.

### Deployment Modes
1. **Standalone** (`WebServer: None`): Runs as HTTP/HTTPS server on configured port
2. **CGI** (`WebServer: Apache2Ubuntu` or `nginxSakura`): Paths adjusted for web server integration

Static files served from `public/` only in standalone mode. Templates in `templates/`.

## Critical Patterns

### Handler Middleware
All handlers wrapped in `commonMiddleware(rateLimiter, handler)` which:
1. Extracts real IP (handles `X-Forwarded-For` for proxies)
2. Rate limits per-IP (configurable `AccessLimit`/`TimeWindow`)
3. Filters bots based on `LvlBots` setting (0=none, 2=medium, 3=strict)
4. Logs to fail2ban file (`/var/log/SRCGI/fail2ban.log`)
5. Async writes to `accesslog` table via `Chlog` channel

Exempt handlers: `GraphSumDataHandler`, `GraphSumData1Handler`, `GraphSumData2Handler` (supporting data endpoints).

### Rate Limiting
`SimpleRateLimiter` tracks requests per IP. Exceeding limits = silent drop (no response). Default: 3 requests/second.

### Bot Filtering
Regex from `bots.yml` compiled to `ShowroomCGIlib.Regexpbots`. Level 3 blocks all bots unconditionally. Level 2 blocks bots except for handlers in `nontargetentry.yml`.

### Async Logging
Access logs written via channel to prevent handler blocking:
```go
ShowroomCGIlib.Chlog <- &al  // Non-blocking send to buffered channel
go ShowroomCGIlib.LogWorker() // Consumer goroutine
```

### Graph Generation
SVG graphs created with sequential numbering from channel `ShowroomCGIlib.Chimgfn` (0-999 rotation). Files written to `public/` (standalone) or web server public directory.

### Graceful Shutdown
Context-based shutdown on SIGINT/SIGTERM with 5-second timeout.

### Daily Log Rotation
`NewLogfileName()` goroutine rotates log files at midnight. Logs to both file and stdout when running in foreground terminal.

## Common Pitfalls

1. **Forget `source my_script.env`**: Config loading will fail with missing env vars
2. **Intervalmin != 5**: Code forces to 5 (legacy constraint from data acquisition daemon)
3. **SSH tunnel cleanup**: Always `defer srdblib.Dialer.Close()` when `UseSSH == true`
4. **Template execution**: Use `t.ExecuteTemplate(w, templateName, data)` - templates are pre-parsed
5. **Handler naming**: Recent migration from `HandlerXXX()` to `XXXHandler()` - both conventions exist
6. **FIXME in main.go:302**: GraphSum data handlers need independent rate limit consideration

## Key Files

- `main.go`: HTTP server setup, middleware, 50+ route handlers
- `ShowroomCGIlib/ShowroomCGIlib.go`: Core library, global config
- `ShowroomCGIlib/dblib.go`: Database operations (point lists, event info, user data)
- `SimpleRateLimiter.go`: Per-IP rate limiting
- `LogForFail2ban.go`: Fail2ban integration logging
- `ReadBots.go`/`ReadEntry.go`: Config file parsers

## External References

- [Zenn article series](https://zenn.dev/chouette2100/books/d8c28f8ff426b7): Setup guides, API documentation
- Live instance: `https://chouette2100.com/top`
- Related repos: See `public/index.html` for complete ecosystem links

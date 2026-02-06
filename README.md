# Asset Manager

A flexible asset management platform for tracking assets, their properties, and assignments to people.

## Architecture

- **API Server**: Go/Gin REST API with JWT authentication
- **Web Frontend**: Svelte/Bulma web application
- **Desktop App**: Wails/Svelte/Bulma desktop application

## Project Structure

```
asset_manager/
├── cmd/
│   ├── api/          # API server entry point
│   └── migrate/      # Database migration tool
├── internal/
│   ├── config/       # Configuration handling
│   ├── database/     # Database connection
│   ├── models/       # Data structs
│   ├── repository/   # Data access layer
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # Gin middleware (auth, etc.)
│   └── auth/         # JWT authentication
├── migrations/       # SQL migration files
├── web/              # Svelte web frontend
├── desktop/          # Wails desktop app
└── shared/           # Shared Svelte components
```

## Configuration

### API Server (`config.yaml`)

```yaml
server:
  port: 8084

database:
  host: localhost
  port: 3306
  user: asset_manager
  password: your_password
  name: asset_manager

jwt:
  secret: your-secret-key
  expiry_hours: 24
```

### Desktop App

On first run, the app will prompt for API configuration. Config is stored in:
- Linux: `~/.config/asset-manager/config.yaml`
- macOS: `~/Library/Application Support/asset-manager/config.yaml`
- Windows: `%APPDATA%/asset-manager/config.yaml`

## Building

```bash
# Build all
make all

# Build API server
make api

# Build web frontend
make web

# Build desktop app
make desktop

# Run database migrations
make migrate
```

## Development

```bash
# Run API server in development
make dev-api

# Run web frontend in development
make dev-web

# Run desktop app in development
make dev-desktop
```

## Database Setup

1. Create MySQL database:
```sql
CREATE DATABASE asset_manager;
CREATE USER 'asset_manager'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON asset_manager.* TO 'asset_manager'@'localhost';
```

2. Run migrations:
```bash
make migrate
```

## Default Users

After migration, a default admin user is created:
- Username: `admin`
- Password: `admin` (change immediately)

A special "Unassigned" person is also created for asset stock management.

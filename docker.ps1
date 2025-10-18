# MyLittlePrice Docker Management Script (PowerShell)
# Easy commands to manage the Docker stack

param(
    [Parameter(Position=0)]
    [string]$Command,

    [Parameter(Position=1)]
    [string]$Service
)

$COMPOSE_FILE = "docker-compose.yml"
$COMPOSE_PROD_FILE = "docker-compose.prod.yml"

function Print-Usage {
    Write-Host "MyLittlePrice Docker Manager" -ForegroundColor Blue
    Write-Host ""
    Write-Host "Usage: .\docker.ps1 [command]"
    Write-Host ""
    Write-Host "Development commands:"
    Write-Host "  up              Start all services (development)"
    Write-Host "  down            Stop all services"
    Write-Host "  restart         Restart all services"
    Write-Host "  logs [service]  View logs (optional: backend, frontend, redis, postgres)"
    Write-Host "  build           Rebuild all containers"
    Write-Host "  ps              Show running containers"
    Write-Host "  clean           Remove all containers and volumes (DANGEROUS!)"
    Write-Host ""
    Write-Host "Production commands:"
    Write-Host "  prod-up         Start all services (production)"
    Write-Host "  prod-down       Stop production services"
    Write-Host "  prod-logs       View production logs"
    Write-Host "  prod-build      Build production images"
    Write-Host ""
    Write-Host "Database commands:"
    Write-Host "  db-migrate      Run database migrations"
    Write-Host "  db-shell        Open PostgreSQL shell"
    Write-Host "  redis-cli       Open Redis CLI"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\docker.ps1 up"
    Write-Host "  .\docker.ps1 logs backend"
    Write-Host "  .\docker.ps1 prod-up"
}

switch ($Command) {
    # Development commands
    "up" {
        Write-Host "🚀 Starting MyLittlePrice (Development)..." -ForegroundColor Green
        docker-compose -f $COMPOSE_FILE up -d
        Write-Host "✅ Services started!" -ForegroundColor Green
        Write-Host "Frontend: http://localhost:3000"
        Write-Host "Backend: http://localhost:8080"
        Write-Host "Redis: localhost:6379"
        Write-Host "PostgreSQL: localhost:5432"
    }

    "down" {
        Write-Host "⏹️ Stopping services..." -ForegroundColor Blue
        docker-compose -f $COMPOSE_FILE down
        Write-Host "✅ Services stopped" -ForegroundColor Green
    }

    "restart" {
        Write-Host "🔄 Restarting services..." -ForegroundColor Blue
        docker-compose -f $COMPOSE_FILE restart
        Write-Host "✅ Services restarted" -ForegroundColor Green
    }

    "logs" {
        if ($Service) {
            docker-compose -f $COMPOSE_FILE logs -f $Service
        } else {
            docker-compose -f $COMPOSE_FILE logs -f
        }
    }

    "build" {
        Write-Host "🔨 Building containers..." -ForegroundColor Blue
        docker-compose -f $COMPOSE_FILE build --no-cache
        Write-Host "✅ Build complete" -ForegroundColor Green
    }

    "ps" {
        docker-compose -f $COMPOSE_FILE ps
    }

    "clean" {
        Write-Host "⚠️ WARNING: This will remove all containers and volumes!" -ForegroundColor Red
        $confirm = Read-Host "Are you sure? (yes/no)"
        if ($confirm -eq "yes") {
            docker-compose -f $COMPOSE_FILE down -v
            Write-Host "✅ Cleaned up" -ForegroundColor Green
        } else {
            Write-Host "Cancelled"
        }
    }

    # Production commands
    "prod-up" {
        Write-Host "🚀 Starting MyLittlePrice (Production)..." -ForegroundColor Green
        if (-not (Test-Path ".env")) {
            Write-Host "❌ Error: .env not found" -ForegroundColor Red
            Write-Host "Copy .env.example to .env and configure it"
            exit 1
        }
        docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD_FILE up -d
        Write-Host "✅ Production services started!" -ForegroundColor Green
    }

    "prod-down" {
        Write-Host "⏹️ Stopping production services..." -ForegroundColor Blue
        docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD_FILE down
        Write-Host "✅ Production services stopped" -ForegroundColor Green
    }

    "prod-logs" {
        docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD_FILE logs -f
    }

    "prod-build" {
        Write-Host "🔨 Building production images..." -ForegroundColor Blue
        docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD_FILE build --no-cache
        Write-Host "✅ Production build complete" -ForegroundColor Green
    }

    # Database commands
    "db-migrate" {
        Write-Host "📊 Running migrations..." -ForegroundColor Blue
        docker-compose -f $COMPOSE_FILE exec postgres psql -U postgres -d mylittleprice -f /docker-entrypoint-initdb.d/002_add_users.sql
        Write-Host "✅ Migrations complete" -ForegroundColor Green
    }

    "db-shell" {
        Write-Host "🐘 Opening PostgreSQL shell..." -ForegroundColor Blue
        docker-compose -f $COMPOSE_FILE exec postgres psql -U postgres -d mylittleprice
    }

    "redis-cli" {
        Write-Host "📦 Opening Redis CLI..." -ForegroundColor Blue
        docker-compose -f $COMPOSE_FILE exec redis redis-cli
    }

    default {
        Print-Usage
        exit 1
    }
}

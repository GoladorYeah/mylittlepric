#!/bin/bash

# MyLittlePrice Docker Management Script
# Easy commands to manage the Docker stack

set -e

COMPOSE_FILE="docker-compose.yml"
COMPOSE_PROD_FILE="docker-compose.prod.yml"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_usage() {
    echo -e "${BLUE}MyLittlePrice Docker Manager${NC}"
    echo ""
    echo "Usage: ./docker.sh [command]"
    echo ""
    echo "Development commands:"
    echo "  up              Start all services (development)"
    echo "  down            Stop all services"
    echo "  restart         Restart all services"
    echo "  logs [service]  View logs (optional: backend, frontend, redis, postgres)"
    echo "  build           Rebuild all containers"
    echo "  ps              Show running containers"
    echo "  clean           Remove all containers and volumes (DANGEROUS!)"
    echo ""
    echo "Production commands:"
    echo "  prod-up         Start all services (production)"
    echo "  prod-down       Stop production services"
    echo "  prod-logs       View production logs"
    echo "  prod-build      Build production images"
    echo ""
    echo "Database commands:"
    echo "  db-migrate      Run database migrations"
    echo "  db-shell        Open PostgreSQL shell"
    echo "  redis-cli       Open Redis CLI"
    echo ""
    echo "Examples:"
    echo "  ./docker.sh up"
    echo "  ./docker.sh logs backend"
    echo "  ./docker.sh prod-up"
}

case "$1" in
    # Development commands
    up)
        echo -e "${GREEN}🚀 Starting MyLittlePrice (Development)...${NC}"
        docker-compose -f $COMPOSE_FILE up -d
        echo -e "${GREEN}✅ Services started!${NC}"
        echo "Frontend: http://localhost:3000"
        echo "Backend: http://localhost:8080"
        echo "Redis: localhost:6379"
        echo "PostgreSQL: localhost:5432"
        ;;

    down)
        echo -e "${BLUE}⏹️  Stopping services...${NC}"
        docker-compose -f $COMPOSE_FILE down
        echo -e "${GREEN}✅ Services stopped${NC}"
        ;;

    restart)
        echo -e "${BLUE}🔄 Restarting services...${NC}"
        docker-compose -f $COMPOSE_FILE restart
        echo -e "${GREEN}✅ Services restarted${NC}"
        ;;

    logs)
        if [ -z "$2" ]; then
            docker-compose -f $COMPOSE_FILE logs -f
        else
            docker-compose -f $COMPOSE_FILE logs -f $2
        fi
        ;;

    build)
        echo -e "${BLUE}🔨 Building containers...${NC}"
        docker-compose -f $COMPOSE_FILE build --no-cache
        echo -e "${GREEN}✅ Build complete${NC}"
        ;;

    ps)
        docker-compose -f $COMPOSE_FILE ps
        ;;

    clean)
        echo -e "${RED}⚠️  WARNING: This will remove all containers and volumes!${NC}"
        read -p "Are you sure? (yes/no): " confirm
        if [ "$confirm" = "yes" ]; then
            docker-compose -f $COMPOSE_FILE down -v
            echo -e "${GREEN}✅ Cleaned up${NC}"
        else
            echo "Cancelled"
        fi
        ;;

    # Production commands
    prod-up)
        echo -e "${GREEN}🚀 Starting MyLittlePrice (Production)...${NC}"
        if [ ! -f .env ]; then
            echo -e "${RED}❌ Error: .env not found${NC}"
            echo "Copy .env.example to .env and configure it"
            exit 1
        fi
        docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD_FILE up -d
        echo -e "${GREEN}✅ Production services started!${NC}"
        ;;

    prod-down)
        echo -e "${BLUE}⏹️  Stopping production services...${NC}"
        docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD_FILE down
        echo -e "${GREEN}✅ Production services stopped${NC}"
        ;;

    prod-logs)
        docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD_FILE logs -f
        ;;

    prod-build)
        echo -e "${BLUE}🔨 Building production images...${NC}"
        docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD_FILE build --no-cache
        echo -e "${GREEN}✅ Production build complete${NC}"
        ;;

    # Database commands
    db-migrate)
        echo -e "${BLUE}📊 Running migrations...${NC}"
        docker-compose -f $COMPOSE_FILE exec postgres /bin/sh -c "
            export PGPASSWORD=postgres
            cd /docker-entrypoint-initdb.d
            if [ -f apply_migrations.sh ]; then
                sh apply_migrations.sh
            else
                echo 'Migration script not found, applying migrations manually...'
                for f in \$(ls -1 *.sql 2>/dev/null | sort); do
                    echo \"Applying \$f...\"
                    psql -U postgres -d mylittleprice -f \"\$f\"
                done
            fi
        "
        echo -e "${GREEN}✅ Migrations complete${NC}"
        ;;

    db-shell)
        echo -e "${BLUE}🐘 Opening PostgreSQL shell...${NC}"
        docker-compose -f $COMPOSE_FILE exec postgres psql -U postgres -d mylittleprice
        ;;

    redis-cli)
        echo -e "${BLUE}📦 Opening Redis CLI...${NC}"
        docker-compose -f $COMPOSE_FILE exec redis redis-cli
        ;;

    *)
        print_usage
        exit 1
        ;;
esac

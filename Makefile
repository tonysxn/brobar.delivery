.PHONY: up-dev up-dev-detached up-prod up-prod-detached down purge

up-dev:
	docker-compose -f docker-compose.dev.yml up

up-dev-rebuild:
	docker-compose -f docker-compose.dev.yml down --remove-orphans
	docker-compose -f docker-compose.dev.yml build --no-cache
	docker-compose -f docker-compose.dev.yml up

up-dev-reset:
	docker-compose -f docker-compose.dev.yml down --remove-orphans --volumes
	docker-compose -f docker-compose.dev.yml build --no-cache
	docker-compose -f docker-compose.dev.yml up

up-dev-detached:
	docker-compose -f docker-compose.dev.yml up -d

up-prod:
	docker-compose -f docker-compose.yml up

up-prod-rebuild:
	docker-compose -f docker-compose.yml down --remove-orphans
	docker-compose -f docker-compose.yml build --no-cache
	docker-compose -f docker-compose.yml up

up-prod-reset:
	docker-compose -f docker-compose.yml down --remove-orphans --volumes
	docker-compose -f docker-compose.yml build --no-cache
	docker-compose -f docker-compose.yml up

up-prod-detached:
	docker-compose -f docker-compose.yml up -d

down:
	docker-compose -f docker-compose.yml down --remove-orphans

purge:
	docker-compose -f docker-compose.yml down --volumes --remove-orphans
	docker volume prune -f

clean:
	docker-compose -f docker-compose.dev.yml down --volumes --remove-orphans
	docker-compose -f docker-compose.yml down --volumes --remove-orphans
	docker container prune -f
	docker volume prune -f
clean-frontend:
	./scripts/clean-next.sh
test:
	cd backend && go test -v ./tests/...

webhook-install:
	cd backend && go run cmd/webhook-tool/main.go install

webhook-remove:
	cd backend && go run cmd/webhook-tool/main.go remove

restore-data:
	./scripts/restore_data.sh
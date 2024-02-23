run:
	go run cmd/api/main.go

docker-up-dev:
	docker-compose -f docker-compose.dev.yml up -d --build

docker-down-dev:
	docker-compose -f docker-compose.dev.yml down

docker-up-local:
	docker-compose -f docker-compose.local.yml up -d --build

docker-down-local:
	docker-compose -f docker-compose.local.yml down

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-build:
	docker build -t rinha-de-backend-2024-q1-poc .

docker-tag:
	docker tag rinha-de-backend-2024-q1-poc matheusfbosa/rinha-de-backend-2024-q1-poc

docker-push:
	docker push matheusfbosa/rinha-de-backend-2024-q1-poc

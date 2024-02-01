delete-dangling-image:
	docker image prune -f

build-develop:
	docker-compose -f docker-compose.dev.yml build

run-develop:
	docker-compose -f docker-compose.dev.yml up -d

down-develop:
	docker-compose -f docker-compose.dev.yml down

restart-develop:
	$(MAKE) down-develop
	$(MAKE) run-develop
delete-dangling-image:
	docker image prune -f

git-pull:
	git pull

build-develop:
	docker-compose -f docker-compose.yml build

run-develop:
	docker-compose -f docker-compose.yml up -d

down-develop:
	docker-compose -f docker-compose.yml down

restart-develop:
	$(MAKE) git-pull
	$(MAKE) down-develop
	$(MAKE) run-develop
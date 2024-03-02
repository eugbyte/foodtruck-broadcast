## call the commands like this: `$ make workspace=<dir> <cmd>`, e.g. `$ make workspace=services/greet dev`
start:
	cd ${workspace} && make start

tidy:
	cd ${workspace} && go mod tidy

lint:
	cd ${workspace} && make lint

test:
	cd ${workspace} && make test

download:
	cd ${workspace} && make download

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin 
	golangci-lint --version

compose:
	docker compose build --no-cache
	docker compose up --detach --force-recreate
	docker system prune --force
	docker volume prune --force

compose-down:
	docker compose down
	docker system prune --force
	docker volume prune --force

compose-refresh:
	make compose-down
	make compose

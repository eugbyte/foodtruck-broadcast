start:
	npx kill-port 6000
	go run cmd/websocket/*.go

simulate:
	go run cmd/simulate/*.go

lint:
	golangci-lint run ./...

test:
	go test -timeout 5000s ./...

tidy:
	go mod tidy

download:
	go mod download

compose:
	cd ./../.. && make compose

compose-down:
	cd ./../.. && make compose-down

compose-refresh:
	cd ./../.. && make compose-refresh

tflocal:	
	tflocal -chdir=infra init
	terraform -chdir=infra validate
	terraform -chdir=infra fmt
	tflocal -chdir=infra destroy -auto-approve 
	tflocal -chdir=infra plan
	tflocal -chdir=infra apply -auto-approve 

tflocal-refresh:
	make compose-down
	rm -rf ./infra/.terraform
	rm -f ./infra/.terraform.lock.hcl ./infra/terraform.tfstate  ./infra/terraform.tfstate.backup
	npx kill-port 6000
	docker-compose --file ../../docker-compose.yml up -d localstack_pro
	make tflocal

docker:
	docker build --file service/stream_aggregator/Dockerfile  --tag foodtruck/stream_aggregator:v0.0.1 . --progress=plain --no-cache

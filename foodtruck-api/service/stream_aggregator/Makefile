start:
	npx kill-port 8080
	go run cmd/websocket/*.go

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

tflocal:	
	tflocal -chdir=infra init
	tflocal -chdir=infra validate
	tflocal -chdir=infra fmt
	tflocal -chdir=infra destroy -auto-approve 
	tflocal -chdir=infra plan
	tflocal -chdir=infra apply -auto-approve 

tflocal-refresh:
	make compose-down
	npx kill-port 8080
# if you don't remove these files, localstack runs into an error
	rm -rf ./infra/.terraform
	rm -f ./infra/.terraform.lock.hcl ./infra/terraform.tfstate  ./infra/terraform.tfstate.backup
	docker-compose --file ../../docker-compose.yml up -d localstack_pro
	make tflocal

docker:
	docker build --file service/stream_aggregator/Dockerfile  --tag foodtruck/stream_aggregator:v0.0.1 . --progress=plain --no-cache

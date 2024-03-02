lint:
	golangci-lint run ./...

test:
	go test -timeout 5000s ./...

build:
	env GO111MODULE=on GOOS=linux go build -o bin/subscription/main -ldflags="-s -w" cmd/subscription/main.go

windows-zip:
	go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest
	echo ${GOPATH}
	${GOPATH}/bin/build-lambda-zip.exe -o bin/subscription/main.zip bin/subscription/main
	
clear:
	rm -rf bin

tidy:
	go mod tidy

download:
	go mod download

dynamodb-dashboard:
	echo "navigate to http://localhost:8001/"
	DYNAMO_ENDPOINT=http://localhost:4566 \
	AWS_REGION=ap-southeast-1 \
	AWS_ACCESS_KEY_ID=123 AWS_SECRET_ACCESS_KEY=123 \
	npx dynamodb-admin

compose:
	cd ./../.. && make compose

compose-down:
	cd ./../.. && make compose-down

tflocal:	
	tflocal -chdir=infra init
	tflocal -chdir=infra validate
	tflocal -chdir=infra destroy -auto-approve 
	tflocal -chdir=infra plan
	tflocal -chdir=infra apply -auto-approve 

tflocal-refresh:
	make compose-down
	docker-compose up -d localstack_pro
	make tflocal

local-list-lambdas:
	aws --endpoint-url=http://localhost:4566 lambda list-functions

local-list-gws:
	aws --endpoint-url=http://localhost:4566 apigatewayv2 get-apis  

local-list-dbs:
	aws --endpoint-url=http://localhost:4566 dynamodb list-tables

local-invoke-subscription:
	aws lambda invoke --function-name subscription --endpoint-url=http://localhost:4566 output.json
	cat output.json && printf "\n"

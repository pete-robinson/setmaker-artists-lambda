.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/process-artists process-artists/lambda_handler.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

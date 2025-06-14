run-linux:
	go build -o ../bin/keylogme-zero ./main.go && sudo CONFIG_FILE=../deploy/default_config_linux.json ../bin/keylogme-zero
run-darwin:
	go build -o ../bin/keylogme-zero ./main.go && sudo CONFIG_FILE=../deploy/default_config_darwin.json.template ../bin/keylogme-zero

test:
	cd .. && go test  ./...

compile:
	GOOS=linux go build -o k0_linux .
	GOOS=darwin go build -o k0_darwin .

gofumpt:
	golangci-lint run -E gofumpt ..


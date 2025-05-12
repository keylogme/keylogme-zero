run:
	go build -o ../bin/keylogme-zero ./main.go && sudo CONFIG_FILE=../deploy/default_config.json ../bin/keylogme-zero

test:
	cd .. && go test  ./...

compile:
	GOOS=linux go build -o k0_linux .
	GOOS=darwin go build -o k0_darwin .

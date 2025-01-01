run:
	go build -o ./gokeny ./main.go && sudo CONFIG_FILE=deploy/default_config.json ./gokeny

test:
	go test  ./...


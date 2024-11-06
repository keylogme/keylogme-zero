ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go build -o ./gokeny ./main.go && sudo ./gokeny


ifneq (,$(wildcard ./.env))
    include .env
    export
endif

dev:
	go build -o ./gokeny ./main.go && sudo ./gokeny ${API_KEY_DEV} ${ENDPOINT_DEV}
run:
	go build -o ./gokeny ./main.go && sudo ./gokeny ${API_KEY_PROD} ${ENDPOINT_PROD}


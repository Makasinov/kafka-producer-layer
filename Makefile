#!make

PROJECT_PATH=`pwd`
PROJECT_NAME=kafka-producer
MAIN_FILE_PATH=cmd/kafka-producer

## run: Launch app in debug mode (with -race)
run:
	@echo " > Launch..."
	@bash -c "go run -race $(MAIN_FILE_PATH)/*.go config/kafka-producer.conf.json"

## install: Downloading dependencies from go.mod
install:
	@echo " > Downloading dependencies"
	@bash -c "go mod download"
	@echo " > Done"

## test: Launch unit tests
test:
	@echo " > Launch unit tests"
	@bash -c "go test ./... -v"
	@echo " > Done"

## reinstall: Update dependencies
reinstall:
	@rm go.sum
	@bash -c "make install"

help: Makefile
	@echo " > Command list "$(PROJECT_NAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

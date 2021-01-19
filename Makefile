BUILDPATH=$(CURDIR)
BINARY=dremio_scrapper
IMAGE_NAME=asdasdasd

build:
	@echo "Creando Binario ..."
	@go build  -ldflags '-s -w' -o $(BUILDPATH)/build/bin/${BINARY} cmd/${BINARY}/main.go
	@echo "Binario generado en build/bin/${BINARY}"

test:
	@echo "Ejecutando tests..."
	@go test ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out

coverage:
	@echo "Coverfile..."
	@go test ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out
	@go tool cover -html=coverfile_out -o coverfile_out.html

docker:
	@docker build . -t ${IMAGE_NAME}:latest -f Dockerfile

.PHONY: test build
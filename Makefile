
.PHONY: all generate fmt build test coverage.console coverage.html

all: image generate fmt build test

pwd := $(CURDIR)
generate = docker run --rm -v $(pwd):/pwd -v $(realpath ../models):/models -w /pwd masagroup/soft.generator.go -m /models/$(2) -o $(1) -ps /pwd/generator.properties

image:
	@docker build --file Dockerfile --tag masagroup/soft.go.dev .

generate:
	@echo "[generate]"
	@$(call generate,/pwd,ecore.ecore)
	@$(call generate,/pwd/test,library.ecore)
	@$(call generate,/pwd/test,tournament.ecore)
	@$(call generate,/pwd/test,empty.ecore)

fmt:
	@echo "[fmt]"
	@docker run --rm -v $(pwd):/pwd -w /pwd masagroup/soft.go.dev go fmt ./...

build:
	@echo "[build]"
	@docker run --rm -v $(pwd):/pwd -w /pwd masagroup/soft.go.dev go build ./...

test:
	@echo "[test]"
	@docker run --rm -p 8080:80 -v $(pwd):/pwd -w /pwd masagroup/soft.go.dev go test -covermode=atomic ./...

coverage.console:
	@echo "[coverage.console]"
	@docker run --rm -v $(pwd):/pwd -w /pwd masagroup/soft.go.dev \
			sh -c 'mkdir -p /pwd/coverage &&\
					go test -coverprofile /pwd/coverage/coverage.out ./... &&\
					go tool cover -func=/pwd/coverage/coverage.out'
coverage.html:
	@echo "[coverage.html]"
	@docker run --rm -v $(pwd):/pwd -w /pwd masagroup/soft.go.dev \
			sh -c 'mkdir -p /pwd/coverage &&\
					go test -coverprofile /pwd/coverage/coverage.out ./... &&\
					go tool cover -html=/pwd/coverage/coverage.out -o /pwd/coverage/coverage.html'

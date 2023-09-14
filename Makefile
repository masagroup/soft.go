# *****************************************************************************
# Copyright(c) 2021 MASA Group
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# *****************************************************************************

GENERATE = docker run --rm -v $(CURDIR):/pwd -v $(realpath ../models):/models -w /pwd masagroup/soft.generator.go -m /models/$(2) -o /pwd/$(1) -P /pwd/generator.properties $(3)

# os detection
ifeq (${OS},Windows_NT)
MKDIR = mkdir $(subst /,\,$(1)) > nul 2>&1 || (exit 0)
WHICH := where
DEVNULL := NUL
else
MKDIR = mkdir -p $(1)
WHICH := which
DEVNULL := /dev/null
endif

# detect go
ifneq ($(shell $(WHICH) go 2>$(DEVNULL)),)
	GO := go
else 
	ifneq ($(shell $(WHICH) go.exe 2>$(DEVNULL)),)
		GO := go.exe
	else
		$(error "go is not in your system PATH")
	endif
endif

.PHONY: all
all: generate fmt lint build test

.PHONY: generate 
generate:
	@echo "[generate]"
	@$(call GENERATE,,ecore.ecore,)
	@$(call GENERATE,test,empty.ecore,)
	@$(call GENERATE,test,library.ecore,)
	@$(call GENERATE,test,tournament.ecore,)
	@$(call GENERATE,,tournament.ecore, \
		-p featureDelegation=Reflective \
		-p module=github.com/masagroup/soft.go \
		-p packages=github.com/masagroup/soft.go/test/tournament-reflective/tournament#github.com/masagroup/soft.go/test/tournament-reflective/internal/impls/impls#github.com/masagroup/soft.go/test/tournament-reflective/internal/mocks/mocks)

.PHONY: fmt
fmt:
	@echo "[fmt]"
	@$(GO) fmt ./...

.PHONY: lint
lint:
	@echo "[lint]"
	@docker run --rm -v $(CURDIR):/pwd -w /pwd golangci/golangci-lint:v1.54.2-alpine golangci-lint run

.PHONY: build
build:
	@echo "[build]"
	@$(GO) build ./...

.PHONY: test
test:
	@echo "[test]"
	@$(GO) test -covermode=atomic ./...

.PHONY: coverage.console
coverage.console:
	@echo "[coverage.console]"
	@$(call MKDIR,coverage)
	@$(GO) test -coverprofile coverage/coverage.out ./...
	@$(GO) tool cover -func=coverage/coverage.out

.PHONY: coverage.html
coverage.html:
	@echo "[coverage.html]"
	@$(call MKDIR,coverage)
	@$(GO) test -coverprofile coverage/coverage.out ./...
	@$(GO) tool cover -html=coverage/coverage.out -o coverage/coverage.html

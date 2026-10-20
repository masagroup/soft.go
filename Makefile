# *****************************************************************************
# Copyright(c) 2021 MASA Group
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# *****************************************************************************

# the tournament model: one package shared by two implementations (typed fields and reflective
# delegation), with the interfaces and the mocks in test/tournament
TOURNAMENT = github.com/masagroup/soft.go/test/tournament
TOURNAMENT_TYPED = $(TOURNAMENT)/tournament\#$(TOURNAMENT)/typed/typed\#$(TOURNAMENT)/tournament
TOURNAMENT_REFLECTIVE = $(TOURNAMENT)/tournament\#$(TOURNAMENT)/reflective/reflective\#$(TOURNAMENT)/tournament

# the modules of the go.work: the core and the codecs, each with its own dependencies
MODULES = . codec/xml codec/xmi codec/json codec/binary codec/sqlite codec/defaults test/integration

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
	@$(call GENERATE,,tournament.ecore, \
		-t generateInterfaces -t generateMocks -t generateMocksTests \
		-p packageInstance=interfaces \
		-p module=github.com/masagroup/soft.go \
		-p packages=$(TOURNAMENT_TYPED))
	@$(call GENERATE,,tournament.ecore, \
		-t generateImplementations -t generateTests \
		-p packageInstance=interfaces \
		-p module=github.com/masagroup/soft.go \
		-p packages=$(TOURNAMENT_TYPED))
	@$(call GENERATE,,tournament.ecore, \
		-t generateImplementations -t generateTests \
		-p packageInstance=interfaces \
		-p featureDelegation=Reflective \
		-p module=github.com/masagroup/soft.go \
		-p packages=$(TOURNAMENT_REFLECTIVE))

.PHONY: fmt
fmt:
	@echo "[fmt]"
	@$(foreach m,$(MODULES),$(GO) -C $(m) fmt ./... &&) echo done

.PHONY: lint
lint:
	@echo "[lint]"
	@docker run --rm -v $(CURDIR):/pwd -w /pwd golangci/golangci-lint:v2.14.0-alpine golangci-lint run --timeout=5m

.PHONY: build
build:
	@echo "[build]"
	@$(foreach m,$(MODULES),$(GO) -C $(m) build ./... &&) echo done

.PHONY: test
test:
	@echo "[test]"
	@$(foreach m,$(MODULES),$(GO) -C $(m) test -covermode=atomic ./... &&) echo done

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

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

# the generator: the docker image of soft.gen v2 at a fixed version
SOFT_GEN_VERSION = v2.0.0-alpha.3
SOFT_GEN_IMAGE ?= masagroup/soft.gen:$(SOFT_GEN_VERSION)

# the golang image go runs in when it is not installed: its Go must support the go directive of go.work
GOLANG_IMAGE ?= golang:1.27.3-alpine

# os detection
ifeq (${OS},Windows_NT)
MKDIR = mkdir $(subst /,\,$(1)) > nul 2>&1 || (exit 0)
WHICH := where
DEVNULL := NUL
else
MKDIR = mkdir -p $(1)
WHICH := which
DEVNULL := /dev/null
# the generated files belong to the user, not to root
SOFT_GEN_USER := --user $(shell id -u):$(shell id -g)
endif

# detect go: go, else go.exe (the Windows Go, from wsl)
ifneq ($(shell $(WHICH) go 2>$(DEVNULL)),)
	GO := go
else ifneq ($(shell $(WHICH) go.exe 2>$(DEVNULL)),)
	GO := go.exe
endif

# without go, go runs in a golang container: this folder mounted read-write (with its go.work),
# the go modules and the build cache in docker volumes kept between runs, no VCS stamping (no git
# in the image)
ifndef GO
	GO := MSYS_NO_PATHCONV=1 docker run --rm \
		-v "$(CURDIR)":/src -w /src \
		-v soft.go-gomod:/go/pkg/mod \
		-v soft.go-gocache:/root/.cache/go-build \
		-e GOFLAGS=-buildvcs=false \
		$(GOLANG_IMAGE) go
endif

# the soft.gen command: the soft.gen image, this folder mounted at /soft/go
# and the folder of the models where ../models points
SOFT_GEN = docker run --rm $(SOFT_GEN_USER) -v $(CURDIR):/soft/go -v $(realpath ../models):/soft/models -w /soft/go $(SOFT_GEN_IMAGE)
# GENERATE generates go code from a model of ../models with the properties of
# .generator.yaml: $(call GENERATE,<output folder>,<model>,<options>), the
# output folder relative to this one (empty: here), the model a file name in
# ../models, the options extra soft.gen options (-t <template>, -p key=value)
GENERATE = $(SOFT_GEN) go -m ../models/$(2) -o ./$(1) -P .generator.yaml $(3)

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

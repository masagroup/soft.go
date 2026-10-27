<!--
sword:
  schema: 1
  name: soft.go
  family: ecore-esync
  kind: go-lib
  generation: next-gen
  shipped: false
  languages: [Go, Makefile, YAML]
  build: [go, make]
  ng: [history, hub, recorder, datacore, esync, ecore-model, hla]
  apps: []
-->

For the previous documentation of the project (v1), see [README.LEGACY.md](README.LEGACY.md).

# soft.go

> Go implementation of the Eclipse EMF Ecore metamodel runtime (EObject, EClass, EResource, notifications) with model codecs in separate Go modules (XMI, XML, binary, JSON, SQLite) and a SQLite-backed `EStore`; it is the model foundation library of the Next-Gen Go services and of the esync library `soft.sync.go`.

| | |
|---|---|
| **Family** | Ecore / esync (the [Soft](https://github.com/masagroup/soft) family) |
| **Languages** | Go (100% of sources), Makefile, YAML |
| **Build systems** | Go (`go.work`) · Makefile (generation with soft.gen v2, run from its Docker image `masagroup/soft.gen:v2.0.0-alpha.4`; `go` runs in a `golang` container when Go is missing). No Ant, no Dockerfile, no docker-compose, no CMake, no npm |

**[Architecture](ARCHITECTURE.md)** · [Contributing](CONTRIBUTING.md) · [previous README (v1)](README.LEGACY.md) · [known issues and open questions](https://claude.ai/artifact/UfvKmky7e8LwEXf192hQjQ#soft.go) · [repository on GitHub](https://github.com/masagroup/soft.go)

## Overview

soft.go provides the Go runtime of the Ecore metamodel: the generated metamodel classes (`EClass`, `EPackage`, `EAttribute`, ...), several `EObject` flavours (basic, compact, dynamic, reflective, store-backed), collections, EMF-style notifications and adapters, resources and resource sets, URI handlers, ID managers, and the codec abstraction. The codecs that read and write model resources are separate packages and modules under `codec/` (see [Codecs](#codecs)). `codec/sqlite` also ships `sqlite.Store` (named `SQLStore` up to v1.12), an `EStore` implementation where every feature access of a store-backed `EObject` becomes a SQL query on a SQLite database.

<details>
<summary>More</summary>

The repository is a pure library: no `package main`, no network listener, no `flag` / `os.Getenv` usage in non-test code. It is consumed with `go get github.com/masagroup/soft.go` plus `go get` of the codec modules it uses (Go module proxy / GitHub tags). Module paths: `github.com/masagroup/soft.go` (package `ecore`), `github.com/masagroup/soft.go/codec/{xml,xmi,binary,sqlite,json,defaults}`, and `github.com/masagroup/soft.go/test/integration` (tests only), all `go 1.27.0`, tied together by `go.work`.

Who uses it (from their `go.mod`):

| Consumer | Version required |
|---|---|
| `soft/soft.sync.go` (esync) | v1.12.0 on `master`; v1.13.3 and the `codec/binary`, `codec/xmi`, `codec/xml` v1.0.0 modules on `release/v2` (module `soft.sync.go/v2`) |
| `soft/soft.sync.studio/devserver` | v1.13.3 (codec modules indirect) |
| `soft/soft.gen` (v2 generator) | v1.13.3-0.20261022175029-abd1d4df298b, with `codec/xmi` and `codec/xml` pseudo-versions (loads the `.ecore` models) |
| `sw/sw.ng/src` | v1.11.0 |
| `sw/sw.datacore` | v1.11.0, replaced by fork `github.com/BenSouchet/soft.go v1.10.3-0.20261020150046-020a766d99c1` |
| `sw/sw.analysis/src`, `sw/sw.api.ng.go` | v1.11.0; v1.13.3 and `codec/xml` v1.0.0 on the `release/v2` branch of sw.api.ng.go (module `sw.api.ng.go/v2`) |
| `sw/sw.recorder` | v1.13.3 and `codec/xml` v1.0.0 |
| `sw/sw.hub` | v1.13.3 (indirect, through `soft.sync.go/v2` and `sw.api.ng.go/v2`) |
| `io/io.connector`, `io/io.nato` | v1.11.0 |
| `io/io.hla` | v1.10.2 |
| `soft/soft` (umbrella repo, submodule `go`) | v1.7.0 (stale) |

The sw.* / io.* consumers of the v1 line are on v1.11.0 (io.hla on v1.10.2) while HEAD is v1.13.3 (2026-10-23); the v2 line (sw.hub, sw.recorder, the `release/v2` branches of soft.sync.go and sw.api.ng.go, the soft.sync.studio devserver) is on v1.13.3 with the codec modules. Since v1.13.0 the codecs are no longer in `ecore`: a consumer upgrading from v1.10 / v1.11 / v1.12 must require the codec modules and rename the codec symbols (see [Migration from v1.12](#migration-from-v112)).

</details>

The diagram of this overview is in [Architecture](ARCHITECTURE.md#overview).

## Build

The only build system is the root `Makefile` plus the standard Go toolchain (GNU make required; the Makefile detects `go` / `go.exe` and Windows vs POSIX `mkdir`, so it runs on both Windows and Linux).

### Linux (Docker)

Not applicable: no `Dockerfile` or `docker-compose.yaml`. Docker is used by `make generate` (image `masagroup/soft.gen:v2.0.0-alpha.4`, always), by `make lint` (image `golangci/golangci-lint:v2.14.0-alpine`), and to run `go` in a container (image `golang:1.27.3-alpine`) when Go is not installed (see below).

### Makefile (Windows and Linux)

Run from the repository root:

```powershell
# full pipeline as in CI (needs Docker, Go or the golang container, and the soft.models repo checked out as ../models)
make generate build test

# lint (Docker)
make lint
```

<details>
<summary>Details</summary>

Makefile targets:

| Target | What it does |
|---|---|
| `all` (default) | `generate fmt lint build test` |
| `generate` | Runs `soft.gen go` 6 times on `.ecore` models from `../models` with `-P .generator.yaml`: `ecore.ecore` -> repo root (regenerates `ecore/`); `empty.ecore`, `library.ecore` -> `test/`; `tournament.ecore` 3 times with `-p packageInstance=interfaces -p module=github.com/masagroup/soft.go` and the package paths of `-p packages=...`: interfaces and mocks (`-t generateInterfaces -t generateMocks -t generateMocksTests`) into `test/tournament`, typed implementation (`-t generateImplementations -t generateTests`) into `test/tournament/typed`, reflective implementation (same templates plus `-p featureDelegation=Reflective`) into `test/tournament/reflective` |
| `fmt` | `go -C <module> fmt ./...` for each module of `MODULES` |
| `lint` | `docker run golangci/golangci-lint:v2.14.0-alpine golangci-lint run --timeout=5m` |
| `build` | `go -C <module> build ./...` for each module of `MODULES` |
| `test` | `go -C <module> test -covermode=atomic ./...` for each module of `MODULES` |
| `coverage.console` / `coverage.html` | coverage profile of the root module in `coverage/` (gitignored) |

Makefile variables:

| Variable | Value |
|---|---|
| `MODULES` | `. codec/xml codec/xmi codec/json codec/binary codec/sqlite codec/defaults test/integration` |
| `SOFT_GEN_VERSION` | `v2.0.0-alpha.4` (`v2.0.0-alpha.2` up to v1.13.1, `v2.0.0-alpha.3` in v1.13.2) |
| `SOFT_GEN_IMAGE` | `masagroup/soft.gen:$(SOFT_GEN_VERSION)` (overridable with `?=`) |
| `SOFT_GEN` | `docker run --rm $(SOFT_GEN_USER) -v $(CURDIR):/soft/go -v $(realpath ../models):/soft/models -w /soft/go $(SOFT_GEN_IMAGE)`: the generator runs from its image, the repository mounted on `/soft/go` and the models folder on `/soft/models`; on Linux `SOFT_GEN_USER` (`--user <uid>:<gid>`) gives the generated files to the user. Up to v1.13.1 the Makefile ran `go run github.com/masagroup/soft.gen/v2/cmd/soft-gen@v2.0.0-alpha.2` instead (the command is `cmd/soft-gen` up to soft.gen `v2.0.0-alpha.2`, `cmd/soft.gen` from `v2.0.0-alpha.3`) |
| `GO` | the `go` found (`go`, else `go.exe`); without Go, `docker run --rm` of `GOLANG_IMAGE` with the repository mounted on `/src`, the Go module cache and build cache in the volumes `soft.go-gomod` and `soft.go-gocache`, and `GOFLAGS=-buildvcs=false`: `fmt`, `build`, `test` and the coverage targets then run in the container |
| `GOLANG_IMAGE` | `golang:1.27.3-alpine` (overridable with `?=`; its Go must support the `go` directive of `go.work`) |
| `GENERATE` | `$(call GENERATE,<output folder>,<model>,<options>)` = `$(SOFT_GEN) go -m ../models/<model> -o ./<output folder> -P .generator.yaml <options>` (extra options: `-t <template>`, `-p key=value`) |

Gotchas:

- The `test/` tree mostly contains generated, gitignored code (`test/*/.gitignore` ignores `*.go` and `go.sum`, whitelisting only hand-written files: `test/library/writer_ext.go` and `test/tournament/tournamentpackage_impl.go`; in `test/tournament` the `*/.soft.gen.yaml` of the implementations are ignored too). The generated models belong to the root module, and the integration tests of `test/integration` import them, so `make build` / `make test` (and `go build ./...` in the root module) on a fresh clone only work after `make generate`.
- `make generate` reads the models from `../models` (it mounts `$(realpath ../models)` on `/soft/models` in the soft.gen container). In the ecosystem workspace the models repo lives at `..\soft.models`, not `..\models` (CI checks it out as `models`). Running `make generate` locally therefore needs a `models` folder or link next to the repo (UNVERIFIED how developers handle this).
- The mockery mocks are generated with the config in `.mockery.yaml` (headers say mockery v2.53.2); the mockery version is not pinned in the Makefile.

</details>

### Go (Windows and Linux)

Run from the repository root:

```powershell
# core and codecs only (ecore/ and codec/ are committed, no generation needed)
go build ./ecore/...
go test -covermode=atomic ./ecore/...
go -C codec/sqlite test -covermode=atomic ./...   # same for codec/xml, codec/xmi, codec/binary, codec/json
```

## Quickstart

Not applicable: the repository has no executable and no service. The only I/O performed by the library is:

```go
import (
  "github.com/masagroup/soft.go/ecore"
  _ "github.com/masagroup/soft.go/codec/defaults" // XMI, XML, binary and SQLite codecs, registered by extension
)
```

<details>
<summary>Details</summary>

- file system via `FileURIHandler` (`ecore/fileurihandler.go`) and an in-memory `MemoryURIHandler` (`ecore/memoryurihandler.go`, `memory:` scheme);
- SQLite database files via `zombiezen.com/go/sqlite` (pure Go, no CGO) in `codec/sqlite`; the SQLite codec creates temporary `.sqlite` files in `os.TempDir()` when decoding from an `io.Reader`, encoding to an `io.Writer`, and in `sqlite.Store.Serialize` (`sqlTmpDB`, `codec/sqlite/sqlcodec.go`).

It opens no port and reads no environment variable.

Minimal usage:

</details>

## Interfaces and access methods

No executable is produced by this repository (library, contract, data or tooling): see the build and the dependencies.

## Documentation and guides

- [Architecture](ARCHITECTURE.md): position in the ecosystem, repository layout, internal structure, how it works (diagrams), dependencies.
- [Contributing](CONTRIBUTING.md): development environment, Git and Jira workflow, conventions, CI/CD.
- [Previous README (v1)](README.LEGACY.md): the documentation of the repository before the wiki rewrite.
- [Known issues and open questions](https://claude.ai/artifact/UfvKmky7e8LwEXf192hQjQ#soft.go): the private artifact that holds them (every page of the wiki).
- Related repositories: [io.connector](https://github.com/masagroup/io.connector), [io.hla](https://github.com/masagroup/io.hla), [io.nato](https://github.com/masagroup/io.nato), [soft](https://github.com/masagroup/soft), [soft.gen](https://github.com/masagroup/soft.gen), [soft.models](https://github.com/masagroup/soft.models), [soft.sync.go](https://github.com/masagroup/soft.sync.go), [soft.sync.studio](https://github.com/masagroup/soft.sync.studio), [sw.analysis](https://github.com/masagroup/sw.analysis), [sw.api.ng.go](https://github.com/masagroup/sw.api.ng.go), [sw.datacore](https://github.com/masagroup/sw.datacore), [sw.hub](https://github.com/masagroup/sw.hub), [sw.ng](https://github.com/masagroup/sw.ng), [sw.recorder](https://github.com/masagroup/sw.recorder).
- Transverse pages: [Build on Windows](https://github.com/masagroup/wiki/wiki/build-windows), [Build on Linux](https://github.com/masagroup/wiki/wiki/build-linux), [CI/CD pipelines](https://github.com/masagroup/wiki/wiki/cicd-pipelines), [Ecosystem today](https://github.com/masagroup/wiki/wiki/ecosystem-current), [Ecosystem next gen](https://github.com/masagroup/wiki/wiki/ecosystem-nextgen), [Glossary](https://github.com/masagroup/wiki/wiki/glossary).
- Repository on GitHub: [masagroup/soft.go](https://github.com/masagroup/soft.go).
- The same documentation in the SWORD wiki: [soft.go](https://github.com/masagroup/wiki/wiki/soft.go), [architecture](https://github.com/masagroup/wiki/wiki/soft.go-architecture), [contributing](https://github.com/masagroup/wiki/wiki/soft.go-contributing), [previous README](https://github.com/masagroup/wiki/wiki/soft.go-legacy).

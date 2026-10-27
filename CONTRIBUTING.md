# soft.go: contributing

> How to set up a development environment, the Git and Jira workflow and the conventions of `soft.go`. Back to the [README](README.md) · [Architecture](ARCHITECTURE.md)

## Development environment

The build of the repository is described in the [README](README.md#build); the tooling of the ecosystem is on the transverse pages:

- Windows, Ant and sw.poney: Apache Ant, JDK 8, Visual Studio 2022 Professional with the MSVC v140 toolset and the Windows SDK 10.0.19041.0, `gh` authenticated, access to `\\file1` or a local `sw.poney` clone in `PONEY_HOME`; initialise the MSVC environment before any `ant` (see [Build on Windows](https://github.com/masagroup/wiki/wiki/build-windows)).
- Linux and Docker Compose: Docker with BuildKit, the GitHub PAT file named by `SECRET_GITHUB_PAT` (`~/.masabot/github.pat`) or the SSH key of `SECRET_GITHUB_SSH` for the private Go modules and npm packages, the profiles `prod`, `test`, `pack` and `dev` (see [Build on Linux](https://github.com/masagroup/wiki/wiki/build-linux)); never push an image, the CI does it.
- Go: the version of `go.mod`, `GOPRIVATE=github.com/masagroup/*` with `gh auth setup-git` or an SSH `insteadOf` rewrite for the private modules; `staticcheck` installed by the build.
- Editor rules for Claude Code: the shared `claude/` folder of [cicd](https://github.com/masagroup/cicd) installs the same `CLAUDE.md` for every developer (rules, roles, build pages, conventions).

## Git and Jira workflow

- The organisation repository `masagroup/soft.go` is the reference (CI, releases); work on a fork or on a branch of it, never directly on `master` nor on a release branch (`6.28`...).
- Branch name: `<TICKET>-<short-description>` (for example `SWNG-421-tests-and-fixes`), created from the target branch of the ticket.
- Commit message: `<type>(<scope>): <TICKET-ID>: <imperative description in English>`, with `type` in `feat`, `fix`, `refactor`, `test`, `build`, `ci`, `perf`, `docs`, `chore`; the body says what and why, and the CI impact when there is one. One commit is one coherent change; no `Co-Authored-By` nor other automatic signature.
- Pull request: masabot reads the Jira keys (`SW`, `SWNG`, `IO`, `AI`, `GEFORA`, `COA`) in the title, the branch, the body and the commits, comments the Jira issue and moves it (*In Progress*, *In Review* with `[DONE]`, *In QA* or *In Validation* after the merge). The PR is built rebased on its base branch; the comment commands `/rebase`, `/rebase-all` and `/automerge` are available. See [CI/CD pipelines](https://github.com/masagroup/wiki/wiki/cicd-pipelines).
- Never push to `master` by hand, never create a release or push an image: the CI does it.

## Conventions

- Code, comments, commit messages and file names in English; documentation of the repository in English.
- Formatting of the repository (`clang-format`, `gofmt`/`goimports`, `prettier`/`eslint`) on the lines you touch; no reformatting of a whole file for a one-line change; no dead code, no `TODO` without a ticket; logs at the right level, never a secret.
- Go: the version of `go.mod` and of the Docker image, `GOPRIVATE=github.com/masagroup/*`; `go vet`, `staticcheck`, `gofmt` and `go test -race ./...` green before a commit; wrapped errors (`fmt.Errorf("...: %w", err)`), no `log.Fatal` outside `main`, clean stop on `SIGINT` and `SIGTERM`, timeouts on every network I/O; REST APIs from an OpenAPI contract (`oapi-codegen`).
- Shared contracts (`.proto`, `.ecore`, OpenAPI, exported Go modules, file formats) are changed with the consumers in mind: grep the workspace, state the compatibility strategy, update the dependants.
- Documentation: this README set (README, ARCHITECTURE, CONTRIBUTING, README.LEGACY) is the source of the [SWORD wiki](https://github.com/masagroup/wiki/wiki/soft.go); keep the metadata block at the top of the README up to date (the `apps` list feeds the Interfaces section).

## CI/CD

Single workflow: `.github/workflows/build_and_test.yaml` ("Build and Test").

| Item | Value |
|---|---|
| Triggers | pull request to `master` (opened, synchronize, reopened, ready_for_review), push to `master`, `workflow_dispatch` |
| Matrix | `go-version: [1.27.x]`, `os: [ubuntu-latest]` |
| Steps | `actions/setup-go@v5`; checkout this repo into `go/`; checkout `masagroup/soft.models` at pinned SHA `90286c344159492747920472760ac3904b19e784` into `models/`; run `make generate build test` in `go/` (soft.gen runs from the image `masagroup/soft.gen:v2.0.0-alpha.4` with the Docker of the runner, build and test with the installed Go; build and test cover every module of `MODULES`) |
| Shared workflows | None: it does not use the `masagroup/cicd` reusable workflows (contrary to the ecosystem convention); it is a self-contained job. See [cicd](https://github.com/masagroup/cicd) and the ecosystem wiki page [CI/CD pipelines](https://github.com/masagroup/wiki/wiki/cicd-pipelines) |
| Lint / coverage | No lint step in CI (`make lint` is not invoked), no coverage upload |
| Artifacts | None: no zip, deb, Docker image or release asset. Releases are plain git tags consumed through the Go module proxy: `v1.13.3`, `v1.13.2`, `v1.13.1`, `v1.13.0`, `v1.12.0`, `v1.11.1`, ... for the root module, `codec/<name>/vX.Y.Z` for the codec modules (`codec/xml/v1.0.0`, `codec/xmi/v1.0.0`, `codec/binary/v1.0.0`, `codec/sqlite/v1.0.0`, `codec/defaults/v1.0.0`, `codec/json/v1.1.0`) |

## Known issues and open questions

They are not kept in the repository: see the [private artifact](https://claude.ai/artifact/UfvKmky7e8LwEXf192hQjQ#soft.go) that holds them for every page of the wiki.

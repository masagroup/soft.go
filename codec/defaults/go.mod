module github.com/masagroup/soft.go/codec/defaults

go 1.27.0

require (
	github.com/masagroup/soft.go/codec v0.0.0-00010101000000-000000000000
	github.com/masagroup/soft.go/codec/binary v0.0.0-00010101000000-000000000000
	github.com/masagroup/soft.go/codec/sql v0.0.0-00010101000000-000000000000
)

require (
	github.com/chebyrash/promise v0.0.0-20260915144232-ce593674804a // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/masagroup/soft.go v1.13.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.1.0 // indirect
	github.com/oklog/ulid/v2 v2.1.2 // indirect
	github.com/panjf2000/ants/v2 v2.12.1 // indirect
	github.com/petermattis/goid v0.0.0-20260918085751-abfca077860b // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rqlite/sql v0.0.0-20261016220845-5c0832fea6e1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/stretchr/objx v0.8.0 // indirect
	github.com/stretchr/testify v1.12.1 // indirect
	github.com/ugurcsen/gods-generic v0.10.4 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	modernc.org/libc v1.79.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
	modernc.org/sqlite v1.60.1 // indirect
	zombiezen.com/go/sqlite v1.4.2 // indirect
)

replace (
	github.com/masagroup/soft.go => ../..
	github.com/masagroup/soft.go/codec => ../../codec
	github.com/masagroup/soft.go/codec/binary => ../../codec/binary
	github.com/masagroup/soft.go/codec/json => ../../codec/json
	github.com/masagroup/soft.go/codec/sql => ../../codec/sql
)

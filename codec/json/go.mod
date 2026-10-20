module github.com/masagroup/soft.go/codec/json

go 1.27.0

require (
	github.com/karlseguin/jsonwriter v1.0.4-0.20170525085137-6f05566bac1c
	github.com/masagroup/soft.go v1.13.0
	github.com/masagroup/soft.go/codec/xmi v0.0.0-00010101000000-000000000000
	github.com/masagroup/soft.go/codec/xml v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/karlseguin/expect v1.0.8 // indirect
	github.com/oklog/ulid/v2 v2.1.2 // indirect
	github.com/stretchr/objx v0.8.0 // indirect
	github.com/ugurcsen/gods-generic v0.10.4 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)

replace (
	github.com/masagroup/soft.go => ../..
	github.com/masagroup/soft.go/codec/binary => ../../codec/binary
	github.com/masagroup/soft.go/codec/defaults => ../../codec/defaults
	github.com/masagroup/soft.go/codec/sql => ../../codec/sql
	github.com/masagroup/soft.go/codec/xmi => ../../codec/xmi
	github.com/masagroup/soft.go/codec/xml => ../../codec/xml
)

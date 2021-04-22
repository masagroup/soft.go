
all: image ecore library tournament empty


pwd := $(shell pwd)

image:
	@docker build --file Dockerfile --tag masagroup/soft.go .

ecore: ecore.gen ecore.fmt ecore.build ecore.tests

ecore.gen:
	@echo "[ecore.gen]"
	@docker run --rm -v $(pwd):/pwd -v $(realpath ../models):/models -w /pwd masagroup/soft.generator.go -m /models/ecore.ecore -o /pwd/ -ps /pwd/generator.properties

ecore.fmt:
	@echo "[ecore.fmt]"
	@docker run --rm -v $(pwd):/pwd -w /pwd masagroup/soft.go go fmt ./...

ecore.build:
	@echo "[ecore.build]"
	@docker run --rm -v $(pwd):/pwd -w /pwd masagroup/soft.go go build ./...

ecore.tests:
	@echo "[ecore.tests]"
	@docker run --rm -v $(pwd):/pwd -w /pwd --env CGO_ENABLED=0 masagroup/soft.go go test ./...

library: library.gen library.fmt library.build library.tests

library.gen:
	@echo "[library.gen]"
	@docker run --rm -v $(pwd):/pwd -v $(realpath ../models):/models -w /pwd masagroup/soft.generator.go -m /models/library.ecore -o /pwd/test -ps /pwd/generator.properties

library.fmt:
	@echo "[library.fmt]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/library masagroup/soft.go go fmt ./...

library.build:
	@echo "[library.build]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/library masagroup/soft.go go build ./...

library.tests:
	@echo "[library.tests]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/library --env CGO_ENABLED=0 masagroup/soft.go go test ./...

tournament: tournament.gen tournament.fmt tournament.build tournament.tests

tournament.gen:
	@echo "[tournament.gen]"
	@docker run --rm -v $(pwd):/pwd -v $(realpath ../models):/models -w /pwd masagroup/soft.generator.go -m /models/tournament.ecore -o /pwd/test -ps /pwd/generator.properties

tournament.fmt:
	@echo "[tournament.fmt]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/tournament masagroup/soft.go go fmt ./...

tournament.build:
	@echo "[tournament.build]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/tournament masagroup/soft.go go build ./...

tournament.tests:
	@echo "[tournament.tests]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/tournament --env CGO_ENABLED=0 masagroup/soft.go go test ./...

empty: empty.gen empty.fmt empty.build empty.tests

empty.gen:
	@echo "[empty.gen]"
	@docker run --rm -v $(pwd):/pwd -v $(realpath ../models):/models -w /pwd masagroup/soft.generator.go -m /models/empty.ecore -o /pwd/test -ps /pwd/generator.properties

empty.fmt:
	@echo "[empty.fmt]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/empty masagroup/soft.go go fmt ./...

empty.build:
	@echo "[empty.build]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/empty masagroup/soft.go go build ./...

empty.tests:
	@echo "[empty.tests]"
	@docker run --rm -v $(pwd):/pwd -w /pwd/test/empty --env CGO_ENABLED=0 masagroup/soft.go go test ./...

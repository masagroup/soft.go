
all: image ecore library tournament empty


pwd := $(shell pwd)
generate = docker run --rm -v $(pwd):/pwd -v $(realpath ../models):/models -w /pwd masagroup/soft.generator.go -m /models/$(2) -o $(1) -ps /pwd/generator.properties
go_fmt   = docker run --rm -v $(pwd):/pwd -w $(1) masagroup/soft.go go fmt ./...
go_build = docker run --rm -v $(pwd):/pwd -w $(1) masagroup/soft.go go build ./...
go_test  = docker run --rm -v $(pwd):/pwd -w $(1) --env CGO_ENABLED=0 masagroup/soft.go\
					 sh -c 'mkdir -p /pwd/coverage &&\
					 		go test -coverprofile /pwd/coverage/$(2)-coverage ./... &&\
				   		    go tool cover -html=/pwd/coverage/$(2)-coverage -o /pwd/coverage/$(2)-coverage.html'

image:
	@docker build --file Dockerfile --tag masagroup/soft.go .

ecore: ecore.gen ecore.fmt ecore.build ecore.tests

ecore.gen:
	@echo "[ecore.gen]"
	@$(call generate,/pwd,ecore.ecore)

ecore.fmt:
	@echo "[ecore.fmt]"
	@$(call go_fmt,/pwd)
	
ecore.build:
	@echo "[ecore.build]"
	@$(call go_build,/pwd)
	
ecore.tests:
	@echo "[ecore.tests]"
	@$(call go_test,/pwd,ecore)

library: library.gen library.fmt library.build library.tests

library.gen:
	@echo "[library.gen]"
	@$(call generate,/pwd/test,library.ecore)

library.fmt:
	@echo "[library.fmt]"
	@$(call go_fmt,/pwd/test/library)

library.build:
	@echo "[library.build]"
	@$(call go_build,/pwd/test/library)

library.tests:
	@echo "[library.tests]"
	@$(call go_test,/pwd/test/library,library)

tournament: tournament.gen tournament.fmt tournament.build tournament.tests

tournament.gen:
	@echo "[tournament.gen]"
	@$(call generate,/pwd/test,tournament.ecore)

tournament.fmt:
	@echo "[tournament.fmt]"
	@$(call go_fmt,/pwd/test/tournament)

tournament.build:
	@echo "[tournament.build]"
	@$(call go_build,/pwd/test/tournament)

tournament.tests:
	@echo "[tournament.tests]"
	@$(call go_test,/pwd/test/tournament,tournament)

empty: empty.gen empty.fmt empty.build empty.tests

empty.gen:
	@echo "[empty.gen]"
	@$(call generate,/pwd/test,empty.ecore)

empty.fmt:
	@echo "[empty.fmt]"
	@$(call go_fmt,/pwd/test/empty)

empty.build:
	@echo "[empty.build]"
	@$(call go_build,/pwd/test/empty)

empty.tests:
	@echo "[empty.tests]"
	@$(call go_test,/pwd/test/empty,empty)


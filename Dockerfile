FROM golang:1.16.2-alpine3.13

# module dependencies
WORKDIR /build
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
RUN go mod download

#vscode tools
RUN go get -v github.com/uudashr/gopkgs/v2/cmd/gopkgs && \
    go get -v github.com/ramya-rao-a/go-outline && \
    go get -v github.com/cweill/gotests/gotests && \
    go get -v github.com/fatih/gomodifytags && \
    go get -v github.com/josharian/impl && \
    go get -v github.com/haya14busa/goplay/cmd/goplay && \
    go get -v github.com/go-delve/delve/cmd/dlv && \
    go get -v golang.org/x/tools/gopls && \
    go get -v honnef.co/go/tools/cmd/staticcheck

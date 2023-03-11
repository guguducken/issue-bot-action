FROM golang:1.19 as builder

WORKDIR /workspace

RUN go env -w GOPROXY="https://proxy.golang.org,direct"

COPY release .

# Build
RUN go build -o app main.go

ENTRYPOINT ["/workspace/app"]
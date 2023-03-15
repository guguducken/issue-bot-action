FROM golang:1.19 as builder

WORKDIR /workspace

RUN go env -w GOPROXY="https://proxy.golang.org,direct"

COPY release .

# Build
RUN go build -o app ./pkg/cmd/issue_time_check.go

ENTRYPOINT ["/workspace/app"]
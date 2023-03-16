FROM golang:1.19 as builder

WORKDIR /workspace

RUN go env -w GOPROXY="https://proxy.golang.org,direct"

COPY . .

# Build
RUN go build -tags issue_time -o app ./pkg/cmd/issue_time_check.go

ENTRYPOINT ["/workspace/app"]
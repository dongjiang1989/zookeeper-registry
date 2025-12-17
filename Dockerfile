# Build the manager binary
FROM golang:1.25 AS builder

WORKDIR /workspace

# Copy source files
COPY . .

# Build
RUN GOOS=linux CGO_ENABLED=0 go build -o zookeeper-registry ./main.go


FROM alpine:latest

WORKDIR /
COPY --from=builder /workspace/zookeeper-registry /bin/zookeeper-registry

USER nobody:nobody

LABEL org.opencontainers.image.source="https://github.com/dongjiang1989/zookeeper-registry" \
    org.opencontainers.image.url="https://kubeservice.cn/" \
    org.opencontainers.image.documentation="https://kubeservice.cn/" \
    org.opencontainers.image.licenses="Apache-2.0"
    

ENTRYPOINT ["/bin/zookeeper-registry"]

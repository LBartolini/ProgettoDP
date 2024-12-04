FROM golang:1.22-bookworm AS build

RUN apt-get update && apt-get install -y zip && \
    mkdir /opt/protoc && cd /opt/protoc && wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip && \
    unzip protoc-3.7.0-linux-x86_64.zip

WORKDIR /proto
COPY proto ./
RUN /opt/protoc/bin/protoc --go_out=/proto --proto_path=/proto --go_opt=paths=source_relative --go-grpc_out=/proto --go-grpc_opt=paths=source_relative service.proto

WORKDIR /app
COPY orchestrator ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./main

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /app /root
WORKDIR /root
CMD ["./main"]
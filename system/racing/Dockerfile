FROM golang:1.23-bookworm AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main

FROM build AS test
CMD ["go", "test", "-v", "./..."]

FROM alpine:latest AS run
RUN apk --no-cache add ca-certificates
COPY --from=build /app /root
WORKDIR /root
CMD ["./main"]
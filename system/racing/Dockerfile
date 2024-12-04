FROM golang:1.22-bookworm as build
WORKDIR /app
COPY . ./
COPY ../proto ./proto
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main
CMD ["./main"]

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/main .
CMD ["./main"]
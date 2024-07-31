# from golang:1.18.6 AS builder
FROM golang:1.18.6

WORKDIR /src/
COPY . .
RUN GOOS=linux go build -a -o bin/app *.go

# FROM debian:buster-slim
# WORKDIR /api/
# COPY --from=builder ["/src/bin/app", "/api"]

EXPOSE 8089
CMD ["./bin/app"]

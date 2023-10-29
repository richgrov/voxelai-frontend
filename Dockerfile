FROM golang:1.20 as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/local/bin/app

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/lists/*
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
COPY assets assets
COPY dist dist
COPY templates templates
EXPOSE 8080
ENTRYPOINT ["app"]

FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/app cmd/main/app.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY config.yml /config.yml

CMD ["/app"]
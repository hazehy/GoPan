FROM golang:1.25.7-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY gopan ./gopan
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/gopan ./gopan

FROM alpine:3.21

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata gettext

COPY --from=builder /out/gopan /app/gopan
COPY gopan/etc/gopan-api.docker.yaml /app/etc/gopan-api.docker.yaml
COPY gopan/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 8888

CMD ["/app/entrypoint.sh"]
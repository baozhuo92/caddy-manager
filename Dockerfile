FROM alpine:3.21

RUN apk add --no-cache caddy curl ca-certificates

RUN mkdir -p /app/server /app/caddy_config/sites /app/data

COPY caddy_server /app/caddy_server
COPY config/config.yaml /app/config/config.yaml
COPY templates/ /app/templates/
COPY static/ /app/static/

WORKDIR /app

ENV CONFIG_PATH=/app/config/config.yaml

EXPOSE 80 443 8080

ENTRYPOINT ["./caddy_server"]

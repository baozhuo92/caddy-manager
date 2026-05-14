FROM alpine

RUN mkdir -p /app/manager/html

COPY ./docker/* /app/

RUN chmod +x /app/caddy_linux_amd64

COPY ./dist/* /app/manager/html

EXPOSE 80
EXPOSE 443
WORKDIR /app
CMD ["./caddy_linux_amd64", "run"]
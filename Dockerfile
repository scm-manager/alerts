FROM alpine:3.23.3

RUN apk add --no-cache wget ca-certificates \
 && apk add --no-cache --upgrade zlib \
 && adduser -S -s /bin/false -D -H -u 1000 alerts

COPY alerts /alerts
COPY website/content /content

USER 1000
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

ENTRYPOINT [ "/alerts", "/content" ]

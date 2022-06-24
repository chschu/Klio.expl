# syntax=docker/dockerfile:1.4.2

FROM golang:1.18.3-alpine3.16 AS builder

COPY builder-config/ /

WORKDIR /opt/expl

RUN --mount=type=cache,target=/go/pkg go install

FROM alpine:3.16

LABEL org.opencontainers.image.authors="Christian Schuster <christian@dnup.de>"

RUN apk --no-cache add tzdata

RUN adduser -s /bin/bash -g expl -D -H expl
USER expl

EXPOSE 8000

COPY --from=builder /go/bin/expl /usr/bin/expl
CMD ["/usr/bin/expl"]

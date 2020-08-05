FROM golang:1.14-alpine3.12 AS builder

COPY *.go /go/src/github.com/juli3nk/docker-coturn/

WORKDIR /go/src/github.com/juli3nk/docker-coturn

RUN apk --update add \
		ca-certificates \
		gcc \
		git \
		musl-dev

RUN go get \
	&& go build -ldflags "-linkmode external -extldflags -static -s -w" -o /tmp/coturn-generate-conf


FROM instrumentisto/coturn

COPY --from=builder /tmp/coturn-generate-conf /usr/local/bin/coturn-generate-conf
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

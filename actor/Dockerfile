FROM golang:1.10.3-alpine AS builder
WORKDIR /go/src/github.com/dynastymasra/shajaro/actor
# install dependecies
RUN set -ex \
    && apk add --update bash git curl \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
    && git config --global http.https://gopkg.in.followRedirects true \
    && rm -rf /var/cache/apk/*
COPY . Gopkg.lock Gopkg.toml ./
RUN dep ensure -v

## build linux app source code
RUN GOOS=linux go build -tags=main -o actor

FROM alpine:3.7
RUN set -ex && apk add --update bash ca-certificates tzdata \
            && rm -rf /var/cache/apk/*
# app
WORKDIR /app
COPY --from=builder /go/src/github.com/dynastymasra/shajaro/actor/actor /app/
COPY --from=builder /go/src/github.com/dynastymasra/shajaro/actor/migration /app/migration
COPY --from=builder /go/src/github.com/dynastymasra/shajaro/actor/countries.json /app/
## runtime configs
EXPOSE 8080
ENTRYPOINT ./actor
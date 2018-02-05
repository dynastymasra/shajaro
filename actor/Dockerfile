FROM golang:1.9.3-alpine AS builder
WORKDIR /go/src/sirius/actor
# install dependecies
RUN set -ex \
    && apk add --update bash git glide \
    && git config --global http.https://gopkg.in.followRedirects true \
    && rm -rf /var/cache/apk/*
COPY . glide.lock glide.yaml ./
RUN glide install
RUN glide update
## build app source code
RUN go build -tags=main -o actor

FROM alpine:3.7
RUN set -ex && apk add --update bash ca-certificates tzdata \
            && rm -rf /var/cache/apk/*
# app
WORKDIR /app
COPY --from=builder /go/src/sirius/actor/actor /app/
## runtime configs
EXPOSE 8080
ENTRYPOINT ./actor
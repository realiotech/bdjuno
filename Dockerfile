FROM golang:1.18-alpine AS builder
RUN apk update && apk add --no-cache make git build-base
WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./
RUN go mod download
RUN go mod tidy
RUN make build

FROM alpine:latest
WORKDIR /bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno

ADD bin /bdjuno/bin
ADD hasura /bdjuno/hasura
ENTRYPOINT ["./bin/run.sh"]

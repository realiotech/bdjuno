FROM golang:1.22.5-alpine AS builder
RUN apk update && apk add --no-cache make git build-base
WORKDIR /go/src/github.com/forbole/callisto
COPY . ./
RUN go mod download
RUN go mod tidy
RUN make build

FROM alpine:latest
WORKDIR /callisto
COPY --from=builder /go/src/github.com/forbole/callisto/build/callisto /usr/bin/callisto
RUN curl -L https://github.com/hasura/graphql-engine/raw/stable/cli/get.sh | sh
ADD bin /callisto/bin
ADD hasura /callisto/hasura
ENTRYPOINT ["./bin/run.sh"]

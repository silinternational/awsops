# Ensure this version tracks with go.mod and /Dockerfile-test
FROM golang:1.19 as builder

# Ensure go build env is correct
ENV GOOS linux
ENV GARCH amd64
ENV CGO_ENABLED 0

# copy in cli source
RUN mkdir -p /src
COPY ./ /src

WORKDIR /src

RUN go get ./...
RUN go build -ldflags="-s -w" -o awsops

FROM alpine:3
RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /src/awsops /awsops

ENTRYPOINT ["/awsops"]

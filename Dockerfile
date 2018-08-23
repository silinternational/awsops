FROM golang:latest as builder

# Ensure go build env is correct
ENV GOOS linux
ENV GARCH amd64
ENV CGO_ENABLED 0

# Install deps
RUN go get github.com/golang/dep/cmd/dep

# copy in cli source
RUN mkdir -p /go/src/github.com/silinternational/awsops
COPY ./ /go/src/github.com/silinternational/awsops/

WORKDIR /go/src/github.com/silinternational/awsops

RUN dep ensure
RUN go build -ldflags="-s -w" -o awsops cli/main.go

FROM alpine:latest
RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/silinternational/awsops/awsops /awsops

ENTRYPOINT ["/awsops"]
FROM golang:latest

RUN apt-get update && apt-get install -y awscli
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/mitchellh/gox

WORKDIR /go/src/github.com/silinternational/awsops
COPY ./ /go/src/github.com/silinternational/awsops/
RUN dep ensure

CMD ["go", "test", "./..."]
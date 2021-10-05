FROM golang:latest

RUN apt-get update && apt-get install -y awscli

WORKDIR /go/src/github.com/silinternational/awsops
COPY ./ /go/src/github.com/silinternational/awsops/
RUN go get ./...

CMD ["go", "test", "./..."]

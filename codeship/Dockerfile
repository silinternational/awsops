FROM golang:latest

RUN apt-get update && apt-get install -y awscli

WORKDIR /src
COPY ./ /src
RUN go get ./...

CMD ["go", "test", "./..."]

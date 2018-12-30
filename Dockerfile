FROM golang:1

RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go install -ldflags "-X github.com/cirocosta/cr.version=$(cat ./VERSION)" -v

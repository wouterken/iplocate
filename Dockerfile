FROM golang:1.12.9-stretch

WORKDIR /app

ARG IPLOCATE_PORT='3000'

RUN apt-get update -qq && apt-get install -y --no-install-recommends curl && rm -rf /var/lib/apt/lists/*
RUN go get -v "github.com/oschwald/maxminddb-golang"
COPY ./GeoLite2-City.mmdb .
COPY . .
RUN go build iploc.go

EXPOSE $IPLOCATE_PORT
ENV IPLOCATE_PORT ${IPLOCATE_PORT}

CMD ./iploc -port=$IPLOCATE_PORT

FROM golang:1.11-alpine AS build-stage

RUN apk add --update make
RUN mkdir -p /go/src/github.com/agparadiso/geolocation
WORKDIR /go/src/github.com/agparadiso/geolocation

COPY . /go/src/github.com/agparadiso/geolocation

RUN make build

# Final Stage
FROM alpine

RUN apk --update add ca-certificates
RUN mkdir /app
WORKDIR /app

COPY --from=build-stage  /go/src/github.com/agparadiso/geolocation/geolocation .

EXPOSE 3000

CMD ./geolocation
FROM golang:1.12.6-alpine3.9 AS build
RUN apk add git
WORKDIR /go/src/github.com/aleph-engineering/rocketchat-notification/
COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build github.com/aleph-engineering/rocketchat-notification

FROM alpine:3.9.4
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/aleph-engineering/rocketchat-notification/rocketchat-notification /usr/bin/
CMD ["rocketchat-notification"]

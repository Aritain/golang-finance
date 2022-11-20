FROM golang:alpine as app-builder
WORKDIR /go/src/app
COPY . .
RUN apk add alpine-sdk
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:3.16
COPY --from=app-builder /go/bin/asset_tracker /asset_tracker
ENTRYPOINT ["/asset_tracker"]

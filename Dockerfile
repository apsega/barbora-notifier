# Start by building the application.
FROM golang:1.13-buster as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
EXPOSE 3000
ENTRYPOINT ["/app"]

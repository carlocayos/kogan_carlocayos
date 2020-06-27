FROM golang:1.14-alpine as build

RUN mkdir /myapp
WORKDIR /myapp

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

# Build the binary
RUN go build -o /go/bin/app /myapp/app

# copy only the binary to final image
FROM golang:1.14-alpine
COPY --from=build /go/bin/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]
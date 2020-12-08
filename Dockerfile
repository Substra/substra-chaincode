FROM golang:1.13.8-alpine AS build

WORKDIR /go/src/github.com/chaincode

# Build application
COPY ./chaincode/go.mod ./chaincode/go.sum /go/src/github.com/chaincode/
RUN go mod download

COPY ./chaincode /go/src/github.com/chaincode
RUN go build -o chaincode -v .

# Production ready image
# Pass the binary to the prod image
FROM alpine:3.11 as prod

COPY --from=build /go/src/github.com/chaincode/chaincode /app/chaincode

USER 1000

WORKDIR /app
CMD ./chaincode

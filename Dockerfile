FROM golang:1.14 AS builder
RUN go get github.com/gobuffalo/packr/v2/packr2

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN packr2 install && \
    CGO_ENABLED=0 GOOS=linux packr2 build -a -installsuffix cgo -o phsite

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/phsite /phsite
CMD ["/phsite"]
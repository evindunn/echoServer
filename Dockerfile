FROM golang:1.13-alpine

WORKDIR /app
ADD server.go /app/
RUN go build server.go

ENTRYPOINT ["/app/server"]
CMD ["80"]

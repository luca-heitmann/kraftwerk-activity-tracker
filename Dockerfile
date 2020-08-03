# Build application
FROM golang:1.14.6 AS builder
WORKDIR $GOPATH/src/github.com/luca-heitmann/kraftwerk-activity-tracker
# Update CA certs
RUN update-ca-certificates
# Download modules
COPY src/go.mod .
RUN go mod download
RUN go mod verify
# Compile sources
COPY src .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/main .
# Create target image
FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/main main
CMD ["./main"]
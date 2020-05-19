FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR $GOPATH/src/github.com/YaleOpenLab/opensolar
COPY . .
RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o opensolar
RUN ["cp", "dummyconfig.yaml", "config.yaml"]
RUN ["mv", "config.yaml", "/"]
RUN ["mv", "opensolar", "/"]
WORKDIR /
RUN ["ls"]
# Step 2: build a smaller image
FROM alpine:3.11
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
WORKDIR /
COPY --from=builder /config.yaml .
COPY --from=builder /opensolar .
# EXPOSE 8080
RUN ["ls"]
ENTRYPOINT ["/opensolar"]
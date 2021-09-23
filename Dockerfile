FROM golang:1.17-alpine AS builder

# Install git.
RUN apk add --no-cache git

# Setup directory.
WORKDIR /go/src/github.com/AkashGit21/ms-template
COPY . .

# Compile for Linux.
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

# Install application.
RUN go get ./...
RUN go build -installsuffix cgo \
  -ldflags="-w -s" \
  -o /go/bin/ms-server \
  ./cmd/ms-project

  # Start a fresh image, and only copy the built binary.
FROM scratch
COPY --from=builder /go/bin/ms-server /go/bin/ms-server

# Expose ports
EXPOSE 8081
EXPOSE 8082

# Run the server.
ENTRYPOINT ["/go/bin/ms-server"]
CMD ["run"]
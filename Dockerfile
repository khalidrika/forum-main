# Builder stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o forum -ldflags="-w -s" ./main.go

# Final image
FROM alpine:latest
# Update package lists and install bash
RUN apk update && apk add bash

# For flyctl deploy
ENV PORT 8080

WORKDIR /app
COPY --from=builder /app/forum .
COPY ./database /app/database
COPY ./static /app/static
CMD ["./forum"]

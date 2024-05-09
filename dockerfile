# Stage 1: Build the application (builder)
FROM golang:alpine AS builder
LABEL project-name="cal-api"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o cal-api cmd/main.go

# Stage 2: Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app .

# Install Firefox
RUN apk add --no-cache firefox-esr

RUN apk add --no-cache ca-certificates curl xdg-utils

# Install xdg-utils package to get xdg-open command
RUN apk add --no-cache xdg-utils

# Add xdg-open command to PATH
ENV PATH="/usr/bin:${PATH}"

CMD ["./cal-api"]
EXPOSE 8080

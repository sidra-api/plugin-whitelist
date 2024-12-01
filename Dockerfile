# Gunakan base image untuk Golang
FROM golang:1.23 AS builder

# Set working directory di dalam container
WORKDIR /app

# Copy semua file plugin ke container
COPY . .

# Build binary plugin
RUN go mod tidy && go build -o plugin-whitelist main.go

# Gunakan image minimal untuk hasil akhir
FROM alpine:latest

# Copy binary dari stage builder ke stage ini
COPY --from=builder /app/plugin-whitelist /usr/local/bin/plugin-whitelist

# Jalankan binary
ENTRYPOINT ["/usr/local/bin/plugin-whitelist"]

# Gunakan image ringan Go
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy file project
COPY . .

# Download dependensi
RUN go mod download

# Build aplikasi
RUN go build -o main .

# Jalankan binary
CMD ["/app/main"]

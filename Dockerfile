# Start with the official Golang image
FROM golang:latest

# Install ImageMagick
RUN apt-get update && apt-get install -y imagemagick

# Copy the app files to the container
WORKDIR /app
COPY . .

# Build the app
RUN go build -o file-storage-service

# Start the app
CMD ["./file-storage-service"]

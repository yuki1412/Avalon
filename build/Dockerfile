# Use latest Go version
FROM golang:latest

WORKDIR /app

RUN go install github.com/air-verse/air@latest

# Copy only go.mod first
COPY go.mod go.sum ./

# Copy the rest of the application
COPY . .

RUN go mod download


CMD ["air", "-c", "build/.air.toml"]

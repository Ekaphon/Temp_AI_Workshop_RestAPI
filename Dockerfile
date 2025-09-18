FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server main.go db.go auth.go models.go
EXPOSE 3000
CMD ["./server"]

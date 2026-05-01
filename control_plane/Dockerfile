FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY control_plane/ ./control_plane/
RUN CGO_ENABLED=0 GOOS=linux go build -o /control-plane ./control_plane/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /control-plane .
EXPOSE 8080
CMD ["./control-plane"]

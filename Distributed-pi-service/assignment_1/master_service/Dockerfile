FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go 

FROM alpine:latest AS runtime

COPY --from=builder /app/main /app/main

COPY --from=builder /app/input.txt ./

EXPOSE 8080

CMD ["/app/main", "tail", "-f", "/dev/null"]

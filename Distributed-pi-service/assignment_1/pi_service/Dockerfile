FROM golang:alpine AS builder

ARG CHUNK_NAME

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go 

FROM alpine:latest AS runtime

ARG CHUNK_NAME

COPY ../data/${CHUNK_NAME} ./

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]

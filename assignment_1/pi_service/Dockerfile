FROM  python:3.10-slim AS py_stage

WORKDIR /app

COPY chunking-script.py ./

COPY master_file.txt ./

RUN python3 chunking-script.py

FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go 

FROM alpine:latest AS runtime

COPY --from=py_stage /app/data/${CHUNK_FILE} ./${CHUNK_FILE}

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]

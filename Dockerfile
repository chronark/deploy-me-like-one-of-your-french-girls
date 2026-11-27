FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o server .


FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/server .

RUN apt update
RUN apt install -y ca-certificates curl nodejs npm
RUN npm install -g @openai/codex


EXPOSE 3000
CMD ["./server"]

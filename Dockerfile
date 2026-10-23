FROM golang:latest AS builder
WORKDIR /app
COPY main.go .
RUN go build -o server main.go


FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/server .

RUN sudo apt update
RUN sudo apt install -y curl nodejs npm
RUN npm install -g @openai/codex


EXPOSE 3000
CMD ["./server"]

FROM golang:latest AS builder
WORKDIR /app
COPY main.go .
RUN go build -o server main.go

FROM golang:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 3000
CMD ["./server"]

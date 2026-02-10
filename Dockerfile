FROM alpine:latest AS logs
RUN for i in $(seq 1 10); do \
      echo "[$(date -Iseconds)] log_level=$(shuf -n1 -e INFO WARN ERROR DEBUG) msg=\"$(head -c 32 /dev/urandom | base64 | tr -d '=/+')\""; \
    done

FROM golang:latest

WORKDIR /app
COPY . .

CMD ["go", "run", "."]

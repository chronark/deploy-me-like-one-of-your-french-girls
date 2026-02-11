FROM golang:latest
RUN for i in $(seq 1 10); do \
      echo "[$(date -Iseconds)] log_level=$(shuf -n1 -e INFO WARN ERROR DEBUG) msg=\"$(head -c 32 /dev/urandom | base64 | tr -d '=/+')\""; \
    done

WORKDIR /app
COPY . .

RUN go build -o server .

CMD ["server"]

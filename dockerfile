FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o go-api ./app/main.go

# ========================================================
FROM scratch
WORKDIR /app
COPY --from=builder /app/go-api .
COPY --from=builder /app .

CMD ["/app/go-api"]

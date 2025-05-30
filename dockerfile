FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o go-api ./main.go

# ========================================================
FROM scratch
WORKDIR /app
COPY --from=builder /app/go-api .
COPY --from=builder /app .

CMD ["/app/go-api"]
FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o go-api ./app/main.go

# ========================================================
FROM scratch
WORKDIR /app
COPY --from=builder /app/go-api .
COPY --from=builder /app .

CMD ["/app/go-api"]

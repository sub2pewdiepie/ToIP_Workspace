FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o go-api ./app/main.go

# ========================================================
FROM scratch

COPY --from=builder /app/go-api /go-api
COPY --from=builder /app .

WORKDIR /
EXPOSE 8080
CMD ["/go-api"]
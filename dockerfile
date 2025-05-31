FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# RUN go mod tidy
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
     CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o go-api ./main.go
#RUN go build -o go-api ./main.go

# ========================================================
FROM scratch
WORKDIR /app
COPY --from=builder /app/go-api .
COPY --from=builder /app .
CMD ["/app/go-api"]
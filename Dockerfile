# ---------- BUILD STAGE ----------
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Install git (needed for go mod if private deps)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app

# ---------- RUNTIME STAGE ----------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/config.yaml .

# create directories inside container
RUN mkdir -p /app/input /app/output

EXPOSE 8080

CMD ["./app"]

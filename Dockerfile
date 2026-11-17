FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/forhum .

FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/forhum /app/bin/forhum
COPY --from=builder /app/client/templates /app/client/templates
COPY --from=builder /app/static /app/static

EXPOSE 8080

CMD ["/app/bin/forhum"]
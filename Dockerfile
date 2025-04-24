FROM golang:1.23.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ocr-demo-be .

FROM alpine:3.17

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/ocr-demo-be .

RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

CMD ["./ocr-demo-be"]
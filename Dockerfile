#BUILD Stage

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/server .

# Final run time image

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["/app/server"]
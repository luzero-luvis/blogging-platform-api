# builder
FROM golang:1.25.5-alpine3.21 AS builder

RUN apk add --no-cache ca-certificates 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -trimpath \
  -o blog-api \
  .

#runtime
FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs

COPY --from=builder /etc/passwd  /etc/passwd

COPY --from=builder /etc/group /etc/group

WORKDIR /app

COPY --from=builder /app/blog-api .

USER nobody:nobody

EXPOSE 5000

ENTRYPOINT [ "./blog-api" ]


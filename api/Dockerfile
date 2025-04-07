FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api .

FROM alpine:latest
RUN apk add --no-cache tini
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -q -O- http://localhost:8080/metrics || exit 1
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["./api"]
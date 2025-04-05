FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o api .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
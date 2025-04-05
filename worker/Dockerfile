FROM golang:1.21 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o worker .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/worker .
COPY --from=builder /app/config ./config

EXPOSE 9090

# HEALTHCHECK --interval=30s --timeout=3s \
#  CMD wget -q -O- http://localhost:9090/metrics || exit 1

CMD [ "./worker" ]
FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o worker .

FROM alpine:latest
RUN apk add --no-cache tini
WORKDIR /root/
COPY --from=builder /app/worker .
EXPOSE 9091
HEALTHCHECK --interval=30s --timeout=3s \
 CMD wget -q -O- http://localhost:9091/metrics || exit 1
ENTRYPOINT [ "/sbin/tini", "--" ]
CMD [ "./worker" ]
# Build stage.
FROM golang:1.14.4 AS builder
WORKDIR /src
COPY . .
COPY /static /bin/static
COPY /data /bin/data
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/server .

# Final stage.
FROM alpine:3.12.0
ENV PORT 42069
WORKDIR /app
EXPOSE 42069
COPY --from=builder /bin/ .
CMD ["./server"]

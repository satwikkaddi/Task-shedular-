# Build stage
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker ./cmd/worker

# Runtime
FROM alpine:latest
COPY --from=build /bin/server /bin/server
COPY --from=build /bin/worker /bin/worker
EXPOSE 8080
CMD ["/bin/server"]
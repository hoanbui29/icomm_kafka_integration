# Stage 1: Build the Go binary
FROM golang:1.23.5 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main .

# Stage 2: Run with a lightweight image
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/main /main

ENTRYPOINT ["/main"]


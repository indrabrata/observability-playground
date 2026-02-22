FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a static binary.
# Alpine's gcc already links against musl libc system-wide, so CC=musl-gcc is not needed.
ENV CGO_ENABLED=1
RUN go build -ldflags="-linkmode external -extldflags '-static'" -o server .
RUN chmod +x server


FROM scratch AS runtime

WORKDIR /app

COPY --from=builder /app/server .

# SQLite file lives on a mounted volume (/data), not inside the image.
# Docker creates the mount point automatically when a volume is attached.

EXPOSE 8080

CMD ["./server"]

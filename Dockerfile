# Stage 1: Build the Go app
FROM golang:1.21-alpine AS builder

# Set timezone as an ARG
ARG TZ=Europe/Oslo
ENV TZ=${TZ}

# Install dependencies
RUN apk add --no-cache git tzdata

# Set working dir inside container
WORKDIR /app/

# Clone Repo
RUN git clone https://github.com/skandix/Beetroot.git


# Set working dir inside container
WORKDIR /app/Beetroot


# Tidy up Go modules
RUN go mod tidy

# Set working dir inside beetroot
WORKDIR /app/Beetroot/cmd/beetroot/

# Build the binary
RUN go build -o beetroot .

# Stage 2: Runtime
FROM alpine:3.18

# Set timezone
ARG TZ=Europe/Oslo
ENV TZ=${TZ}
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone


# Set working dir /app/
WORKDIR /app/

# Copy binary from builder
COPY --from=builder /app/Beetroot/cmd/beetroot/beetroot /app/


# Expose port
EXPOSE 8080

# Run it
CMD ["./beetroot"]

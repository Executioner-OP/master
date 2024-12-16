# Use Ubuntu 20.04 as the base image
FROM ubuntu:20.04

# Set environment variable to non-interactive to prevent prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Update package list and install dependencies
RUN apt-get update && apt-get install -y  \
    wget \
    libcap-dev \
    asciidoc-base \
    libsystemd-dev \
    pkg-config \
    git \
    build-essential && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Download and install Go 1.23.4
RUN wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz && \
    rm go1.23.4.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Set working directory
WORKDIR /app

# Copy Go code and related files from the host to the container
COPY . /app

# Ensure proper permissions for the Go app
RUN chmod -R 755 /app

# Verify Go installation and download dependencies
RUN go version && go mod tidy

# Expose ports 9001 (RPC) and 3000 (HTTP)
EXPOSE 9001 3000

# Set entrypoint for isolate command (optional)
ENTRYPOINT ["go"]

# Default command to run the main Go application
CMD ["run", "./main.go"]

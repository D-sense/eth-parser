# Build the Go Binary.
FROM golang:1.20 as build_eth_parser
ENV CGO_ENABLED 0
ARG BUILD_REF

# Create the core directory and the copy the module files first and then
# download the dependencies. If this doesn't change, we won't need to do this
# again in future builds.
RUN mkdir /server
COPY go.* /api/
WORKDIR /api/server
RUN go mod download

# Copy the source code into the container.
COPY . /api

# Build the service binary.
WORKDIR /api/api
RUN go build -o eth_parser .

# Run the Go Binary in Alpine.
FROM alpine:3.14
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S user && \
    adduser -u 1000 -h /server -G user -S user

COPY --from=build_eth_parser /api /api
WORKDIR /api/api
USER user
CMD ["./eth_parser"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="ethereum-parser-api" \
      org.opencontainers.image.authors="Adesina Hammed H. <delameh@icloud.com>" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Trust Wallet"



